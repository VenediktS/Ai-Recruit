package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type TrackingInfo struct {
	Email     string    `json:"email"`
	Campaign  string    `json:"campaign"`
	UTMSource string    `json:"utm_source"`
	UTMMedium string    `json:"utm_medium"`
	Clicked   bool      `json:"clicked"`
	CreatedAt time.Time `json:"created_at"`
	ClickedAt time.Time `json:"clicked_at"`
}

type Store struct {
	mu    sync.RWMutex
	items map[string]*TrackingInfo
}

func NewStore() *Store {
	return &Store{items: make(map[string]*TrackingInfo)}
}

func (s *Store) Create(info *TrackingInfo) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := generateID()
	s.items[id] = info
	return id
}

func (s *Store) MarkClicked(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if item, ok := s.items[id]; ok {
		item.Clicked = true
		item.ClickedAt = time.Now()
	}
}

func (s *Store) Get(id string) (*TrackingInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.items[id]
	return info, ok
}
func generateID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	store := NewStore()
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
		id := store.Create(info)
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
		store.MarkClicked(id)
		http.Redirect(w, r, "https://example.com/", http.StatusSeeOther)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
