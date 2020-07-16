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

// Balance implements the data structure for the `/balance` Monzo API response
// for current accounts. Docs at https://docs.monzo.com/#balance.
type Balance struct {
	Balance                         float64    `json:"balance"`
	TotalBalance                    float64    `json:"total_balance"`
	BalanceIncludingFlexibleSavings float64    `json:"balance_including_flexible_savings"`
	Currency                        string     `json:"currency"`
	SpendToday                      int        `json:"spend_today"`
	LocalCurrency                   string     `json:"local_currency"`
	LocalExchangeRate               int        `json:"local_exchange_rate"`
	LocalSpend                      []struct{} `json:"local_spend"`
}

// Ping implements the data structure for the `/ping/whoami` Monzo API response.
// Docs at https://docs.monzo.com/#authenticating-requests.
type Ping struct {
	Authenticated bool   `json:"authenticated"`
	ClientID      string `json:"client_id"`
	UserID        string `json:"user_id"`
}

// Pots implements the data structure for the `/pots` Monzo API response.
// Docs at https://docs.monzo.com/#pots.
type Pots struct {
	Pots []struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		Style     string  `json:"style"`
		Balance   float64 `json:"balance"`
		Currency  string  `json:"currency"`
		CreatedAt string  `json:"created"`
		UpdatedAt string  `json:"updated"`
		Deleted   bool    `json:"deleted"`
	} `json:"pots"`
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
	getPotsBalance(accountID, apiToken, monzoAPI)
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

	parsedBalance := Balance{}
	json.Unmarshal(resp.Body(), &parsedBalance)

	fmt.Println(fmt.Sprintf("Current account balance: %v %s.", parsedBalance.Balance/100, parsedBalance.Currency))
	fmt.Println(fmt.Sprintf("Total balance (including savings): %v %s.", parsedBalance.TotalBalance/100, parsedBalance.Currency))

	return ""
}

func getPotsBalance(accountID string, apiToken string, monzoAPI string) string {
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"current_account_id": accountID,
		}).
		SetAuthToken(apiToken).Get(monzoAPI + "/pots")

	if err != nil {
		fmt.Println("Something went wrong.")
	}

	parsedPots := Pots{}
	json.Unmarshal(resp.Body(), &parsedPots)

	for i := 0; i < len(parsedPots.Pots); i++ {

		if parsedPots.Pots[i].Deleted == true {
			continue
		}

		fmt.Println(fmt.Sprintf("The %s pot contains %v %s.", parsedPots.Pots[i].Name, parsedPots.Pots[i].Balance/100, parsedPots.Pots[i].Currency))
	}

	return ""

}
