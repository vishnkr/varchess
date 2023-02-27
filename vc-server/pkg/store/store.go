package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

type PostgresStore struct {
    db *sql.DB
}

func NewStore() (*PostgresStore,error) {
    var err error
	openStr := fmt.Sprintf("user=%s dbname=%s password=%s host=localhost port=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"),os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", openStr)
	if err != nil {
		return nil,err
	} 
	if err := db.Ping(); err != nil {
		return nil, err
	}
	wd, _ := os.Getwd()
    migrationsDir := filepath.Join(wd, "migrations")
	err = goose.Run("up", db, migrationsDir)
	if err != nil {
        log.Fatalf("failed to apply migrations: %s", err)
    }

	db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(time.Hour)
    return &PostgresStore{db: db},nil
}

type Storage interface{
	//User
	CreateUser(*User) error
	GetUserByUsername(string) (*User,error)
	DeleteUser(int) error
}

func (s *PostgresStore) CreateUser(user *User) error{
	stmt, err:= s.db.Prepare("INSERT INTO users(username,email,password) VALUES ($1, $2, $3)")
	if err!=nil{
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Username, user.Email, user.Password)
	return err
}

func (s *PostgresStore) GetUserByUsername(username string) (*User,error){
	var user User
	err := s.db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (s *PostgresStore) DeleteUser(id int) error{
	return nil
}
