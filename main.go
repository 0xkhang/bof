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
	db   database.Client
	port string
	host string
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

	cfg := apiConfig{
		host: host,
		port: port,
		db:   dbClient,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", cfg.HandleHello)

	s := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	log.Printf("Running server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
