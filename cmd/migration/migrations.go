package main

import (
	"main/internal/ds"
	"main/internal/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	env := dsn.FromEnv()

	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.AutoMigrate(
		&ds.User{},
		&ds.Operation{},
		&ds.Banknote{},
		&ds.OperationBanknote{},
	); err != nil {
		panic("cant migrate db:" + err.Error())
	}

	users := []ds.User{
		{Login: "user1", Password: "password1", IsModerator: true, Balance: 9999},
		{Login: "user2", Password: "password2", IsModerator: false, Balance: 1000},
		{Login: "user3", Password: "password3", IsModerator: false, Balance: 90},
	}

	money := []ds.Banknote{
		{Nominal: 10, ImageURL: "http://127.0.0.1:9000/banknote-name-server/206517eb-31e3-4ba2-9302-e8d7ad8a1dac.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 50, ImageURL: "http://127.0.0.1:9000/banknote-name-server/4c6be3eb-a199-4893-be97-7570aa97db6c.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 100, ImageURL: "http://127.0.0.1:9000/banknote-name-server/a370d56a-c798-41f1-8b94-8accc041e2bf.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 500, ImageURL: "http://127.0.0.1:9000/banknote-name-server/918daaaa-9a38-4127-ad52-7f2d1c0d7161.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 1000, ImageURL: "http://127.0.0.1:9000/banknote-name-server/d33d6d2e-3dc2-468d-8862-cc2a06b631a2.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 2000, ImageURL: "http://127.0.0.1:9000/banknote-name-server/08163a8b-9a6c-44a9-b101-dd2eb605c206.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 5000, ImageURL: "http://127.0.0.1:9000/banknote-name-server/bcdd454c-4ef3-4810-b05d-86730a5d8adf.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 1000, ImageURL: "http://127.0.0.1:9000/banknote-name-server/1000_yen.png", Description: "Японская купюра", Currency: "YEN", Status: "действует"},
		{Nominal: 100, ImageURL: "http://127.0.0.1:9000/banknote-name-server/b1712b8d-5f3d-48cf-bcc1-1e7bea91b4dd.jpg", Currency: "USD", Status: "действует"},
	}

	db.Create(&users)
	db.Create(&money)
}
