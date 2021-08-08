package models

import (
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Migrate() error {
	err := DB.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return initRootUser()
}

// Connect to the database
//
// If the database does not exist it will be created
func Connect() error {
	var err error

	if DB != nil {
		return nil
	}

	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	return err
}

// initRootUser will create the root user based on the .env details if one does not already exist
//
// this dows not make use of the CreateUser helper by design as that is incapable of creating a root user
func initRootUser() error {
	var user User

	DB.Where("root = 1").Find(&user)
	if user.ID != 0 {
		return nil
	}

	if os.Getenv("ROOT_USER") == "" || os.Getenv("ROOT_PASS") == "" {
		return errors.New("ROOT_USER or ROOT_PASS is not defined")
	}

	hash, err := hashPassword(os.Getenv("ROOT_PASS"))
	if err != nil {
		return err
	}

	user.Name = os.Getenv("ROOT_USER")
	user.Password = hash
	user.Root = true

	tx := DB.Create(&user)
	return tx.Error
}

// hashPassword will create a salted hash using bcrypt
func hashPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
