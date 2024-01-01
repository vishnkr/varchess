package db

import (
	"context"
	"fmt"
	"varchess/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Database struct{
	*pgxpool.Pool
}

func Connect(cfg config.DBConfig) (*Database, error){
	dbUrl := getConnURL(cfg)
	db, err:= pgxpool.New(context.Background(),dbUrl)
	if err!=nil{
		return nil,err
	}
	return &Database{db},nil
}

func getConnURL(cfg config.DBConfig) string{
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name)
}