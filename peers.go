package simplecache

type Picker interface {
	Pick(key string) (Fetcher, bool)
}

type Fetcher interface {
	Fetch(group string, key string) ([]byte, error)
}
