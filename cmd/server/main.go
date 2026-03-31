package main

import (
	"context"
	_ "embed"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Butters19/url-shortener/internal/handler"
	"github.com/Butters19/url-shortener/internal/service"
	"github.com/Butters19/url-shortener/internal/storage"
	"github.com/Butters19/url-shortener/internal/storage/memory"
	"github.com/Butters19/url-shortener/internal/storage/postgres"
)

//go:embed migrations/001_unit.sql
var migration string

func main() {
	storageType := flag.String("storage", "memory", "storage type: memory or postgres")
	dsn := flag.String("dsn", "", "postgres DSN")
	addr := flag.String("addr", ":8080", "http server address")
	flag.Parse()

	var store storage.Storage

	switch *storageType {
	case "postgres":
		if *dsn == "" {
			log.Fatal("-dsn flag is required for postgres storage")
		}
		pg, err := postgres.New(context.Background(), *dsn)
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}
		defer pg.Close()

		if _, err := pg.Pool().Exec(context.Background(), migration); err != nil {
			log.Fatalf("failed to run migration: %v", err)
		}

		store = pg

	case "memory":
		store = memory.New()

	default:
		log.Fatalf("unknown storage type: %s", *storageType)
	}

	svc := service.New(store)
	h := handler.New(svc)

	log.Printf("starting server on %s with %s storage", *addr, *storageType)
	if err := http.ListenAndServe(*addr, h.Routes()); err != nil {
		log.Fatalf("server error: %v", err)
		os.Exit(1)
	}
}