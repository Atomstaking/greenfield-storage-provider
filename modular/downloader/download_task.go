package downloader

import (
	"context"
	"errors"
	"net/http"

	"github.com/bnb-chain/greenfield-storage-provider/base/types/gfsperrors"
	"github.com/bnb-chain/greenfield-storage-provider/core/module"
	"github.com/bnb-chain/greenfield-storage-provider/core/piecestore"
	"github.com/bnb-chain/greenfield-storage-provider/core/spdb"
	corespdb "github.com/bnb-chain/greenfield-storage-provider/core/spdb"
	"github.com/bnb-chain/greenfield-storage-provider/core/task"
	"github.com/bnb-chain/greenfield-storage-provider/core/taskqueue"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"github.com/bnb-chain/greenfield-storage-provider/store/sqldb"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
)

var (
	ErrDanglingPointer   = gfsperrors.Register(module.DownloadModularName, http.StatusInternalServerError, 30001, "OoooH.... request lost")
	ErrObjectUnsealed    = gfsperrors.Register(module.DownloadModularName, http.StatusBadRequest, 30002, "object unsealed")
	ErrRepeatedTask      = gfsperrors.Register(module.DownloadModularName, http.StatusBadRequest, 30003, "request repeated")
	ErrExceedBucketQuota = gfsperrors.Register(module.DownloadModularName, http.StatusBadRequest, 30004, "bucket quota overflow")
	ErrExceedQueue       = gfsperrors.Register(module.DownloadModularName, http.StatusServiceUnavailable, 30005, "request exceed the limit, try again later")
	ErrInvalidParam      = gfsperrors.Register(module.DownloadModularName, http.StatusBadRequest, 30006, "request params invalid")
	ErrNoSuchPiece       = gfsperrors.Register(module.DownloadModularName, http.StatusBadRequest, 30007, "request params invalid, no such piece")
	ErrPieceStore        = gfsperrors.Register(module.DownloadModularName, http.StatusInternalServerError, 35101, "server slipped away, try again later")
	ErrGfSpDB            = gfsperrors.Register(module.DownloadModularName, http.StatusInternalServerError, 35201, "server slipped away, try again later")
)

func (d *DownloadModular) PreDownloadObject(ctx context.Context, downloadObjectTask task.DownloadObjectTask) error {
	if downloadObjectTask == nil || downloadObjectTask.GetObjectInfo() == nil || downloadObjectTask.GetStorageParams() == nil {
		log.CtxErrorw(ctx, "failed pre download object due to pointer dangling")
		return ErrDanglingPointer
	}
	if downloadObjectTask.GetObjectInfo().GetObjectStatus() != storagetypes.OBJECT_STATUS_SEALED {
		log.CtxErrorw(ctx, "failed to pre download object due to object unsealed")
		return ErrObjectUnsealed
	}
	if d.downloadQueue.Has(downloadObjectTask.Key()) {
		log.CtxErrorw(ctx, "failed to pre download object due to task repeated")
		return ErrRepeatedTask
	}
	// TODO:: spilt check and add record steps, check in pre download, add record in post download
	if err := d.baseApp.GfSpDB().CheckQuotaAndAddReadRecord(
		&spdb.ReadRecord{
			BucketID:        downloadObjectTask.GetBucketInfo().Id.Uint64(),
			ObjectID:        downloadObjectTask.GetObjectInfo().Id.Uint64(),
			UserAddress:     downloadObjectTask.GetUserAddress(),
			BucketName:      downloadObjectTask.GetBucketInfo().GetBucketName(),
			ObjectName:      downloadObjectTask.GetObjectInfo().GetObjectName(),
			ReadSize:        uint64(downloadObjectTask.GetSize()),
			ReadTimestampUs: sqldb.GetCurrentTimestampUs(),
		},
		&spdb.BucketQuota{
			ReadQuotaSize: downloadObjectTask.GetBucketInfo().GetChargedReadQuota() + d.bucketFreeQuota,
		},
	); err != nil {
		log.CtxErrorw(ctx, "failed to check bucket quota", "error", err)
		if errors.Is(err, sqldb.ErrCheckQuotaEnough) {
			return ErrExceedBucketQuota
		}
		// ignore the access db error, it is the system's inner error, will be let the request go.
	}
	// report the task to the manager for monitor the download task
	d.baseApp.GfSpClient().ReportTask(ctx, downloadObjectTask)
	return nil
}

func (d *DownloadModular) HandleDownloadObjectTask(ctx context.Context, downloadObjectTask task.DownloadObjectTask) ([]byte, error) {
	var err error
	defer func() {
		if err != nil {
			downloadObjectTask.SetError(err)
		}
		log.CtxDebugw(ctx, downloadObjectTask.Info())
	}()
	if err = d.downloadQueue.Push(downloadObjectTask); err != nil {
		log.CtxErrorw(ctx, "failed to push download queue", "error", err)
		return nil, ErrExceedQueue
	}
	defer d.downloadQueue.PopByKey(downloadObjectTask.Key())
	pieceInfos, err := SplitToSegmentPieceInfos(downloadObjectTask, d.baseApp.PieceOp())
	if err != nil {
		log.CtxErrorw(ctx, "failed to generate piece info to download", "error", err)
		return nil, err
	}
	var data []byte
	for _, pInfo := range pieceInfos {
		piece, err := d.baseApp.PieceStore().GetPiece(ctx, pInfo.SegmentPieceKey,
			int64(pInfo.Offset), int64(pInfo.Length))
		if err != nil {
			log.CtxErrorw(ctx, "failed to get piece data from piece store", "error", err)
			return nil, ErrPieceStore
		}
		data = append(data, piece...)
	}
	return data, nil
}

