package main

import (
	"flag"
	"log"

	"github.com/iPatrushevSergey/adpace/app/internal/pkg/migrate"
)

func main() {
	dsn := flag.String("d", "", "database DSN")
	dir := flag.String("dir", "", "path to migration files")
	flag.Parse()

	if *dsn == "" {
		log.Fatal("migrate: -d is required")
	}
	if *dir == "" {
		log.Fatal("migrate: -dir is required")
	}

	if err := migrate.PostgresUp(*dsn, *dir); err != nil {
		log.Fatalf("migrate: %v", err)
	}
}
