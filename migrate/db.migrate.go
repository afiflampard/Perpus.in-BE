package migrate

import (
	"log"
	"onboarding/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	tableExist := (db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Role{}) && db.Migrator().HasTable(&models.Book{}) && db.Migrator().HasTable(&models.Borrow{}) && db.Migrator().HasTable(&models.OrderState{}) && db.Migrator().HasTable(&models.OrderDetail{}) && db.Migrator().HasTable(&models.History{}) && db.Migrator().HasTable(&models.Stock{}))

	if !tableExist {
		// dbMigrate := db.Debug().Migrator().DropTable(&models.User{}, &models.Role{}, &models.Book{}, &models.Borrow{}, &models.OrderState{}, &models.OrderDetail{}, &models.History{}, &models.Stock{})
		// if dbMigrate != nil {
		// 	log.Fatal("Cannot drop Table")
		// }
		db.AutoMigrate(&models.User{}, &models.Role{}, &models.Book{}, &models.Borrow{}, &models.OrderState{}, &models.OrderDetail{}, &models.History{}, &models.Stock{})

		var roles = []models.Role{
			models.Role{
				Role: "Petugas",
			},
			models.Role{
				Role: "Member",
			},
		}

		pass, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}
		var Users = []models.User{
			models.User{
				Email:    "afiflampard32@gmail.com",
				Name:     "Afif",
				Mobile:   "08576543434",
				Address:  "Jabon",
				Password: string(pass),
				RoleID:   1,
			},
			models.User{
				Email:    "afiflampard123@gmail.com",
				Name:     "Fifa",
				Mobile:   "08576543439",
				Address:  "Jabon",
				Password: string(pass),
				RoleID:   2,
			},
		}
		var orderstates = []models.OrderState{
			{
				No:   1,
				Name: "Dipinjam",
			},
			{
				No:   2,
				Name: "Dikembalikan",
			},
		}

		for _, role := range roles {
			err := db.Debug().Create(&role).Error
			if err != nil {
				log.Fatalf("Failed to create Role")
			}
		}

		for _, user := range Users {
			err := db.Debug().Create(&user).Error
			if err != nil {
				log.Fatalf("Failed to create User")
			}
		}
		for _, orderState := range orderstates {
			err := db.Debug().Create(&orderState).Error
			if err != nil {
				log.Fatal("Failed to Create OrderState")
			}
		}

	}

}
