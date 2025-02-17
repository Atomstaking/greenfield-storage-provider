package test

import (
	"errors"
	"fmt"
	golog "log"
	"math/big"
	"time"

	"github.com/forbole/juno/v4/common"
	"github.com/forbole/juno/v4/models"
	"github.com/spaolacci/murmur3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/bnb-chain/greenfield-storage-provider/store/bsdb"
)

var verifyFuncs = []func(db *gorm.DB) error{verify1, verify2, verify3, verify4, verify5, verify6, verify7, verify8, verify9, verify10, verify11, verify12, verify13, verify14, verify15, verify16, verify17, verify18, verify19, verify20, verify21, verify22, verify23, verify24, verify25}

func Verify() error {
	dsn := "root:root@tcp(localhost:3306)/block_syncer?charset=utf8mb4&parseTime=True&loc=Local&interpolateParams=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	tick := time.NewTicker(time.Millisecond * 500)
	defer tick.Stop()
	for i := 1; i <= LatestHeight; i++ {
		tryCount := 0
		for tryCount < 10 {
			<-tick.C
			var epoch bsdb.Epoch
			db.Table("epoch").First(&epoch)
			tryCount++
			if epoch.BlockHeight == int64(i) {
				break
			}
		}

		// verify data
		f := verifyFuncs[i-1]
		if err := f(db); err != nil {
			return err
		}

		golog.Printf("%d case pass", i)

		// height increase
		StatusRes = fmt.Sprintf("{\"sync_info\":{\"latest_block_height\":\"%d\"}}", i+1)
	}

	// verify height
	var epoch bsdb.Epoch
	db.Table("epoch").First(&epoch)
	if epoch.BlockHeight != int64(LatestHeight) {
		return errors.New("height error")
	}

	return nil
}

func verify1(db *gorm.DB) error {
	var bucket models.Bucket
	if err := db.Table(bsdb.BucketTableName).Where("bucket_name = ?", "full-node-v1-acc0000000072-buc0000000000").Find(&bucket).Error; err != nil {
		return err
	}
	golog.Println(bucket)
	if bucket.Status != "BUCKET_STATUS_CREATED" {
		return errors.New("bucket status error")
	}
	return nil
}

