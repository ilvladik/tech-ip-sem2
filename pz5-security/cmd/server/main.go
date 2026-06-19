package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"example.com/pz5-security/internal/config"
	"example.com/pz5-security/internal/httpapi"
	"example.com/pz5-security/internal/httpserver"
	"example.com/pz5-security/internal/student"
)

func main() {
	cfg := config.New()

	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("database ping failed:", err)
	}

	repo := student.NewRepo(db)

	stmtByID, err := repo.PrepareGetByID()
	if err != nil {
		log.Fatal(err)
	}
	defer stmtByID.Close()

	stmtByEmail, err := repo.PrepareGetByEmail()
	if err != nil {
		log.Fatal(err)
	}
	defer stmtByEmail.Close()

	handler := httpapi.NewHandler(repo, stmtByID, stmtByEmail)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/students", handler.GetStudentByID)
	mux.HandleFunc("/students/by-email", handler.GetStudentByEmail)
	mux.HandleFunc("/students/unsafe", handler.GetStudentUnsafe)

	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal(err)
	}

	httpserver.StartRedirectServer(cfg.HTTPRedirectAddr, "https://localhost:8443")

	server := &http.Server{
		Addr:    cfg.HTTPSAddr,
		Handler: mux,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		},
	}

	log.Println("HTTPS server started on https://localhost:8443")
	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatal(err)
	}
}
