package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {

	url := "http://blogbid.000webhostapp.com/api/categories/read.php"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	xyz := string(body)
	fmt.Println(json.Marshal(xyz))
}
