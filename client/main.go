package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RootResponse struct {
	Success bool
	Err     error
}

func main() {
	res, err := http.Get("http://localhost:8080")
	if err != nil {
		fmt.Println("error", err)
		return
	}
	res_bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	var rootRes RootResponse
	err = json.Unmarshal(res_bytes, &rootRes)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	print(rootRes.Success)
}
