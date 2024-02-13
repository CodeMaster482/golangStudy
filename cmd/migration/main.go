package main

import (
	"main/internal/app/ds"
	"main/internal/app/dsn"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	env, err := dsn.FromEnv()
	if err != nil {
		panic("Error from reading env")
	}

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
		{
			Login:    "user1",
			Name:     "name1",
			Password: "e38ad214943daad1d64c102faec29de4afe9da3d", // password1
			Role:     ds.Admin,
			Balance:  9999,
		},
		{
			Login:    "user2",
			Name:     "name2",
			Password: "2aa60a8ff7fcd473d321e0146afd9e26df395147", // password2
			Role:     ds.Buyer,
			Balance:  1000,
		},
		{
			Login:    "user3",
			Name:     "name3",
			Password: "1119cfd37ee247357e034a08d844eea25f6fd20f", // password3
			Role:     ds.Moderator,
			Balance:  90,
		},
	}

	money := []ds.Banknote{
		{Nominal: 10, ImageURL: "http://127.0.0.1:9000/banknote-name-server/9b92bb9e-a016-4193-b224-ae9c09f69970.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 50, ImageURL: "http://127.0.0.1:9000/banknote-name-server/4c6be3eb-a199-4893-be97-7570aa97db6c.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 100, ImageURL: "http://127.0.0.1:9000/banknote-name-server/a370d56a-c798-41f1-8b94-8accc041e2bf.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 500, ImageURL: "http://127.0.0.1:9000/banknote-name-server/918daaaa-9a38-4127-ad52-7f2d1c0d7161.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 1000, ImageURL: "http://127.0.0.1:9000/banknote-name-server/d33d6d2e-3dc2-468d-8862-cc2a06b631a2.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 2000, ImageURL: "http://127.0.0.1:9000/banknote-name-server/08163a8b-9a6c-44a9-b101-dd2eb605c206.jpg", Currency: "RUB", Status: "действует"},
		{Nominal: 5000, ImageURL: "http://127.0.0.1:9000/banknote-name-server/bcdd454c-4ef3-4810-b05d-86730a5d8adf.jpg", Currency: "RUB", Status: "действует"},
	}
	db.Create(&users)
	db.Create(&money)
}
