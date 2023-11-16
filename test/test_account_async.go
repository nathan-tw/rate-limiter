package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func sendRequest(url, account string, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("account", account)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response for account %s: %s\n", account, string(body))
}

func main() {
	baseURL := "http://127.0.0.1:8080/api/" +
		"path1"
	account := "Account1"

	var wg sync.WaitGroup

	for i := 0; i < 11; i++ {
		wg.Add(1)
		go sendRequest(baseURL, account, &wg)
	}

	wg.Wait()
}
