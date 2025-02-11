package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/varnit-ta/Ecom-API/cmd/api"
	"github.com/varnit-ta/Ecom-API/config"
	"github.com/varnit-ta/Ecom-API/db"
)

/*
Main Function
*/
func main() {
	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":4000", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database is successfully connected!")
}
