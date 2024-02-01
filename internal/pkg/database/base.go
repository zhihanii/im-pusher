package database

import "time"

type Base struct {
	Id          uint64    `json:"id" gorm:"primaryKey"`
	IsDeleted   uint8     `json:"is_deleted" gorm:"column:is_deleted"`
	GMTCreate   time.Time `json:"gmt_create" gorm:"column:gmt_create"`
	GMTModified time.Time `json:"gmt_modified" gorm:"column:gmt_modified"`
}
