package main

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
	client := resty.New()
	resp, err := client.R().SetAuthToken(apiToken).Get(monzoAPI + "/ping/whoami")

	if err != nil {
		fmt.Println("Something went wrong.")
	}

	fmt.Println(resp)
	return true
}
