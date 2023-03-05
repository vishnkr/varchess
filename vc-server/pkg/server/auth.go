package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"unicode"
	"varchess/pkg/store"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Success  bool   `json:"success"`
}

func CreateJwtToken(userid int) (string, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userid
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Server) CreateAccountHandler(w http.ResponseWriter, r *http.Request) error {
	user := &store.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}
	if !validatePassword(user.Password) {
		return WriteJSON(w, http.StatusUnprocessableEntity, 
			ApiError { Error: "Invalid password"})

	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	user.Password = string(hashPassword)
	return s.store.CreateUser(user)

}

func (s *Server) AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) error {
	user := &store.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, AuthResponse{Success: false})
	}
	result, err := s.store.GetUserByUsername(user.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			//invalid username
			return WriteJSON(w, http.StatusUnauthorized, AuthResponse{Success: false})
		}
		return WriteJSON(w, http.StatusInternalServerError, AuthResponse{Success: false})
	}
	ok := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if ok != nil {
		return WriteJSON(w, http.StatusUnauthorized, AuthResponse{Success: false})
	} else {
		token, err := CreateJwtToken(result.ID)
		if err != nil {
			return WriteJSON(w, http.StatusUnprocessableEntity, AuthResponse{Success: false})
		}
		return WriteJSON(w, http.StatusOK, AuthResponse{UserID: result.ID, Username: result.Username, Token: token, Success: true})
	}
}

func displayResult(result *sql.Rows) {
	for result.Next() {
		var (
			userid   int
			password string
		)
		if err := result.Scan(&userid, &password); err != nil {
			panic(err)
		}
		fmt.Printf("%d is %s\n", userid, password)
	}
}

func validatePassword(password string) bool {
    hasUpper, hasLower, hasNumber, hasSpecial := false, false, false, false

    for _, ch := range password {
        switch {
        case unicode.IsUpper(ch):
            hasUpper = true
        case unicode.IsLower(ch):
            hasLower = true
        case unicode.IsNumber(ch):
            hasNumber = true
        case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
            hasSpecial = true
        }
    }
    return hasUpper && hasLower && hasNumber && hasSpecial && len(password) >= 7
}
