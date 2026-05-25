package services

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"

	"project/utils"
)

func callAPI(url string, body []byte, token string, wg *sync.WaitGroup) {

	defer wg.Done()

	// Create request
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(body),
	)

	if err != nil {

		utils.AddLog(
			fmt.Sprintf(
				"Error creating request for %s : %v",
				url,
				err,
			),
		)

		return
	}

	// Add headers
	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	// Add JWT token
	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	// HTTP client
	client := &http.Client{}

	// Execute request
	response, err := client.Do(req)

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
		token string
	}{
		{
			url: "http://localhost:8080/api/student/create",
			body: []byte(`{
				"name":"Raone",
				"age":25,
				"email":"raone123@gmail.com"
			}`),
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imxtbm9AZ21haWwuY29tIiwiZXhwIjoxNzc5NzY4NDMyLCJyb2xlIjoiYWRtaW4ifQ.qWBO_5RHCipgiDAyiGP74MkEQnZc-1CT49mzLaybUMk",
		},
		{
			url: "http://localhost:8080/admin/login",
			body: []byte(`{
				"email":"admin@gmail.com",
				"password":"admin123"
			}`),
		},
	}

	for _, api := range apis {

		wg.Add(1)

		go callAPI(
			api.url,
			api.body,
			api.token,
			&wg,
		)
	}

	wg.Wait()
}