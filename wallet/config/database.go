package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Conn   *gorm.DB
	ConnSy *gorm.DB
	Err    error
	ErrSy  error
}

var database = Database{}

func ConnectDB() (c *gorm.DB, err error) {
	DB_CONNECTION := os.Getenv("DB_CONNECTION")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_DATABASE := os.Getenv("DB_DATABASE")
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	DB_TEST := os.Getenv("DB_TEST")
	DB_DETAIL := DB_USERNAME + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_DATABASE + "?parseTime=true"
	if DB_CONNECTION == "" {
		DB_DETAIL = DB_TEST
		conn, err := gorm.Open(sqlite.Open("./../../test.db"), &gorm.Config{})
		if err != nil || conn == nil {
			fmt.Println("Error connecting to DB")
			fmt.Println(err.Error())
		}
		return conn, err
	} else {
		if database.Conn == nil {
			database.Conn, database.Err = gorm.Open(mysql.Open(DB_DETAIL), &gorm.Config{})
			if database.Err != nil || database.Conn == nil {
				fmt.Println("Error connecting to DB")
				fmt.Println(database.Err.Error())
			}
		}
		return database.Conn, database.Err
	}
}

func ConnectDBSY() (c *gorm.DB, err error) {
	DB_CONNECTION := os.Getenv("DB_CONNECTION_SY")
	DB_HOST := os.Getenv("DB_HOST_SY")
	DB_PORT := os.Getenv("DB_PORT_SY")
	DB_DATABASE := os.Getenv("DB_DATABASE_SY")
	DB_USERNAME := os.Getenv("DB_USERNAME_SY")
	DB_PASSWORD := os.Getenv("DB_PASSWORD_SY")

	DB_TEST := os.Getenv("DB_TEST")
	DB_DETAIL := DB_USERNAME + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_DATABASE + "?parseTime=true"
	if DB_CONNECTION == "" {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(pwd)
		DB_DETAIL = DB_TEST
		conn, err := gorm.Open(sqlite.Open("./../../testSy.db"), &gorm.Config{})
		if err != nil || conn == nil {
			fmt.Println("Error connecting to DB")
			fmt.Println(err.Error())
		}
		return conn, err
	} else {
		if database.ConnSy == nil {
			database.ConnSy, database.ErrSy = gorm.Open(mysql.Open(DB_DETAIL), &gorm.Config{})
			if database.ErrSy != nil || database.ConnSy == nil {
				fmt.Println("Error connecting to DB")
				fmt.Println(database.Err.Error())
			}
		}
		return database.ConnSy, database.ErrSy
	}
}
