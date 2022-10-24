package authentication

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var DbConnectionString string

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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type UserGorm struct {
	Db *gorm.DB
}

// CreateUser adds user in db
func (ug *UserGorm) CreateUser(user *User) (*User, error) {
	result := ug.Db.Create(&user)
	if result.Error != nil {
		errorLog.Printf("Could not add '%s' user in database.", user.FirstName)
		return nil, result.Error
	}
	infoLog.Printf("Added '%s' user in database.", user.FirstName)
	return user, nil
}

// Connect establishes connection to mysql
func (ug *UserGorm) Connect(DbConnectionString string) bool {
	var err error
	ug.Db, err = gorm.Open(mysql.Open(DbConnectionString), &gorm.Config{})
	if err != nil {
		errorLog.Fatal("Failed to connect database", err)
		return false
	}
	infoLog.Println("Successfully connected to db.")
	return true
}

// InitialMigration creates the table if it doesn't exist
func (ug *UserGorm) InitialMigration() {
	ug.Db.AutoMigrate(&User{})
}

// ByEmail looks up user by Email address
func (ug *UserGorm) ByEmail(email string) (*User, error) {
	var user User
	result := ug.Db.First(&user, email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
