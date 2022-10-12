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
	book_update := []byte(`{"title": "Animal farm", "author": "George Orwell"}`)

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

	if TestPutRequest(baseURL+"/books/1", book_update) == http.StatusOK {
		TestOK()
	} else {
		TestFail()
	}

	if TestDeleteRequest(baseURL+"/books/1") == http.StatusNoContent {
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

	r, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))

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
		fmt.Printf(" Error: status code %d", r.StatusCode)
	}

	return r.StatusCode
}

func TestPutRequest(url string, jsonBody []byte) int {
	fmt.Printf("PUT " + url)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.ContentLength = int64(bytes.NewBuffer(jsonBody).Len())

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf(" Response: " + bodyString)
	} else {
		fmt.Printf(" Error: status code %d", res.StatusCode)
	}

	return res.StatusCode
}

// func TestPutRequest(url string, jsonBody []byte) int {
// 	fmt.Printf("PUT " + url)

// 	r, err := http.Put(url, "application/json", bytes.NewReader(jsonBody))

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if r.StatusCode == http.StatusOK {
// 		bodyBytes, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		bodyString := string(bodyBytes)
// 		fmt.Printf(" Response: " + bodyString)
// 	} else {
// 		fmt.Printf(" Error: status code %d", r.StatusCode)
// 	}

// 	return r.StatusCode
// }

func TestDeleteRequest(url string) int {
	fmt.Printf("DELETE " + url)
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusNoContent {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Printf(" Response: " + bodyString)
	} else {
		fmt.Printf(" Error: status code %d", res.StatusCode)
	}

	return res.StatusCode
}
