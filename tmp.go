package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	name := "Oleg"
	request := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	response, err := http.Get(request)
	fmt.Println(response.StatusCode)
	fmt.Println(err)
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))

}
