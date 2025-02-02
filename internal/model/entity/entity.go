package entity

import (
	"cynxhostagent/internal/constant/types"
	"time"
)

type TblUser struct {
	Id          int       `gorm:"primaryKey" visibility:"1"`
	Username    string    `gorm:"size:255;not null;unique" visibility:"2"`
	Password    string    `gorm:"size:255;not null" visibility:"10"`
	Coin        int       `gorm:"default:0" visibility:"2"`
	CreatedDate time.Time `gorm:"autoCreateTime" visibility:"1"`
	UpdatedDate time.Time `gorm:"autoUpdateTime" visibility:"1"`
}

type TblScript struct {
	Id             int       `gorm:"primaryKey" visibility:"1"`
	CreatedDate    time.Time `gorm:"autoCreateTime" visibility:"1"`
	UpdatedDate    time.Time `gorm:"autoUpdateTime" visibility:"1"`
	Name           string    `gorm:"size:255;not null" visibility:"1"`
	Variables      string    `gorm:"type:json;not null" visibility:"2"`
	SetupScript    string    `gorm:"type:text;not null" visibility:"2"`
	StartScript    string    `gorm:"type:text;not null" visibility:"2"`
	StopScript     string    `gorm:"type:text;not null" visibility:"2"`
	ShutdownScript string    `gorm:"type:text;not null" visibility:"2"`
}

type TblServerTemplate struct {
	Id          int       `gorm:"primaryKey" visibility:"1"`
	CreatedDate time.Time `gorm:"autoCreateTime" visibility:"1"`
	UpdatedDate time.Time `gorm:"autoUpdateTime" visibility:"1"`
	Name        string    `gorm:"size:255;not null" visibility:"1"`
	MinimumRam  int       `gorm:"not null" visibility:"1"`
	MinimumCpu  int       `gorm:"not null" visibility:"1"`
	MinimumDisk int       `gorm:"not null" visibility:"1"`
	ScriptId    int       `gorm:"not null" visibility:"1"`
	Script      TblScript `gorm:"foreignKey:ScriptId" visibility:"1"`
}

type TblInstanceType struct {
	Id               int       `gorm:"primaryKey" visibility:"1"`
	CreatedDate      time.Time `gorm:"autoCreateTime" visibility:"1"`
	UpdatedDate      time.Time `gorm:"autoUpdateTime" visibility:"1"`
	Name             string    `gorm:"size:255;not null" visibility:"1"`
	VcpuCount        int       `gorm:"not null" visibility:"1"`
	MemorySizeGb     int       `gorm:"not null" visibility:"1"`
	NetworkSpeedMbps int       `gorm:"not null" visibility:"1"`
	SpotPrice        float64   `gorm:"type:decimal(10,2);not null" visibility:"10"`
	SellPrice        float64   `gorm:"type:decimal(10,2);not null" visibility:"1"`
	Status           string    `gorm:"type:enum('ACTIVE', 'INACTIVE', 'HIDDEN');not null" visibility:"1"`
}

type TblInstance struct {
	Id             int                  `gorm:"primaryKey" visibility:"1"`
	CreatedDate    time.Time            `gorm:"autoCreateTime" visibility:"1"`
	UpdatedDate    time.Time            `gorm:"autoUpdateTime" visibility:"1"`
	Name           string               `gorm:"size:255;not null" visibility:"1"`
	AwsInstanceId  string               `gorm:"size:255;not null" visibility:"10"`
	PublicIp       string               `gorm:"size:255;not null" visibility:"2"`
	PrivateIp      string               `gorm:"size:255;not null" visibility:"10"`
	InstanceTypeId int                  `gorm:"not null" visibility:"1"`
	Status         types.InstanceStatus `gorm:"size:255;not null" visibility:"1"`
	InstanceType   TblInstanceType      `gorm:"foreignKey:InstanceTypeId" visibility:"1"`
}

type TblStorage struct {
	Id               int                 `gorm:"primaryKey"`
	CreatedDate      time.Time           `gorm:"autoCreateTime"`
	UpdatedDate      time.Time           `gorm:"autoUpdateTime"`
	Name             string              `gorm:"size:255;not null"`
	SizeMb           int                 `gorm:"not null"`
	AwsEbsId         string              `gorm:"size:255" visibility:"10"`
	AwsEbsSnapshotId string              `gorm:"size:255" visibility:"10"`
	Status           types.StorageStatus `gorm:"size:255;not null"`
}

type TblPersistentNode struct {
	Id               int                        `gorm:"primaryKey"`
	CreatedDate      time.Time                  `gorm:"autoCreateTime"`
	UpdatedDate      time.Time                  `gorm:"autoUpdateTime"`
	Name             string                     `gorm:"size:255;not null"`
	OwnerId          int                        `gorm:"not null"`
	ServerTemplateId int                        `gorm:"not null"`
	InstanceId       int                        `gorm:"default:null"`
	InstanceTypeId   int                        `gorm:"not null"`
	StorageId        int                        `gorm:"not null"`
	Status           types.PersistentNodeStatus `gorm:"size:255;not null"`
	Owner            TblUser                    `gorm:"foreignKey:OwnerId"`
	ServerTemplate   TblServerTemplate          `gorm:"foreignKey:ServerTemplateId"`
	Instance         TblInstance                `gorm:"foreignKey:InstanceId"`
	InstanceType     TblInstanceType            `gorm:"foreignKey:InstanceTypeId"`
	Storage          TblStorage                 `gorm:"foreignKey:StorageId"`
}

type TblParameter struct {
	Id          string    `gorm:"primaryKey"`
	Value       string    `gorm:"type:text;not null"`
	Desc        string    `gorm:"type:text;not null"`
	CreatedDate time.Time `gorm:"autoCreateTime"`
	UpdatedDate time.Time `gorm:"autoUpdateTime"`
}
