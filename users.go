package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errEmailRequired = errors.New("email is required")
var errPasswordRequired = errors.New("password is required")

type UserService struct {
	store Store
}

func NewUserService(store Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", s.handleUserRegister).Methods("POST")
}

func (s *UserService) handleUserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error reading request payload"})
		return
	}
	defer r.Body.Close()

	// create a user
	var user *User
	err = json.Unmarshal(body, &user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid Request payload"})
		return
	}

	if err := validateUserPayload(user); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error hashing password"})
		return
	}

	user.Password = hashedPassword

	// save the user
	u, err := s.store.CreateUser(user)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating user"})
		return
	}

	// return the user
	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating session"})
		return
	}

	WriteJSON(w, http.StatusCreated, token)
}

func validateUserPayload(u *User) error {
	if u.Email == "" {
		return errEmailRequired
	}

	if u.Password == "" {
		return errPasswordRequired
	}

	return nil
}

func createAndSetAuthCookie(userID string, w http.ResponseWriter) (string, error) {
	secret := []byte(Envs.JWTSecret)
	token, err := CreateJWT(secret, userID)

	if err != nil {
		return "", err
	}

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    token,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	return token, nil
}
