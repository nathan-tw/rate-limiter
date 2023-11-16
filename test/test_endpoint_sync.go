package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}

	url := "http://127.0.0.1:8080/api/path1"
	headers1 := map[string]string{"account": "Account1"}
	headers2 := map[string]string{"account": "Account2"}

	for i := 1; i <= 5; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		for key, value := range headers1 {
			req.Header.Add(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response for account %s: %s\n", req.Header.Get("account"), string(body))
		time.Sleep(100 * time.Millisecond) // Add a delay between requests
	}

	for i := 6; i <= 11; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		for key, value := range headers2 {
			req.Header.Add(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response for account %s: %s\n", req.Header.Get("account"), string(body))
		time.Sleep(100 * time.Millisecond)
	}
}
