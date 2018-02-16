package lib

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func makeCoindeskURL(currency string) string {
	return "https://api.coindesk.com/v1/bpi/currentprice/" + currency + ".json"
}

func makeHistoricalCoindeskURL(currency string) string {
	return "https://api.coindesk.com/v1/bpi/historical/close.json?currency=" + currency + "&for=yesterday"
}

type coindesk struct {
	client   *http.Client
	currency string
	url      string
}

type coindeskBpi struct {
	Rate float64 `json:"rate_float"`
}

type coindeskResponse struct {
	Time map[string]string      `json:"time"`
	BPI  map[string]coindeskBpi `json:"bpi"`
}

type coindeskHistoricalResponse struct {
	BPI  map[string]float64 `json:"bpi"`
	Time map[string]string  `json:"time"`
}

func (c *coindesk) Fetch(ch chan FetchResponse) {
	// fetch json
	url := makeCoindeskURL(c.currency)
	log.Printf("Fetching JSON from CoinDesk: %s\n", url)
	r, err := c.client.Get(url)
	if err != nil {
		log.Printf("Error fetching: %s\n", err.Error())
		ch <- FetchResponse{Price{}, err}
		return
	}
	log.Println("Successfully fetched JSON")
	log.Println("Marshalling JSON")
	// marshal body into coindeskResponse
	var resp coindeskResponse
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s\n", err.Error())
		ch <- FetchResponse{Price{}, err}
		return
	}
	// get bpi for currency
	// (we are no being defensive, since bpi for currency
	// should exist for the resource)
	bpi := resp.BPI[c.currency]

	log.Println("Parsing date")
	updated, err := time.Parse("2006-01-02T15:04:05+00:00", resp.Time["updatedISO"])
	if err != nil {
		log.Printf("Error parsing data: %s\n", err.Error())
		ch <- FetchResponse{Price{}, err}
		return
	}
	log.Println("Successfully parsed date")
	log.Println("Successfully fetched price")
	// return price
	ch <- FetchResponse{
		Price{
			Updated:  updated,
			Rate:     bpi.Rate,
			Currency: c.currency,
		},
		nil,
	}
}

func (c *coindesk) FetchWithHistory(ch chan FetchWithHistoryResponse) {
	pch := make(chan FetchResponse, 1)
	log.Println("Getting price")
	go c.Fetch(pch)
	resp := <-pch
	close(pch)
	if resp.Err != nil {
		log.Printf("Error getting price: %s\n", resp.Err.Error())
		ch <- FetchWithHistoryResponse{Price{}, 0.0, resp.Err}
		return
	}
	log.Println("Fetching historical data")
	r, err := c.client.Get(makeHistoricalCoindeskURL(c.currency))
	if err != nil {
		log.Printf("Error fetching historical data: %s\n", err.Error())
		ch <- FetchWithHistoryResponse{Price{}, 0.0, err}
		return
	}

	log.Println("Successfully fetched historical data")
	log.Println("Marshalling historical data")
	var histResp coindeskHistoricalResponse
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&histResp)
	if err != nil {
		log.Printf("Error marshalling historical data: %s\n", err.Error())
		ch <- FetchWithHistoryResponse{Price{}, 0.0, err}
		return
	}
	log.Println("Successfully marshalled historical data")
	var histRate float64
	for _, rate := range histResp.BPI {
		histRate = rate
		break
	}
	log.Println("Succesfully fetched historical data, returning")
	ch <- FetchWithHistoryResponse{resp.Price, histRate, err}
}

func NewCoindesk(currency string) Fetcher {
	log.Printf("Creating new CoinDesk fetcher for currency %s\n", currency)
	url := makeCoindeskURL(currency)
	return &coindesk{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url:      url,
		currency: currency,
	}
}
