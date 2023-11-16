package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func sendRequest2(url, account string, wg *sync.WaitGroup) {
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
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response for account %s: %s\n", account, string(body))
}

func main() {
	baseURL := "http://127.0.0.1:8080/api/path1"
	account1 := "Account1"
	account2 := "Account2"

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go sendRequest2(baseURL, account1, &wg)
	}

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go sendRequest2(baseURL, account2, &wg)
	}

	wg.Wait()
}
