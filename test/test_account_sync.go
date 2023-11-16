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
	headers := map[string]string{"account": "Account1"}

	for i := 1; i <= 11; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		for key, value := range headers {
			req.Header.Add(key, value)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}

		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response for account Account1: %s\n", string(body))
		time.Sleep(100 * time.Millisecond)
	}
}
