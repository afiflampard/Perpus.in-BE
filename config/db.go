package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	pgName := os.Getenv("PGNAME")
	pgPassword := os.Getenv("PGPASSWORD")
	pgDB := os.Getenv("PGDATABASE")
	pgHost := os.Getenv("PGHOST")
	pgPort := os.Getenv("PGPORT")

	postgresConname := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%v", pgHost, pgName, pgDB, pgPassword, pgPort)
	fmt.Println("canname is\t\t", postgresConname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 0,             // Slow SQL threshold
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(postgresConname), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	return db
}
