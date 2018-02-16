package lib

type Fetcher interface {
	Fetch() (Price, error)
}
