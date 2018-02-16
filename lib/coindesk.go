package lib

import (
	"encoding/json"
	"net/http"
	"time"
)

func makeCoindeskURL(currency string) string {
	return "https://api.coindesk.com/v1/bpi/currentprice/" + currency + ".json"
}

type coindesk struct {
	client   *http.Client
	currency string
}

type coindeskBpi struct {
	Rate float64 `json:"rate_float"`
}

type coindeskResponse struct {
	Time map[string]string      `json:"time"`
	BPI  map[string]coindeskBpi `json:"bpi"`
}

func (c *coindesk) Fetch() (Price, error) {
	// fetch json
	url := makeCoindeskURL(c.currency)
	r, err := c.client.Get(url)
	if err != nil {
		return Price{}, err
	}
	// marshal body into coindeskResponse
	var resp coindeskResponse
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&resp)
	if err != nil {
		return Price{}, err
	}
	// get bpi for currency
	// (we are no being defensive, since bpi for currency
	// should exist for the resource)
	bpi := resp.BPI[c.currency]
	// return price
	return Price{
		Updated:  time.Now(),
		Rate:     bpi.Rate,
		Currency: c.currency,
	}, nil
}

func NewCoindesk(currency string) Fetcher {
	return &coindesk{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		currency: currency,
	}
}
