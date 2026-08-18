package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/0xkhangle/bof/internal/auth"
	"github.com/0xkhangle/bof/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type apiConfig struct {
	db       database.Client
	hasher   auth.Hasher
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

func (cfg *apiConfig) HandleResetDB(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`You do not have the permisson. imposter!`))
		return
	}

	if err := cfg.db.Reset(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`DB is clean!`))
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
		log.Fatalf("Couldnt't create DB client: %s", err)
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
	if platform == "" {
		platform = "dev"
	}

	var hasher auth.Hasher
	switch os.Getenv("PASSWORD_ALGO") {
	case "argon2id":
		hasher = auth.NewArgon2Hasher()
	default:
		hasher = auth.NewBcryptHasher(bcrypt.DefaultCost)
	}

	cfg := apiConfig{
		host:     host,
		port:     port,
		db:       dbClient,
		hasher:   hasher,
		platform: platform,
	}

	v1 := http.NewServeMux()
	v1.HandleFunc("GET /hello", cfg.HandleHello)
	v1.HandleFunc("POST /admin/reset", cfg.HandleResetDB)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	s := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	log.Printf("Running server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
