package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           string `json:"id" gorm:"size:11"`
	FirstName    string `json:"firstName" gorm:"size:128"`
	LastName     string `json:"lastName" gorm:"size:128"`
	Email        string `json:"email" gorm:"size:128; unique"`
	Password     []byte `json:"-" gorm:"size:64"` // don't return password on json
	AccessToken  string `json:"accessToken" gorm:"-"`
	RefreshToken string `json:"refreshToken" gorm:"-"`

	Roles []Role `json:"roles" gorm:"many2many:user_roles;"`
}

// UserWithPassword used on change password requests
type UserWithPassword struct {
	*User
	Password string `json:"password"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// UUID short version
	if user.ID == "" {
		user.ID = SortableShortUUID()
	}

	return nil
}

func (user *User) SetPassword(password string) {
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = hashPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

func (user *User) Count(db *gorm.DB) int64 {
	var total int64
	db.Model(&User{}).Count(&total)

	return total
}

func (user *User) Take(db *gorm.DB, limit int, offset int) interface{} {
	var users []User

	db.Preload("Role").Offset(offset).Limit(limit).Find(&users)

	return users
}
