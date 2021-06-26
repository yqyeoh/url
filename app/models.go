package app

type ShortenURLReq struct {
	URL string `json:"url"`
}

type URLInfo struct {
	URL  string `json:"url" db:"url"`
	Code string `json:"code" db:"code"`
}

type ShortenURLResponse struct {
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
}
