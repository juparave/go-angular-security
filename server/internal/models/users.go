package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           string `json:"id" gorm:"size:11"`
	AccountID    string `json:"accountId" gorm:"size:11;index"`
	RoleID       uint   `json:"roleId" gorm:"default:1"`
	FirstName    string `json:"firstName" gorm:"size:128"`
	LastName     string `json:"lastName" gorm:"size:128"`
	Email        string `json:"email" gorm:"size:128;unique"`
	Password     []byte `json:"-" gorm:"size:64"` // don't return password on json
	AccessToken  string `json:"accessToken" gorm:"-"`
	RefreshToken string `json:"refreshToken" gorm:"-"`

	// Security flags
	PasswordChangeRequired bool       `json:"passwordChangeRequired" gorm:"default:false"`
	PasswordDigest         string     `json:"-" gorm:"size:64"` // For detecting token reuse in password reset
	Enabled                bool       `json:"enabled" gorm:"default:true"`
	LastLoginAt            *time.Time `json:"lastLoginAt"`

	// Relationships
	Account Account `json:"account" gorm:"foreignKey:AccountID"`
	Role    Role    `json:"role" gorm:"foreignKey:RoleID"`

	// Legacy field for backward compatibility during migration
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

