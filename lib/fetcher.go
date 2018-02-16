package lib

type FetchResponse struct {
	Price Price
	Err   error
}

type FetchWithHistoryResponse struct {
	Price    Price
	HistRate float64
	Err      error
}

type Fetcher interface {
	Fetch(chan FetchResponse)
	FetchWithHistory(chan FetchWithHistoryResponse)
}
