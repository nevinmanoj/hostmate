package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	. "github.com/nevinmanoj/hostmate/api"
	"github.com/nevinmanoj/hostmate/internal/app/httputil"
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
	ctx := r.Context()
	userIdStr := chi.URLParam(r, "userId")
	log.Println("HandlerGetUser::Fetching user with ID:", userIdStr)
	w.Header().Set("Content-Type", "application/json")
	var resp any
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		resp = httputil.GetErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}
	result, err := h.service.GetUserByID(ctx, userId)
	if err != nil {
		resp = httputil.GetErrorResponse(err)
	} else {
		userResponse := ToUserResponse(result)
		resp = GetResponsePage[UserResponse]{
			StatusCode: 200,
			Message:    "Payment fetched successfully",
			Data:       userResponse,
		}
	}

	json.NewEncoder(w).Encode(resp)
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
	userResponse := ToUserResponse(createdUser)
	json.NewEncoder(w).Encode(PostResponsePage[UserResponse]{
		Message:    "User created successfully",
		Data:       userResponse,
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
		resp := httputil.GetErrorResponse(err)
		json.NewEncoder(w).Encode(resp)
		return
	}
	logingResponse := ToLoginUserResponse(user, token)

	json.NewEncoder(w).Encode(PostResponsePage[LoginUserResponse]{
		Message:    "User created successfully",
		Data:       logingResponse,
		StatusCode: http.StatusCreated,
	})
}
