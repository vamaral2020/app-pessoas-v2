package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"app-pessoas/internal/db"
	"app-pessoas/internal/httpapi"
	"app-pessoas/internal/repo"
)

func main() {

	dsn := getenv("DB_DSN", "postgres://app:app@localhost:5432/appdb?sslmode=disable")
	port := getenv("PORT", "8080")

	database, err := db.Open(dsn)
	if err != nil {
		log.Fatal("Erro conectando ao banco", err)
	}
	defer database.Close()

	pessoaRepo := repo.NewPessoaRepo(database)
	handler := httpapi.NewHandler(pessoaRepo)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	srv := httpapi.LoggingMiddleware(mux)

	log.Println("API rodando em http://localhost:", port)
	log.Fatal(http.ListenAndServe(":"+port, srv))

}

func getenv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}
