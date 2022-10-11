package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	baseURL := "http://localhost:8001/api/v1"
	book := []byte(`{"title": "1984", "author": "George Orwell"}`)

	if TestGetRequest(baseURL+"/books") == http.StatusOK {
		TestOK()
	} else {
		TestFail()
	}

	if TestPostRequest(baseURL+"/books", book) == http.StatusCreated {
		TestOK()
	} else {
		TestFail()
	}

	if TestGetRequest(baseURL+"/books/1") == http.StatusOK {
		TestOK()
	} else {
		TestFail()
	}
}

func TestFail() {
	fmt.Printf(" -\033[31m FAIL \033[37m\n")
}

func TestOK() {
	fmt.Printf(" -\033[32m OK \033[37m\n")
}

func TestGetRequest(url string) int {
	fmt.Printf("GET " + url)

	r, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf(" Response: " + bodyString)
	} else {
		fmt.Printf(" Error: status code %d", r.StatusCode)
	}

	return r.StatusCode
}

func TestPostRequest(url string, jsonBody []byte) int {
	fmt.Printf("POST " + url)

	bodyReader := bytes.NewReader(jsonBody)
	r, err := http.Post(url, "application/json", bodyReader)

	if err != nil {
		log.Fatal(err)
	}

	if r.StatusCode == http.StatusCreated {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf(" Response: " + bodyString)
	} else {
		fmt.Printf("Error: status code %d", r.StatusCode)
	}

	return r.StatusCode
}
