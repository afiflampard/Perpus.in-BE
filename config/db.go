package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db, err := gorm.Open(postgres.Open(postgresConname), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
