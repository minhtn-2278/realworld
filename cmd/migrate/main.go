package main

import (
	"flag"
	"fmt"
	"os"

	"realworldapp/config"
	"realworldapp/migrations"
	database "realworldapp/pkg/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fail(err)
	}

	direction := flag.String("direction", "up", "migration direction: up or down")
	flag.Parse()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		fail(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		fail(err)
	}
	defer sqlDB.Close()

	switch *direction {
	case "up":
		err = migrations.Run(db)
	case "down":
		err = migrations.RollbackLast(db)
	default:
		fail(fmt.Errorf("unsupported migration direction %q; use up or down", *direction))
	}

	if err != nil {
		fail(err)
	}

	fmt.Printf("migration %s completed successfully\n", *direction)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
