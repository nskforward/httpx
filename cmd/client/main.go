package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	req, err := http.NewRequest("GET", "http://", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}
