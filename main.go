package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

func main() {
	const monzoAPI = "https://api.monzo.com"

	apiToken := os.Getenv("MONZO_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Get a Monzo API token from https://developers.monzo.com/api/playground and set it in your environment as `MONZO_API_TOKEN`.")
		return
	}

	checkMonzoTokenWorks(apiToken, monzoAPI)
	getUserDetails(apiToken, monzoAPI)
}

func checkMonzoTokenWorks(apiToken string, monzoAPI string) bool {
	type Ping struct {
		Authenticated bool   `json:"authenticated"`
		ClientID      string `json:"client_id"`
		UserID        string `json:"user_id"`
	}

	client := resty.New()
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

func getUserDetails(apiToken string, monzoAPI string) string {
	type Accounts struct {
		Accounts []struct {
			ID          string `json:"id"`
			Closed      bool   `json:"closed"`
			Created     string `json:"created"`
			Description string `json:"description"`
			Type        string `json:"type"`
			Currency    string `json:"currency"`
			CountryCode string `json:"country_code"`
			Owners      []struct {
				UserID             string `json:"user_id"`
				PreferredName      string `json:"preferred_name"`
				PreferredFirstName string `json:"preferred_first_name"`
			} `json:"owners"`
			AccountNumber int `json:"account_number"`
			SortCode      int `json:"sort_code"`
		} `json:"accounts"`
	}

	client := resty.New()
	resp, err := client.R().SetAuthToken(apiToken).Get(monzoAPI + "/accounts")

	if err != nil {
		fmt.Println("Something went wrong.")
	}

	fmt.Println(resp)

	parsedAccounts := Accounts{}
	json.Unmarshal(resp.Body(), &parsedAccounts)

	for i := 0; i < len(parsedAccounts.Accounts); i++ {
		if parsedAccounts.Accounts[i].Type == "uk_retail" {
			currentAccount := parsedAccounts.Accounts[i]
			fmt.Println("Found a current account belonging to " + currentAccount.Owners[0].PreferredName + ".")
		} else {
			continue
		}
	}

	return ""
}
