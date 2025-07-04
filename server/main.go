package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func sendReadFail() {

}

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
	response_json, err := json.Marshal(response)
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
  user     = "postgres"
  password = "your-password"
  dbname   = "calhounio_demo"
)

var pool *sql.DB

func connectToDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
		return err
	}
	
}

func connectToDBOld() {
	id := flag.Int64()
	dsn := flag.String("dsn", )
	pool, err := sql.Open("pq", *dsn)
	appSignal := make(chan os.Signal, 3)
	ctx, stop := context.WithCancel(context.Background())

	go func() {
		<-appSignal
		stop()
	}

	Ping(ctx)
	Query(ctx, *id)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", serveRoot)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Could not set up server on port.")
		return
	}
	connectToDB()
}
