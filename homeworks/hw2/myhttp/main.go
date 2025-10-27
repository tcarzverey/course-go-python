package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/tcarzverey/course-go-python/homeworks/hw2/myhttp/client"
	"github.com/tcarzverey/course-go-python/homeworks/hw2/myhttp/server"
)

// main Пример использования связки нашего сервера-клиента-обработчика
func main() {
	myServer := server.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		fmt.Println("Server started at port", port)
		http.HandleFunc("/test", MyHandler)
		err := myServer.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	myClient := client.New()

	// Создаем запрос
	reqURL, _ := url.Parse(fmt.Sprintf("http://localhost:%s/test?name=Test", port))
	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {"abc"},
		},
	}

	resp, err := myClient.Do(req)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Headers: %v\n", resp.Header)
	fmt.Printf("Body: %v\n", string(body))
}
