package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

func main() {

	// CLI input
	url := flag.String("url", "", "URL to shorten")
	flag.Parse()

	if *url == "" {
		fmt.Println("Please provide a URL using -url")
		return
	}

	// create JSON body
	data := map[string]string{
		"url": *url,
	}

	jsonData, _ := json.Marshal(data)

	// POST request
	resp, err := http.Post(
		"http://localhost:3000/api/v1",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println("Response:")
	fmt.Println(string(body))
}
