package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"ai-recruit/internal/config"
	"ai-recruit/internal/db"
)

// basicAuth wraps an HTTP handler with Basic Authentication using provided credentials.
func basicAuth(h http.Handler, user, pass string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || u != user || p != pass {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

type TrackingInfo = db.TrackingInfo

func main() {
	rand.Seed(time.Now().UnixNano())

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "config.json"
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	if cfg.DatabaseURL == "" {
		log.Fatal("database_url not set in config")
	}

	repo, err := db.NewRepository(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer repo.Close()

	adminUser := cfg.AdminUsername
	adminPass := cfg.AdminPassword

	http.Handle("/admin/api/trackings", basicAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		infos, err := repo.List(r.Context())
		if err != nil {
			log.Printf("list trackings: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(infos)
	}), adminUser, adminPass))

	fs := http.FileServer(http.Dir("admin"))
	http.Handle("/admin/", basicAuth(http.StripPrefix("/admin/", fs), adminUser, adminPass))

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		campaign := r.URL.Query().Get("campaign")
		if email == "" || campaign == "" {
			http.Error(w, "email and campaign required", http.StatusBadRequest)
			return
		}
		info := &TrackingInfo{
			Email:     email,
			Campaign:  campaign,
			UTMSource: "email",
			UTMMedium: "newsletter",
			CreatedAt: time.Now(),
		}
		id, err := repo.Create(r.Context(), info)
		if err != nil {
			log.Printf("create track: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		url := "https://example.com/welcome?utm_source=" + info.UTMSource + "&utm_medium=" + info.UTMMedium + "&utm_campaign=" + campaign + "&uid=" + id
		resp := map[string]string{"url": url}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("uid")
		if id == "" {
			http.Error(w, "missing uid", http.StatusBadRequest)
			return
		}
		if err := repo.MarkClicked(r.Context(), id); err != nil {
			log.Printf("mark clicked: %v", err)
		}
		http.Redirect(w, r, "https://example.com/", http.StatusSeeOther)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
