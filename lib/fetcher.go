package lib

type Fetcher interface {
	Fetch() (Price, error)
	FetchWithHistory() (Price, float64, error)
}
