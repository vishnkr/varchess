package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var db *sql.DB

var googleAuthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:5000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func OauthGoogleLogin(w http.ResponseWriter, r *http.Request){
	fmt.Println(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	fmt.Println("Sec",os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),googleAuthConfig)
	oauthState := genOauthCookie(w)
	u := googleAuthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func genOauthCookie(w http.ResponseWriter) string{
	var exp =time.Now().Add(365 * 24 * time.Hour)
	b:=make([]byte,16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: exp}
	http.SetCookie(w, &cookie)
	return state
}

func OauthGoogleCallback(w http.ResponseWriter, r *http.Request){
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value{
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func getUserDataFromGoogle(code string)([]byte,error){
	token, err := googleAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
type Credentials struct {
	ID       uint64 `json:"user_id" db:"user_id"`
	Token    string `json:"token"`
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
	IsValid  bool   `json:"valid"`
}

func initDB() {
	var err error
	openStr := fmt.Sprintf("user=%s dbname=%s password=%s host=localhost sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	db, err = sql.Open("postgres", openStr)
	if err != nil {
		panic(err)
	}

}

func CreateJwtToken(userid uint64) (string, error) {
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

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	initDB()
	defer db.Close()
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	println(err, r, creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(creds.Password), 0)
	fmt.Println("Password", creds.Password)
	fmt.Println("Hash", string(hashPassword))
	if _, err := db.Query("insert into users(username,password) values ($1, $2)", creds.Username, string(hashPassword)); err != nil {
		return
	}

}

func AuthUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	initDB()
	defer db.Close()
	creds := &Credentials{IsValid: false}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := db.QueryRow("select user_id, password from users where username=$1", creds.Username)
	checkCreds := &Credentials{}
	err = result.Scan(&checkCreds.ID, &checkCreds.Password)
	if err != nil {
		log.Fatal(err)
		if err == sql.ErrNoRows {
			//invalid username
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ok := bcrypt.CompareHashAndPassword([]byte(checkCreds.Password), []byte(creds.Password))
	if ok != nil {
		creds.IsValid = false
	} else {
		creds.IsValid = true
		token, err := CreateJwtToken(checkCreds.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		creds.ID = checkCreds.ID
		creds.Token = token
	}
	json.NewEncoder(w).Encode(creds)
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
