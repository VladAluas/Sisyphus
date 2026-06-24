// Package db manages the connection with the metadata database
package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDatabase(
	host string,
	port string,
	name string,
	user string,
	pass string,
) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host,
		port,
    name,
    user,
    pass,
 )

	return sql.Open("pgx", dsn)
}
