package copper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Person ... Generic new person struct
type Person struct {
	Name   string        `json:"name"`
	Emails []interface{} `json:"emails"`
}

// SetEmail ... sets email for newly created person
func (person *Person) SetEmail(email string) {
	person.Emails = []interface{}{
		map[string]string{
			"email": email,
		},
	}
}

// Bufferize ... create byte array and subsequent buffer read stream from person struct
func (person *Person) Bufferize() *bytes.Buffer {
	p, _ := json.Marshal(person)
	return bytes.NewBuffer(p)
}

// CreateCopperPerson ... Create person record in Copper CRM
func CreateCopperPerson(user map[string]interface{}, apiKey string, apiEmail string) int {
	// Initialize a Person
	person := new(Person)
	// Set person's email no matter what
	person.SetEmail(user["email"].(string))

	// Check user role for proper name assignment
	if user["role"] == "guest" {
		person.Name = "CASHEWSCHASHEWS"
	} else {
		fullName := user["first"].(string) + user["last"].(string)
		person.Name = fullName
	}
	fmt.Println(person)
	// Initialize http client
	httpClient := &http.Client{}
	// Make a Post request to Copper
	request, err := http.NewRequest("POST", "https://api.prosperworks.com/developer_api/v1/people", person.Bufferize())
	request.Header.Add("X-PW-AccessToken", apiKey)
	request.Header.Add("X-PW-Application", "developer_api")
	request.Header.Add("X-PW-UserEmail", apiEmail)
	request.Header.Add("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error building Copper request: %v", err)
	}

	// Send the request to Copper
	resp, err := httpClient.Do(request)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("Error sending request to Copper: %v", err)
	}
	// defer resp.Body.Close()

	// Parse & validate the Copper response, get the copper ID
	decoder := json.NewDecoder(resp.Body)
	var copperResponse map[string]interface{}
	decoder.Decode(&copperResponse)

	// Copper ID response is float 64
	c := copperResponse["id"].(float64)

	// Convert float 64 to int
	var ID int = int(c)

	// Return copper id
	return ID
}
