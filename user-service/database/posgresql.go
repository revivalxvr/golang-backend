package database

import (
	"user-service/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresql() *gorm.DB {
	gormDB, err := gorm.Open(postgres.Open(env.Conf.Postgres.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return gormDB
}

func ClosePostgresql(gormDB *gorm.DB) {
	db, err := gormDB.DB()
	if err != nil {
		panic("failed to close database")
	}
	db.Close()
}
