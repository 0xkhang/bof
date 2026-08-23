package main

import (
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"

	"github.com/0xkhangle/bof/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wneessen/go-mail"
)

type apiConfig struct {
	db            database.Client
	port          string
	host          string
	platform      string
	smtp_user     string
	smtp_password string
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

	stmpUser := os.Getenv("SMTP_USER")
	stmpPass := os.Getenv("SMTP_PASSWORD")

	cfg := apiConfig{
		host:          host,
		port:          port,
		db:            dbClient,
		platform:      platform,
		smtp_user:     stmpUser,
		smtp_password: stmpPass,
	}

	v1 := http.NewServeMux()
	v1.HandleFunc("GET /hello", cfg.HandlerHello)
	v1.HandleFunc("POST /admin/reset", cfg.HandlerResetDB)
	v1.HandleFunc("POST /user", cfg.HandlerSignUp)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	s := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	message := mail.NewMsg()
	if err := message.From("lenguyenkhang102002@gmail.com"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := message.To("ngukhangle@gmail.com"); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}
	message.Subject("This is my first mail with go-mail!")
	message.SetBodyString(mail.TypeTextPlain, "Do you like this mail? I certainly do!")

	client, err := mail.NewClient("smtp.gmail.com", mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername("lenguyenkhang102002"), mail.WithPassword("csev ixaa vydh pxsn"))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	if err := client.DialAndSend(message); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}

	log.Printf("Running server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