func verify2(db *gorm.DB) error {
	var bucket models.Bucket
	if err := db.Table(bsdb.BucketTableName).Where("bucket_name = ?", "full-node-v1-acc0000000072-buc0000000000").Find(&bucket).Error; err != nil {
		return err
	}
	if bucket.Visibility != "VISIBILITY_TYPE_PUBLIC_READ" {
		return errors.New("bucket Visibility is error")
	}
	return nil
}
func verify3(db *gorm.DB) error {
	var bucket models.Bucket
	if err := db.Table(bsdb.BucketTableName).Where("bucket_name = ?", "full-node-v1-acc0000000072-buc0000000000").Find(&bucket).Error; err != nil {
		return err
	}
	if !bucket.Removed {
		return errors.New("bucket is not remove")
	}
	return nil
}
func verify4(db *gorm.DB) error {
	var bucket models.Bucket
	if err := db.Table(bsdb.BucketTableName).Where("bucket_name = ?", "cxz").Find(&bucket).Error; err != nil {
		return err
	}
	if bucket.DeleteAt != int64(4) {
		return fmt.Errorf("bucket is not discontinue delete at:%v", bucket.DeleteAt)
	}
	return nil
}
func verify5(db *gorm.DB) error {
	var bucket models.Bucket
	if err := db.Table(bsdb.BucketTableName).Where("bucket_name = ?", "cxz").Find(&bucket).Error; err != nil {
		return err
	}
	if bucket.GlobalVirtualGroupFamilyId != 9 {
		return errors.New("bucket GlobalVirtualGroupFamilyId error")
	}
	return nil
}
func verify6(db *gorm.DB) error {
	var bucket models.Bucket
	if err := db.Table(bsdb.BucketTableName).Where("bucket_name = ?", "cxz").Find(&bucket).Error; err != nil {
		return err
	}
	if bucket.GlobalVirtualGroupFamilyId != 10 {
		return errors.New("bucket GlobalVirtualGroupFamilyId error")
	}
	return nil
}
func verify7(db *gorm.DB) error {
	var g models.Group
	if err := db.Table(bsdb.GroupTableName).Where("group_id = ?", common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002")).Find(&g).Error; err != nil {
		return err
	}
	return nil
}
func verify8(db *gorm.DB) error {
	var g models.Group
	if err := db.Table(bsdb.GroupTableName).Where("group_id = ? and account_id = ?", common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"), common.HexToAddress("0x92b7976702C064C7e2e791854497Ec73C853CEB5")).Find(&g).Error; err != nil {
		return err
	}
	if g.ExpirationTime != 253402300799 {
		return errors.New("member ExpirationTime error")
	}
	return nil
}
func verify9(db *gorm.DB) error {
	var g models.Group
	if err := db.Table(bsdb.GroupTableName).Where("group_id = ? and account_id = ?", common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"), common.HexToAddress("0x4d57d300AfaF9f407e26552965ce355786206cF4")).Find(&g).Error; err != nil {
		return err
	}
	if g.ExpirationTime != 1704067199 {
		return fmt.Errorf("member ExpirationTime error ExpirationTime:%d account id:%v", g.ExpirationTime, g.AccountID)
	}
	return nil
}
func verify10(db *gorm.DB) error {
	var g models.Group
	if err := db.Table(bsdb.GroupTableName).Where("group_id = ? and account_id = ?", common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000002"), common.HexToAddress("0x0000000000000000000000000000000000000000")).Find(&g).Error; err != nil {
		return err
	}
	if !g.Removed {
		return errors.New("group is not removed")
	}
	return nil
}
func verify11(db *gorm.DB) error {
	table := GetObjectsTableName("cxz")
	var o models.Object
	if err := db.Table(table).Where("object_id = ?", common.BigToHash(big.NewInt(1201))).Find(&o).Error; err != nil {
		return err
	}
	return nil
}
func verify12(db *gorm.DB) error {
	table := GetObjectsTableName("cxz")
	var o models.Object
	if err := db.Table(table).Where("object_id = ?", common.BigToHash(big.NewInt(1201))).Find(&o).Error; err != nil {
		return err
	}
	if o.Status != "OBJECT_STATUS_SEALED" {
		return errors.New("object status error")
	}
	return nil
}
func verify13(db *gorm.DB) error {
	table := GetObjectsTableName("cxz")
	var o models.Object
	if err := db.Table(table).Where("object_id = ?", common.BigToHash(big.NewInt(1201))).Find(&o).Error; err != nil {
		return err
	}
	if !o.Removed {
		return errors.New("object is not removed")
	}
	return nil
}
func verify14(db *gorm.DB) error {
	var sr models.StreamRecord
	if err := db.Table((&models.StreamRecord{}).TableName()).Where("account = ?", common.HexToAddress("0x51D2c7dFa2a58C73853c2657Dc505e7430879B1b")).Find(&sr).Error; err != nil {
		return err
	}
	return nil
}
func verify15(db *gorm.DB) error {
	var p models.PaymentAccount
	if err := db.Table((&models.PaymentAccount{}).TableName()).Where("addr = ?", common.HexToAddress("0x68bd245Df652435321989b999F9F70Cd31281b66")).Find(&p).Error; err != nil {
		return err
	}
	return nil
}
func verify16(db *gorm.DB) error {
	var lvg models.LocalVirtualGroup
	if err := db.Table((&models.LocalVirtualGroup{}).TableName()).Where("local_virtual_group_id = ?", 8).Find(&lvg).Error; err != nil {
		return err
	}
	return nil
}
func verify17(db *gorm.DB) error {
	var lvg models.LocalVirtualGroup
	if err := db.Table((&models.LocalVirtualGroup{}).TableName()).Where("local_virtual_group_id = ?", 8).Find(&lvg).Error; err != nil {
		return err
	}
	if lvg.StoredSize != 0 {
		return errors.New("StoredSize error")
	}
	return nil
}
func verify18(db *gorm.DB) error {
	var gvg bsdb.GlobalVirtualGroup
	if err := db.Table((&bsdb.GlobalVirtualGroup{}).TableName()).Where("global_virtual_group_id = ?", 1).Find(&gvg).Error; err != nil {
		return err
	}
	return nil
}
func verify19(db *gorm.DB) error {
	var gvg bsdb.GlobalVirtualGroup
	if err := db.Table((&bsdb.GlobalVirtualGroup{}).TableName()).Where("global_virtual_group_id = ?", 1).Find(&gvg).Error; err != nil {
		return err
	}
	if gvg.StoredSize != 48372064 || gvg.TotalDeposit.Raw().String() != "140737488355328000" {
		return fmt.Errorf("gvg update error StoredSize: %v, TotalDeposit:%v", gvg.StoredSize, gvg.TotalDeposit.Raw().String())
	}
	return nil
}
func verify20(db *gorm.DB) error {
	var gvgf bsdb.GlobalVirtualGroupFamily
	if err := db.Table((&bsdb.GlobalVirtualGroupFamily{}).TableName()).Where("global_virtual_group_family_id = ?", 1).Find(&gvgf).Error; err != nil {
		return err
	}
	return nil
}
func verify21(db *gorm.DB) error {
	var gvgf bsdb.GlobalVirtualGroupFamily
	if err := db.Table((&bsdb.GlobalVirtualGroupFamily{}).TableName()).Where("global_virtual_group_family_id = ?", 1).Find(&gvgf).Error; err != nil {
		return err
	}
	if gvgf.PrimarySpId != 3 {
		return errors.New("PrimarySpId error")
	}
	return nil
}
func verify22(db *gorm.DB) error {
	var gvgf bsdb.GlobalVirtualGroupFamily
	if err := db.Table((&bsdb.GlobalVirtualGroupFamily{}).TableName()).Where("global_virtual_group_family_id = ?", 1).Find(&gvgf).Error; err != nil {
		return err
	}
	if !gvgf.Removed {
		return errors.New("gvgf is not removed")
	}
	return nil
}
func verify23(db *gorm.DB) error {
	var sp models.StorageProvider
	if err := db.Table((&models.StorageProvider{}).TableName()).Where("sp_id = ?", 14).Find(&sp).Error; err != nil {
		return err
	}
	return nil
}
func verify24(db *gorm.DB) error {
	var sp models.StorageProvider
	if err := db.Table((&models.StorageProvider{}).TableName()).Where("sp_id = ?", 14).Find(&sp).Error; err != nil {
		return err
	}
	if sp.Endpoint != "http://spxrmfl.greenfield.io" || sp.BlsKey != "b689357b256f8aabaf02fceb56a9a61c59b2d9b3cc78d4413fefd1e3bd902c90dcc7b346deb672d164c9ff832a8ee1d9" {
		return fmt.Errorf("sp update failed endpoint: %s, blskey:%s", sp.Endpoint, sp.BlsKey)
	}
	return nil
}
func verify25(db *gorm.DB) error {
	var sp models.StorageProvider
	if err := db.Table((&models.StorageProvider{}).TableName()).Where("sp_id = ?", 14).Find(&sp).Error; err != nil {
		return err
	}
	if !sp.Removed {
		return errors.New("sp is not removed")
	}
	return nil
}

const ObjectsNumberOfShards = 64
const ObjectTableName = "objects"

func GetObjectsTableName(bucketName string) string {
	return GetObjectsTableNameByShardNumber(int(GetObjectsShardNumberByBucketName(bucketName)))
}

func GetObjectsShardNumberByBucketName(bucketName string) uint32 {
	return murmur3.Sum32([]byte(bucketName)) % ObjectsNumberOfShards
}

func GetObjectsTableNameByShardNumber(shard int) string {
	return fmt.Sprintf("%s_%02d", ObjectTableName, shard)
}
