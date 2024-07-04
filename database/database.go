package db

import (
	"errors"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/jameshw-dev01/user-api/spec"
	"gorm.io/gorm"
)

type UserDB struct {
	Username string `gorm:"primaryKey"`
	Hash     string
	Email    string
	Name     string
	Age      uint
}

type DBWrapper struct {
	DB *gorm.DB
}

func ToUserDB(user spec.User) UserDB {
	return UserDB{
		Username: user.Username,
		Hash:     user.Hash,
		Email:    user.Email,
		Name:     user.Name,
		Age:      user.Age,
	}
}

func ToSpecUser(user UserDB) spec.User {
	return spec.User{
		Username: user.Username,
		Hash:     user.Hash,
		Email:    user.Email,
		Name:     user.Name,
		Age:      user.Age,
	}
}

// Create implements spec.DbInterface.
func (d DBWrapper) Create(user spec.User) error {
	userDb := ToUserDB(user)
	ret := d.DB.Create(&userDb)
	if ret.Error != nil {
		return ret.Error
	}
	if ret.RowsAffected != 1 {
		return errors.New("wrong number of rows affected")
	}
	return nil
}

// Delete implements spec.DbInterface.
func (d DBWrapper) Delete(user spec.User) error {
	userDb := UserDB{Username: user.Username}
	ret := d.DB.Delete(&userDb)
	if ret.Error != nil {
		return ret.Error
	}
	if ret.RowsAffected != 1 {
		return errors.New("wrong number of rows affected")
	}
	return nil
}

// ReadAll implements spec.DbInterface.
func (d DBWrapper) ReadAll() ([]spec.User, error) {
	var records []UserDB
	ret := d.DB.Find(&records)
	if ret.Error != nil {
		return []spec.User{}, ret.Error
	}
	var specUsers []spec.User
	for _, u := range records {
		specUsers = append(specUsers, ToSpecUser(u))
	}
	return specUsers, nil
}

func (d DBWrapper) Read(username string) (spec.User, error) {
	var user UserDB
	user.Username = username
	ret := d.DB.First(&user)
	if ret.Error != nil {
		return spec.User{}, ret.Error
	}
	if ret.RowsAffected != 1 {
		return spec.User{}, errors.New("wrong number of rows found")
	}
	return ToSpecUser(user), nil
}

// Update implements spec.DbInterface.
func (d DBWrapper) Update(user spec.User) error {
	userDb := ToUserDB(user)
	ret := d.DB.Model(&userDb).Updates(userDb)
	if ret.Error != nil {
		return ret.Error
	}
	if ret.RowsAffected != 1 {
		return errors.New("wrong number of rows affected")
	}
	return nil
}

func initDB() spec.DbInterface {
	os.Remove("gorm.db")
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	err = db.AutoMigrate(&UserDB{})
	if err != nil {
		log.Fatal("Failed to migrate User table")
	}
	wrap := DBWrapper{DB: db}
	return wrap
}
