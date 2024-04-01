package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Email struct {
	To      []Recipient `json:"to"`
	From    Recipient   `json:"from"`
	Subject string      `json:"subject"`
	HTML    string      `json:"html"`
}

type Recipient struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func SendEmail(email, subject, body string) {
	emailData := Email{
		To: []Recipient{
			{
				Name:  email,
				Email: email,
			},
		},
		From: Recipient{
			Name:  "From User Name",
			Email: "from@chrisbrocklesby.com",
		},
		Subject: subject,
		HTML:    body,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", "https://email.chrisbrocklesby.workers.dev", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("EMAIL_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println("Email sent to:", email)
}
