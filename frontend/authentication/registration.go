package authentication

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

type User struct {
	gorm.Model
	FirstName      string
	LastName       string
	Email          string            `gorm:"not null;unique"`
	Password       string            `gorm:"-"` //ignored by db
	HashedPassword string            `gorm:"not null"`
	Errors         map[string]string `gorm:"-"` //ignored by db
}

type Users []*User

type UserModel struct {
	DB *gorm.DB
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// CreateUser adds user in db
func (ug *UserModel) CreateUser(user *User) (*User, error) {
	result := ug.DB.Create(&user)
	if result.Error != nil {
		errorLog.Printf("Could not add '%s' user in database.", user.FirstName)
		return nil, result.Error
	}
	infoLog.Printf("Added '%s' user in database.", user.FirstName)
	return user, nil
}

// InitialMigration creates the table if it doesn't exist
func (ug *UserModel) InitialMigration() {
	ug.DB.AutoMigrate(&User{})
}

// ByEmail looks up user by Email address
func (ug *UserModel) ByEmail(email string) (*User, error) {
	var user User
	result := ug.DB.First(&user, "Email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
