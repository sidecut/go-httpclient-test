package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	client := new(http.Client)
	response, err := client.Get("https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	if response.StatusCode != 200 {
		panic(fmt.Sprintf("response.Status = %v", response.Status))
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s", body)
}
