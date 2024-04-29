package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	var serveradd string

	serveradd = "http://localhost:8080/"
	resp, err := http.Get(serveradd)
	if err != nil {
		log.Fatal(err)
	}
	read, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(read))
}
