package main

import "encoding/json"
import "fmt"
import "github.com/go-resty/resty/v2"
import "os"

func main() {
	const monzoAPI = "https://api.monzo.com"

	apiToken := os.Getenv("MONZO_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Get a Monzo API token from https://developers.monzo.com/api/playground and set it in your environment as `MONZO_API_TOKEN`.")
		return
	}

	checkMonzoTokenWorks(apiToken, monzoAPI)
}

func checkMonzoTokenWorks(apiToken string, monzoAPI string) bool {
	type Ping struct {
		Authenticated bool   `json:"authenticated"`
		ClientID      string `json:"client_id"`
		UserID        string `json:"user_id"`
	}

	client := resty.New()
	var resp *resty.Response

	resp, err := client.R().SetResult(&Ping{}).SetAuthToken(apiToken).Get(monzoAPI + "/ping/whoami")

	if err != nil {
		fmt.Println("Something went wrong.")
	}

	parsedPing := Ping{}
	json.Unmarshal(resp.Body(), &parsedPing)

	if parsedPing.Authenticated == true {
		fmt.Println("You are authenticated with the Monzo API.")
	} else {
		fmt.Println("You are not authenticated with the Monzo API. Try again.")
	}

	return parsedPing.Authenticated
}
