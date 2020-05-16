package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gavinc95/go-blog/db/models"
)

var (
	ErrBadRequest = fmt.Errorf("Invalid request: missing required parameters")
)

type GetUserRequest struct {
	ID string `json:"id"` // required
}

type GetUserResponse struct {
	User *models.User `json:"user"`
}

type CreateUserRequest struct {
	Email string `json:"email"` // required
	Name  string `json:"name"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

type UpdateUserRequest struct {
	ID    string `json:"id"` // required
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UpdateUserResponse struct {
	ID string `json:"id"`
}

type DeleteUserRequest struct {
	ID string `json:"id"` // required
}

type DeleteUserResponse struct {
	ID string `json:"id"`
}

type CreatePostRequest struct {
	UserID  string `json:"user_id"` // required
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CreatePostResponse struct {
	ID string `json:"id"`
}

type UpdatePostRequest struct {
	ID      string `json:"id"` // required
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostResponse struct {
	ID string `json:"id"`
}

type GetPostRequest struct {
	ID string `json:"id"` // required
}

type GetPostResponse struct {
	Post *models.Post `json:"post"`
}

type GetAllPostsRequest struct {
	UserID string `json:"user_id"` // required
}

type GetAllPostsResponse struct {
	Posts []*models.Post `json:"posts"`
}

type DeletePostRequest struct {
	ID string `json:"id"` // required
}

type DeletePostResponse struct {
	ID string `json:"id"`
}

func (a *App) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	var req GetUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.ID == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.BlogStore.GetUser(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := GetUserResponse{user}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.Email == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	userID, err := a.BlogStore.CreateUser(req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := CreateUserResponse{ID: userID}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.ID == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	userID, err := a.BlogStore.UpdateUser(req.ID, req.Name, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := UpdateUserResponse{ID: userID}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	var req DeleteUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.ID == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := a.BlogStore.DeleteUser(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := DeleteUserResponse{ID: id}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleGetAllPosts(w http.ResponseWriter, r *http.Request) {
	var req GetAllPostsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.UserID == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	posts, err := a.BlogStore.GetAllPosts(req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := GetAllPostsResponse{Posts: posts}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	var req GetPostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.ID == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	post, err := a.BlogStore.GetPost(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := GetPostResponse{Post: post}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.UserID == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	postID, err := a.BlogStore.CreatePost(req.UserID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := CreatePostResponse{ID: postID}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleUpdatePost(w http.ResponseWriter, r *http.Request) {
	var req UpdatePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.ID == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	postID, err := a.BlogStore.UpdatePost(req.ID, req.Title, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := UpdatePostResponse{ID: postID}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	var req DeletePostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate the request
	if req.ID == "" {
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	postID, err := a.BlogStore.DeletePost(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := DeletePostResponse{ID: postID}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
