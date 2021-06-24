package models

import (
	"errors"
	"strings"
	"todo/api/utils"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(100);unique" json:"email"`
	Name         string `gorm:"type:varchar(100);not null" json:"name"`
	Password     string `gorm:"type:varchar(100); not null" json:"password"`
	IsRegistered bool   `gorm:"type:boolean;not null" json:"isRegistered"`
}

func (user *User) SaveUser(db *gorm.DB) error {
	err := db.Debug().Create(&user).Error
	return err
}

func GetUserByEmail(email string, db *gorm.DB) (*User, error) {
	foundRecord := &User{}
	err := db.Debug().Table("users").Where("email = ?", email).Find(foundRecord).Error
	if err != nil {
		return nil, err
	}
	return foundRecord, nil
}

func (user *User) UpdateUser(db *gorm.DB) error {
	return db.Debug().Table("users").Where("email = ?", user.Email).Update(&user).Error
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("incorrect password")
	}
	return nil
}

func (u *User) BeforeSave() error {
	password := strings.TrimSpace(u.Password)
	hashedpassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = string(hashedpassword)
	return nil
}

func (u *User) AfterSave() error {
	if !u.IsRegistered {
		return utils.SendEmail(u.Email)
	}
	return nil
}

func (user *User) PrepareData() {
	user.Email = strings.TrimSpace(user.Email)
	user.Name = strings.TrimSpace(user.Name)
	user.Password = strings.TrimSpace(user.Password)
}

func (user *User) ValidateFields(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if user.Email == "" {
			return errors.New("email is required")
		}
		if user.Password == "" {
			return errors.New("password is required")
		}
		return nil
	default:
		if user.Name == "" {
			return errors.New("name is required")
		}
		if user.Email == "" {
			return errors.New("email is required")
		}
		if user.Password == "" {
			return errors.New("password is required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("invalid email format")
		}
		return nil
	}
}
