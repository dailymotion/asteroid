package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	config "github.com/dailymotion/asteroid/pkg/config"
	_ "github.com/lib/pq"
)

// ConnectToDB to Postgres DB
func ConnectToDB(config config.ConfigDB) (*sql.DB, error) {
	dbInfo := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?connect_timeout=10&sslmode=disable",
		config.DbUsername, config.DbPassword,
		config.DbHost, config.DbPort, config.DbName)
	//sql.Open want a string for the DB connection, Sprintf gives it
	conn, err := sql.Open("postgres", dbInfo)
	//defer conn.Close()
	if err != nil {
		return nil, err
	}

	// We test the database's connection
	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

// InsertUserInDB will insert user info into the DB
func InsertUserInDB(db *sql.DB, key string, cidr string, name string) error {
	fmt.Println("Adding user in DB in Function")
	now := time.Now()

	sqlStatement := `
	INSERT INTO users (username, key, cidr, date)
	VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(sqlStatement, name, key, cidr, now)
	if err != nil {
		return err
	}
	fmt.Println("Adding user in DB DONE")
	return nil
}

// ReadKeyFromDB will retrieve user info from the DB
func ReadKeyFromDB(db *sql.DB) ([]config.User, error) {
	var user config.User
	var userList []config.User
	//layout := "2006-01-02"

	rows, err := db.Query("SELECT username, key, cidr, date FROM users;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&user.Username,
			&user.Key,
			&user.CIDR,
			&user.Date); err != nil {
			log.Fatal(err)
		}
		userList = append(userList, user)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println("DATE: ", user.Date.Format(layout))
	return userList, nil
}

// DeleteUserInDB will delete user from the DB
func DeleteUserInDB(db *sql.DB, key string) error {
	sqlStatement := `
	DELETE FROM users WHERE key = $1;`

	_, err := db.Exec(sqlStatement, key)
	if err != nil {
		return err
	}
	fmt.Println("Deleting user from DB DONE")
	return nil
}

