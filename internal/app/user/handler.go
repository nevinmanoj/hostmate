package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	. "github.com/nevinmanoj/hostmate/api"
	user "github.com/nevinmanoj/hostmate/internal/domain/user"
)

type UserHandler struct {
	service   user.UserService
	validator *validator.Validate
}

func NewUserHandler(s user.UserService) *UserHandler {
	return &UserHandler{service: s, validator: validator.New()}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetUser endpoint"))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req CreateUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(&req); err != nil {
		fmt.Print(err)
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid JSON body",
		})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	var email string = req.Email
	var password string = req.Password
	var name string = req.Name
	createdUser, err := h.service.CreateUser(ctx, email, password, name)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(PostResponsePage[user.User]{
		Message:    "User created successfully",
		Data:       *createdUser,
		StatusCode: http.StatusCreated,
	})
}
func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid JSON body1",
		})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}

	var email string = req.Email
	var password string = req.Password

	token, user, err := h.service.LoginUser(ctx, email, password)
	if err != nil {
		json.NewEncoder(w).Encode(ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message:    err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(LoginUserResponse{
		Token: token,
		User:  user,
	})
}
