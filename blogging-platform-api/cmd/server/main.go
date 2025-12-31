package main

import (
	"log"
	"net/http"

	// Sesuaikan path import ini dengan nama module di go.mod kamu
	"blogging-platform-api/internal/common"
	"blogging-platform-api/internal/post"

	_ "github.com/lib/pq"
)

func main() {
	// 1. Inisialisasi Koneksi Database
	// Diambil dari Internal/common/database.go yang sudah dibuat
	db, err := common.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal menghubungkan ke database: %v", err)
	}
	defer db.Close()

	// 2. Inisialisasi Layer (Wiring / Dependency Injection)
	// Database -> Repository -> Service -> Handler
	postRepo := post.NewRepository(db)
	postService := post.NewService(*postRepo)
	postHandler := post.NewHandler(postService)

	// 3. Setup Router Native (Fitur Go 1.22+)
	mux := http.NewServeMux()

	// Mapping Routes ke Handler
	mux.HandleFunc("POST /posts", postHandler.CreatePostHandler)       // Create post
	mux.HandleFunc("GET /posts", postHandler.GetPostsHandler)         // Get All & Filter/Search
	mux.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		postHandler.GetPostByIDHandler(w, r, r.PathValue("id"))
	}) // Get Single post
	mux.HandleFunc("PUT /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		postHandler.UpdatePostHandler(w, r, r.PathValue("id"))
	}) // Update post
	mux.HandleFunc("DELETE /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		postHandler.DeletePostHandler(w, r, r.PathValue("id"))
	}) // Delete post

	// 4. Konfigurasi Server
	port := ":8080"
	log.Printf("Server berjalan di http://localhost%s 🚀", port)

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	// Jalankan Server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server gagal berjalan: %v", err)
	}
}