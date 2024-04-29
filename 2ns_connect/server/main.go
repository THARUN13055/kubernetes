package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func gethello(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://pkg.go.dev/net/http@go1.21.6")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	Body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(Body))
	writefile(Body)
	w.Write(Body)
}

func writefile(content []byte) {
	file, err := os.Create("logs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString(string(content))
}

func main() {

	http.HandleFunc("/", gethello)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
