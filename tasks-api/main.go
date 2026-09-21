package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func mustEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	dsn := mustEnv("DATABASE_URL", "postgres://app:secret@localhost:5432/app?sslmode=disable")
	addr := ":" + mustEnv("PORT", "8080")

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	store := &TaskStore{db: db}
	h := &Handler{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", h.list)
	mux.HandleFunc("POST /tasks", h.create)
	mux.HandleFunc("GET /tasks/{id}", h.getOne)
	mux.HandleFunc("PATCH /tasks/{id}", h.update)
	mux.HandleFunc("DELETE /tasks/{id}", h.remove)

	srv := &http.Server{
		Addr:         addr,
		Handler:      withLogging(withCORS(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("server listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
