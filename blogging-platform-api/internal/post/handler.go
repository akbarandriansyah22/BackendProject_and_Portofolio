package post

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type handler struct {
	service *Service
}

func NewHandler(s *Service) *handler {
	return &handler{service: s}
}

// helper untuk mengirim response JSON
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
} 

// helper untuk mengirim response error
func respondError(w http.ResponseWriter, code int, message string){
	respondJSON(w, code, map[string]string{"error": message})
}

//createposthadler: post /posts
func (h *handler) CreatePostHandler(w http.ResponseWriter, r *http.Request){
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	post, err := h.service.CreatePost(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusCreated, post)
}

// getpostshandler: get /posts (bisa handle search juga via query param ?q=...)
func (h *handler) GetPostsHandler(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("q")

	var posts []*Post
	var err error
	
	if query != "" {
		posts, err = h.service.SearchPosts(r.Context(), query)
	} else {
		posts, err = h.service.GetAllPosts(r.Context())
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, posts)
}

//Getpostbyidhandler: get /posts/{id}
func (h *handler) GetPostByIDHandler(w http.ResponseWriter, r *http.Request, idParam string){
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	post, err := h.service.GetPostByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if post == nil {
		respondError(w, http.StatusNotFound, "Post not found")
		return
	}
	respondJSON(w, http.StatusOK, post)
}

// updateposthandler: put /posts/{id}
func (h *handler) UpdatePostHandler(w http.ResponseWriter, r *http.Request, idParam string){
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	var req updatePostInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	updatedPost, err := h.service.UpdatePost(r.Context(), id, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, updatedPost)
}

// deleteposthandler: delete /posts/{id}
func (h *handler) DeletePostHandler(w http.ResponseWriter, r *http.Request, idParam string){
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid post ID")
		return
	}
	if err := h.service.DeletePost(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}
