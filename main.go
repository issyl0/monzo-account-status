package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
)

// Accounts implements the data structure of the `/accounts` Monzo API response.
// Docs at https://docs.monzo.com/#accounts.
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
		AccountNumber string `json:"account_number"`
		SortCode      string `json:"sort_code"`
	} `json:"accounts"`
}

// Ping implements the data structure for the `/ping/whoami` Monzo API response.
// Docs at https://docs.monzo.com/#authenticating-requests.
type Ping struct {
	Authenticated bool   `json:"authenticated"`
	ClientID      string `json:"client_id"`
	UserID        string `json:"user_id"`
}

func main() {
	const monzoAPI = "https://api.monzo.com"

	apiToken := os.Getenv("MONZO_API_TOKEN")
	if apiToken == "" {
		fmt.Println("Get a Monzo API token from https://developers.monzo.com/api/playground and set it in your environment as `MONZO_API_TOKEN`.")
		return
	}

	checkMonzoTokenWorks(apiToken, monzoAPI)
	accountID := getUserDetails(apiToken, monzoAPI)
	getCurrentAccountBalance(accountID, apiToken, monzoAPI)
}

func checkMonzoTokenWorks(apiToken string, monzoAPI string) bool {
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
	client := resty.New()
	resp, err := client.R().SetAuthToken(apiToken).Get(monzoAPI + "/accounts")

	if err != nil {
		fmt.Println("Something went wrong.")
	}

	parsedAccounts := Accounts{}
	json.Unmarshal(resp.Body(), &parsedAccounts)

	accountID := ""
	for i := 0; i < len(parsedAccounts.Accounts); i++ {
		if parsedAccounts.Accounts[i].Type == "uk_retail" {
			currentAccount := parsedAccounts.Accounts[i]
			accountID = currentAccount.ID
			fmt.Println("Found a current account (" + currentAccount.AccountNumber + ") belonging to " + currentAccount.Owners[0].PreferredName + ".")
		}
	}
	return accountID
}

func getCurrentAccountBalance(accountID string, apiToken string, monzoAPI string) string {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"account_id": accountID,
		}).
		SetAuthToken(apiToken).Get(monzoAPI + "/balance")

	if err != nil {
		fmt.Println("Something went wrong.")
	}

	fmt.Println(resp)

	return ""
}
