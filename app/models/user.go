package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name     string
	Password string
	Root     bool
}

// FindUser by its id
func FindUser(id string) *User {
	var user User

	tx := DB.First(&user, id)
	if tx.Error != nil {
		return nil
	}

	return &user
}

// ListUsers will return all the users
func ListUsers() *[]User {
	var users []User

	tx := DB.Find(&users)
	if tx.Error != nil {
		return nil
	}

	return &users
}

// CreateUser is a simple helper function for adding a standard user to the database
func CreateUser(name, password string) *User {
	hash, err := hashPassword(password)
	if err != nil {
		return nil
	}

	user := User{
		Name:     name,
		Password: hash,
		Root:     false,
	}

	tx := DB.Create(&user)
	if tx.Error != nil {
		return nil
	}

	return &user
}

// UpdateUser is a helper function to make deleting users a little more readable
func UpdateUser(user *User) error {
	tx := DB.Save(user)

	return tx.Error
}

// DeleteUser is a helper function to make deleting users a little more readable
func DeleteUser(user *User) error {
	tx := DB.Delete(user)

	return tx.Error
}

// LoadUserByLoginDetails will check off the given user/pass and attempt to load a valid user
// from them
func LoadUserByLoginDetails(username, passwd string) *User {
	var user User

	tx := DB.Where("name = ?", username).First(&user)
	if tx.Error != nil {
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd)); err != nil {
		return nil
	}

	return &user
}
