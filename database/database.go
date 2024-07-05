package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/jameshw-dev01/user-api/spec"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type userDB struct {
	Username string `gorm:"primaryKey"`
	Hash     string
	Email    string
	Name     string
	Age      uint
}

type dbWrapper struct {
	DB *gorm.DB
}

func toUserDB(user spec.User) userDB {
	return userDB{
		Username: user.Username,
		Hash:     user.Hash,
		Email:    user.Email,
		Name:     user.Name,
		Age:      user.Age,
	}
}

func ToSpecUser(user userDB) spec.User {
	return spec.User{
		Username: user.Username,
		Hash:     user.Hash,
		Email:    user.Email,
		Name:     user.Name,
		Age:      user.Age,
	}
}

// Create implements spec.DbInterface.
func (d dbWrapper) Create(user spec.User) error {
	userDb := toUserDB(user)
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
func (d dbWrapper) Delete(user spec.User) error {
	userDb := userDB{Username: user.Username}
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
func (d dbWrapper) ReadAll() ([]spec.User, error) {
	var records []userDB
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

func (d dbWrapper) Read(username string) (spec.User, error) {
	var user userDB
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
func (d dbWrapper) Update(user spec.User) error {
	userDb := toUserDB(user)
	ret := d.DB.Model(&userDb).Updates(userDb)
	if ret.Error != nil {
		return ret.Error
	}
	if ret.RowsAffected != 1 {
		return errors.New("wrong number of rows affected")
	}
	return nil
}

func InitDB() {
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	db, err := sql.Open("mysql", "root:"+password+"@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS USERDB")
	if err != nil {
		panic(err)
	}
}

func GetDBConnection(resetDb bool) spec.DbInterface {
	InitDB()
	password := os.Getenv("MYSQL_ROOT_PASSWORD")

	dsn := "root:" + password + "@tcp(localhost:3306)/USERDB?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database")
	}
	err = db.AutoMigrate(&userDB{})
	if err != nil {
		log.Fatal("Failed to migrate User table")
	}
	if resetDb {
		db.Delete(&userDB{}, "1=1")
	}
	wrap := dbWrapper{DB: db}
	return wrap
}
