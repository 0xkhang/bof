package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Couldn't load env file: %s", err)
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`
			<html>
			<head></head>
				<body>
					<h1>Hello world</h1>
				</body>
			</html>
		`))
	})

	s := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	log.Printf("Running server on %s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
