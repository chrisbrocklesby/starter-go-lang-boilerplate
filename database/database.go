package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        string `gorm:"type:uuid;default:uuid7()"`
	Email     string
	Password  string
	FirstName *string
	LastName  *string
	Verified  bool   `gorm:"default:false"`
	Code      string `gorm:"default:null"`
}

type Posts struct {
	gorm.Model
	ID    string `gorm:"type:uuid;default:uuid7()"`
	Title string
	Body  string
}

var DB *gorm.DB

func init() {
	db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_URI")+"/test"), &gorm.Config{})
	if err != nil {
		fmt.Println("Database: Failed to connect database.")
	}

	fmt.Println("Database: Connection Successful.")

	DB = db

	// Migrate the schema
	db.AutoMigrate(&Posts{})
	db.AutoMigrate(&User{})
}
