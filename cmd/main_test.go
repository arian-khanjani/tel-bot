package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestConn(t *testing.T) {
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot6277771629:AAGlZHcCRs80t3uvC6xXUfllP0LjdLxyJlY/getMe", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatal("status was not 200")
	}

	var out map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("-- OUT --")
	for i, v := range out {
		fmt.Println(i, v)
	}
}