type SegmentPieceInfo struct {
	SegmentPieceKey string
	Offset          uint64
	Length          uint64
}

func SplitToSegmentPieceInfos(downloadObjectTask task.DownloadObjectTask, op piecestore.PieceOp) ([]*SegmentPieceInfo, error) {
	if downloadObjectTask.GetObjectInfo().GetPayloadSize() == 0 ||
		downloadObjectTask.GetLow() >= int64(downloadObjectTask.GetObjectInfo().GetPayloadSize()) ||
		downloadObjectTask.GetHigh() >= int64(downloadObjectTask.GetObjectInfo().GetPayloadSize()) ||
		downloadObjectTask.GetHigh() < downloadObjectTask.GetLow() {
		log.Errorw("failed to parse params", "object_size",
			downloadObjectTask.GetObjectInfo().GetPayloadSize(), "low", downloadObjectTask.GetLow(), "high", downloadObjectTask.GetHigh())
		return nil, ErrInvalidParam
	}
	segmentSize := downloadObjectTask.GetStorageParams().VersionedParams.GetMaxSegmentSize()
	segmentCount := op.SegmentPieceCount(downloadObjectTask.GetObjectInfo().GetPayloadSize(),
		downloadObjectTask.GetStorageParams().VersionedParams.GetMaxSegmentSize())
	var (
		pieceInfos []*SegmentPieceInfo
		low        = uint64(downloadObjectTask.GetLow())
		high       = uint64(downloadObjectTask.GetHigh())
	)
	for segmentPieceIndex := uint64(0); segmentPieceIndex < uint64(segmentCount); segmentPieceIndex++ {
		currentStart := segmentPieceIndex * segmentSize
		currentEnd := (segmentPieceIndex+1)*segmentSize - 1
		if low > currentEnd {
			continue
		}
		if low > currentStart {
			currentStart = low
		}

		if high <= currentEnd {
			currentEnd = high
			offsetInPiece := currentStart - (segmentPieceIndex * segmentSize)
			lengthInPiece := currentEnd - currentStart + 1
			pieceInfos = append(pieceInfos, &SegmentPieceInfo{
				SegmentPieceKey: op.SegmentPieceKey(
					downloadObjectTask.GetObjectInfo().Id.Uint64(),
					uint32(segmentPieceIndex)),
				Offset: offsetInPiece,
				Length: lengthInPiece,
			})
			// break to finish
			break
		} else {
			offsetInPiece := currentStart - (segmentPieceIndex * segmentSize)
			lengthInPiece := currentEnd - currentStart + 1
			pieceInfos = append(pieceInfos, &SegmentPieceInfo{
				SegmentPieceKey: op.SegmentPieceKey(
					downloadObjectTask.GetObjectInfo().Id.Uint64(),
					uint32(segmentPieceIndex)),
				Offset: offsetInPiece,
				Length: lengthInPiece,
			})
		}
	}
	return pieceInfos, nil
}

func (d *DownloadModular) PostDownloadObject(ctx context.Context, downloadObjectTask task.DownloadObjectTask) {
}

func (d *DownloadModular) PreDownloadPiece(ctx context.Context, downloadPieceTask task.DownloadPieceTask) error {
	if downloadPieceTask == nil || downloadPieceTask.GetObjectInfo() == nil || downloadPieceTask.GetStorageParams() == nil {
		log.CtxErrorw(ctx, "failed pre download piece due to pointer dangling")
		return ErrDanglingPointer
	}
	if downloadPieceTask.GetObjectInfo().GetObjectStatus() != storagetypes.OBJECT_STATUS_SEALED {
		log.CtxErrorw(ctx, "failed to pre download piece due to object unsealed")
		return ErrObjectUnsealed
	}
	if d.downloadQueue.Has(downloadPieceTask.Key()) {
		log.CtxErrorw(ctx, "failed to pre download piece due to task repeated")
		return ErrRepeatedTask
	}

	if downloadPieceTask.GetEnableCheck() {
		if err := d.baseApp.GfSpDB().CheckQuotaAndAddReadRecord(
			&spdb.ReadRecord{
				BucketID:        downloadPieceTask.GetBucketInfo().Id.Uint64(),
				ObjectID:        downloadPieceTask.GetObjectInfo().Id.Uint64(),
				UserAddress:     downloadPieceTask.GetUserAddress(),
				BucketName:      downloadPieceTask.GetBucketInfo().GetBucketName(),
				ObjectName:      downloadPieceTask.GetObjectInfo().GetObjectName(),
				ReadSize:        downloadPieceTask.GetTotalSize(),
				ReadTimestampUs: sqldb.GetCurrentTimestampUs(),
			},
			&spdb.BucketQuota{
				ReadQuotaSize: downloadPieceTask.GetBucketInfo().GetChargedReadQuota() + d.bucketFreeQuota,
			},
		); err != nil {
			log.CtxErrorw(ctx, "failed to check bucket quota", "error", err)
			if errors.Is(err, sqldb.ErrCheckQuotaEnough) {
				return ErrExceedBucketQuota
			}
			// ignore the access db error, it is the system's inner error, will be let the request go.
		}
	}
	// report the task to the manager for monitor the download piece task
	d.baseApp.GfSpClient().ReportTask(ctx, downloadPieceTask)
	return nil
}

