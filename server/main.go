package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq"
)

type RootResponse struct {
	Success bool
	Err     error
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	_, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	response := RootResponse{
		Success: true,
		Err:     nil,
	}
	response_json, err := json.Marsha(response)
	w.Write(response_json)
	fmt.Println("wrote response:", response)
	var response2 RootResponse
	fmt.Println(string(response_json))
	json.Unmarshal(response_json, &response2)
	fmt.Println(response2.Success)
	if err != nil {
		fmt.Println("error", err)
		return
	}
}

const (
	host     = "localhost"
	port     = 5432
	user     = "jublizoo"
	password = "origami2003"
	dbname   = "mydb"
)

var pool *sql.DB

// func getAllUsers(db *sql.DB) error {
// 	db.Exec()
// }

func addUser(db *sql.DB, username string, passwd string) error {
	password_hash := passwd

	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
	`
	result, err := db.Exec(query, username, password_hash)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		return err
	}
	fmt.Println("Add user result:", result)

	return nil
}

func getUsers(db *sql.DB) error {
	query := `
		SELECT username, password_hash FROM users;
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to get users:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var username string
		var pass string
		err := rows.Scan(&username, &pass)
		if err != nil {
			fmt.Println("Error parsing query result:", err)
			return err
		}
		fmt.Println("username:", username, "password:", pass)
	}

	return nil
}

func connectToDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return nil, err
	}
	fmt.Println("Connected to DB")

	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging:", err)
		return nil, err
	}

	fmt.Println("Successfully connected!")
	return db, nil
}

func main() {
	db, err := connectToDB()
	if err != nil {
		return
	}
	defer db.Close()

	addUser(db, "jamjam", "ploopy")
	addUser(db, "jublizoo", "ploopy")
	getUsers(db)

	router := http.NewServeMux()
	router.HandleFunc("GET /", serveRoot)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Could not set up server on port.")
		return
	}
}
