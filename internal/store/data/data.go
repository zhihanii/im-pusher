package data

import "gorm.io/gorm"

type Data struct {
	db *gorm.DB
}

func NewData(db *gorm.DB) *Data {
	return &Data{
		db: db,
	}
}