func (d *DownloadModular) HandleDownloadPieceTask(ctx context.Context, downloadPieceTask task.DownloadPieceTask) ([]byte, error) {
	var (
		err       error
		pieceData []byte
	)

	defer func() {
		if err != nil {
			downloadPieceTask.SetError(err)
		}
		log.CtxDebugw(ctx, downloadPieceTask.Info())
	}()
	if err = d.downloadQueue.Push(downloadPieceTask); err != nil {
		log.CtxErrorw(ctx, "failed to push download queue", "error", err)
		return nil, ErrExceedQueue
	}
	defer d.downloadQueue.PopByKey(downloadPieceTask.Key())

	if pieceData, err = d.baseApp.PieceStore().GetPiece(ctx, downloadPieceTask.GetPieceKey(),
		int64(downloadPieceTask.GetPieceOffset()), int64(downloadPieceTask.GetPieceLength())); err != nil {
		log.CtxErrorw(ctx, "failed to get piece data from piece store", "task_info", downloadPieceTask.Info(), "error", err)
		return nil, ErrPieceStore
	}
	return pieceData, nil
}

func (d *DownloadModular) PostDownloadPiece(ctx context.Context, downloadPieceTask task.DownloadPieceTask) {
}

func (d *DownloadModular) PreChallengePiece(ctx context.Context, downloadPieceTask task.ChallengePieceTask) error {
	if downloadPieceTask == nil || downloadPieceTask.GetObjectInfo() == nil {
		log.CtxErrorw(ctx, "failed to pre challenge piece due to pointer dangling")
		return ErrDanglingPointer
	}
	if downloadPieceTask.GetObjectInfo().GetObjectStatus() != storagetypes.OBJECT_STATUS_SEALED {
		log.CtxErrorw(ctx, "failed to pre challenge piece due to object unsealed")
		return ErrObjectUnsealed
	}
	d.baseApp.GfSpClient().ReportTask(ctx, downloadPieceTask)
	return nil
}

func (d *DownloadModular) HandleChallengePiece(ctx context.Context, downloadPieceTask task.ChallengePieceTask) (
	[]byte, [][]byte, []byte, error) {
	var (
		err       error
		integrity *corespdb.IntegrityMeta
		data      []byte
	)
	defer func() {
		if err != nil {
			downloadPieceTask.SetError(err)
		}
		log.CtxDebugw(ctx, downloadPieceTask.Info())
	}()
	if err = d.challengeQueue.Push(downloadPieceTask); err != nil {
		log.CtxErrorw(ctx, "failed to push challenge piece queue", "error", err)
		return nil, nil, nil, ErrExceedQueue
	}
	defer d.challengeQueue.PopByKey(downloadPieceTask.Key())
	pieceKey := d.baseApp.PieceOp().ChallengePieceKey(
		downloadPieceTask.GetObjectInfo().Id.Uint64(),
		downloadPieceTask.GetSegmentIdx(),
		downloadPieceTask.GetRedundancyIdx())
	integrity, err = d.baseApp.GfSpDB().GetObjectIntegrity(downloadPieceTask.GetObjectInfo().Id.Uint64())
	if err != nil {
		log.CtxErrorw(ctx, "failed to get integrity hash", "error", err)
		return nil, nil, nil, ErrGfSpDB
	}
	if int(downloadPieceTask.GetSegmentIdx()) >= len(integrity.PieceChecksumList) {
		log.CtxErrorw(ctx, "failed to get challenge info due to segment index wrong")
		return nil, nil, nil, ErrNoSuchPiece
	}
	data, err = d.baseApp.PieceStore().GetPiece(ctx, pieceKey, 0, -1)
	if err != nil {
		log.CtxErrorw(ctx, "failed to get piece data", "error", err)
		return nil, nil, nil, ErrPieceStore
	}
	return integrity.IntegrityChecksum, integrity.PieceChecksumList, data, nil
}

func (d *DownloadModular) PostChallengePiece(ctx context.Context, downloadPieceTask task.ChallengePieceTask) {
}

func (d *DownloadModular) QueryTasks(ctx context.Context, subKey task.TKey) (
	[]task.Task, error) {
	downloadTasks, _ := taskqueue.ScanTQueueBySubKey(d.downloadQueue, subKey)
	challengeTasks, _ := taskqueue.ScanTQueueBySubKey(d.challengeQueue, subKey)
	return append(downloadTasks, challengeTasks...), nil
}