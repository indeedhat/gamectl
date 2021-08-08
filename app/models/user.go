package models

import "time"

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
