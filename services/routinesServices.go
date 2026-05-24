package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"project/utils"
)

func callAPI(url string, body []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	response, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		utils.AddLog(
			fmt.Sprintf(
				"Error calling %s : %v",
				url,
				err,
			),
		)
		return
	}

	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)

	if err != nil {
		utils.AddLog(
			fmt.Sprintf(
				"Error reading response from %s",
				url,
			),
		)
		return
	}

	utils.AddLog(
		fmt.Sprintf(
			"API: %s\nResponse: %s",
			url,
			string(respBody),
		),
	)
}

func CallMultipleApis() {

	var wg sync.WaitGroup

	apis := []struct {
		url  string
		body []byte
	}{
		{
			url: "http://localhost:8080/admin/create",
			body: []byte(`{
				"email":"raone123@gmail.com",
				"password":"12345"
			}`),
		},
		{
			url: "http://localhost:8080/admin/login",
			body: []byte(`{
				"email":"akshit1@gmail.com",
				"password":"aksh1234"
			}`),
		},
	}

	for _, api := range apis {

		wg.Add(1)

		go callAPI(
			api.url,
			api.body,
			&wg,
		)
	}

	wg.Wait()
}