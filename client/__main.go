package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"
)

func sendData() {
	// Define the data to be sent
	data := []byte(`{"foo":"bar"}`)

	// Create a new request with the data
	req, err := http.NewRequest("POST", "https://example.com/api/data", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers if necessary
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	fmt.Println("Data sent successfully!")
}

func main() {
	// Create a new cron job that runs sendData() every minute
	c := cron.New()
	c.AddFunc("*/1 * * * *", sendData)
	c.Start()

	// Wait forever
	select {}
}
