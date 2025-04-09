package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"app/internal/dto"
	"app/internal/entity"
	"app/internal/infra/database"

	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB       database.UserInterface
	Jwt          *jwtauth.JWTAuth
	JwtExpiresIn int
}

func NewUserHandler(
	UserDB database.UserInterface,
	Jwt *jwtauth.JWTAuth,
	JwtExpiresIn int,
) *UserHandler {
	return &UserHandler{
		UserDB,
		Jwt,
		JwtExpiresIn,
	}
}

// GetJwt godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body     dto.GetJwtInput  true  "user credentials"
// @Success      200  {object}  dto.GetJwtOutput
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/v1/users/generate-jwt [post]
func (h *UserHandler) GetJwt(w http.ResponseWriter, r *http.Request) {
	var jwtInput dto.GetJwtInput
	err := json.NewDecoder(r.Body).Decode(&jwtInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid request body"}

		json.NewEncoder(w).Encode(err)
		return
	}

	u, err := h.UserDB.FindByEmail(jwtInput.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		err := Error{Message: "Invalid email or password"}
		json.NewEncoder(w).Encode(err)
		return
	}

	if !u.ValidatePassword(jwtInput.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		err := Error{Message: "Invalid email or password"}
		json.NewEncoder(w).Encode(err)
		return
	}

	_, tokenString, err := h.Jwt.Encode(map[string]interface{}{
		"sub": u.ID,
		"exp": time.Now().Add(time.Second * time.Duration(h.JwtExpiresIn)).Unix(),
	})

	accessToken := dto.GetJwtOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /api/v1/users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: "Invalid user data"}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: "Error creating user"}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
