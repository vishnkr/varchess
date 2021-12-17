package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type Credentials struct {
	ID uint64 `json:"user_id" db:"user_id"`
	Token string `json:"token"`
	Password string	`json:"password" db:"password"`
	Username string `json:"username" db:"username"`
	IsValid bool `json:"valid"`
}

func initDB(){
	var err error
	openStr := fmt.Sprintf("user=%s dbname=%s password=%s host=localhost sslmode=disable",os.Getenv("DB_USER"),os.Getenv("DB_NAME"),os.Getenv("DB_PASSWORD"))
	db, err = sql.Open("postgres", openStr)
	if err != nil {
		panic(err)
	}
	
}

func CreateJwtToken(userid uint64)(string,error){
	var err error
	claims:=jwt.MapClaims{}
	claims["authorized"]=true
	claims["user_id"]=userid
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
  	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token,err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err!=nil{
		return "",err
	}
	return token,nil
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain") 
	w.Header().Set("Access-Control-Allow-Origin", "*")
	initDB()
	defer db.Close()
	creds:=&Credentials{}
	err:= json.NewDecoder(r.Body).Decode(creds)
	println(err,r,creds)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
	}
	hashPassword,_ := bcrypt.GenerateFromPassword([]byte(creds.Password),0)
	fmt.Println("Password",creds.Password)
	fmt.Println("Hash",string(hashPassword))
	if _, err := db.Query("insert into users(username,password) values ($1, $2)", creds.Username, string(hashPassword)); err != nil {
		return
	}
	
}

func AuthUserHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json") 
	w.Header().Set("Access-Control-Allow-Origin", "*")
	initDB()
	defer db.Close()
	creds := &Credentials{IsValid:false}
	err:= json.NewDecoder(r.Body).Decode(creds)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result := db.QueryRow("select user_id, password from users where username=$1",creds.Username)
	checkCreds := &Credentials{}
	err =result.Scan(&checkCreds.ID,&checkCreds.Password)
	if err!=nil{
		log.Fatal(err)
		if err == sql.ErrNoRows {
			//invalid username
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ok:=bcrypt.CompareHashAndPassword([]byte(checkCreds.Password),[]byte(creds.Password)) 
	if ok!=nil{
		creds.IsValid = false
	} else {
		creds.IsValid = true
		token, err := CreateJwtToken(checkCreds.ID)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		creds.ID=checkCreds.ID
		creds.Token = token
	}
	json.NewEncoder(w).Encode(creds)
}



func displayResult(result *sql.Rows){
	for result.Next() {
        var (
            userid int
            password  string
        )
        if err := result.Scan(&userid, &password); err != nil {
            panic(err)
        }
        fmt.Printf("%d is %s\n", userid, password)
    }
}
