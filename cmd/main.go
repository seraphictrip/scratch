package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"scratch/api"
	"scratch/config"
	"scratch/db"

	"github.com/go-sql-driver/mysql"
)

func main() {

	cfg := config.Envs

	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPassword,
		Addr:                 cfg.DBAddress,
		DBName:               cfg.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		log.Fatal("failed to start db", err)
	}

	err = initStorage(db)
	if err != nil {
		log.Fatal("failed to start db", err)
	}

	ctx := context.Background()
	addr := ":" + cfg.Port
	fmt.Println(addr)
	server := api.NewAPIServer(addr, db)

	err = server.Run(ctx)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func initStorage(db *sql.DB) error {
	err := db.Ping()

	return err
}
