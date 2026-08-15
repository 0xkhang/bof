package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/0xkhangle/bof/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type apiConfig struct {
	db       database.Client
	port     string
	host     string
	platform string
}

func (cfg *apiConfig) HandleHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
		<html>
		<head></head>
			<body>
				<h1>Hello world</h1>
			</body>
		</html>
	`))
}

func (cfg *apiConfig) HandleRollbackDB(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.Write([]byte(`You do not have the permisson. imposter!`))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := cfg.db.Rollback(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(`DB is clean!`))
	w.WriteHeader(http.StatusOK)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Couldn't load env file: %s", err)
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatalf("Missing DB path. DB path should be set")
	}

	dbClient, err := database.NewClient(dbPath)
	if err != nil {
		log.Fatalf("Couldnt't create DB client")
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	platform := os.Getenv("PLATFORM")
	if port == "" {
		port = "dev"
	}

	cfg := apiConfig{
		host:     host,
		port:     port,
		db:       dbClient,
		platform: platform,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", cfg.HandleHello)
	mux.HandleFunc("/admin/reset", cfg.HandleRollbackDB)

	s := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	log.Printf("Running server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
