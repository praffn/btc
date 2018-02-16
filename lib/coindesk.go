package lib

import (
	"encoding/json"
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
	// 2013-09-18T17:27:00+00:00
	updated, err := time.Parse("2006-01-02T15:04:05+00:00", resp.Time["updatedISO"])
	if err != nil {
		return Price{}, err
	}
	// return price
	return Price{
		Updated:  updated,
		Rate:     bpi.Rate,
		Currency: c.currency,
	}, nil
}

func (c *coindesk) FetchWithHistory() (Price, float64, error) {
	price, err := c.Fetch()
	if err != nil {
		return Price{}, 0.0, err
	}
	r, err := c.client.Get(makeHistoricalCoindeskURL(c.currency))
	if err != nil {
		return Price{}, 0.0, err
	}

	var histResp coindeskHistoricalResponse
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&histResp)
	if err != nil {
		return Price{}, 0.0, err
	}
	var histRate float64
	for _, rate := range histResp.BPI {
		histRate = rate
		break
	}

	return price, histRate, nil
}

func NewCoindesk(currency string) Fetcher {
	url := makeCoindeskURL(currency)
	return &coindesk{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url:      url,
		currency: currency,
	}
}
