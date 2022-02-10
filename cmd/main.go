package main

import (
	"fmt"
	"log"

	"github.com/yarikyarichek/streamer/config"
	"github.com/yarikyarichek/streamer/infostructure/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_DATABASE, config.DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewService(db)

	err = repo.Migrate()
	if err != nil {
		log.Fatal(err)
	}

}
