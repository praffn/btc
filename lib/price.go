package lib

import "time"

// Price represents the price for a single BTC
// in a given currency at a specific point in time
type Price struct {
	Currency string
	Rate     float64
	Updated  time.Time
}
