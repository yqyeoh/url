package app

type ShortenURLReq struct {
	URL string `json:"url"`
}
type GetOriginalURLReq struct {
	URL string `form:"url"`
}

type URLInfo struct {
	URL  string `json:"url" db:"url"`
	Code string `json:"code" db:"code"`
}

type ShortenURLResponse struct {
	Original  string `json:"original"`
	Shortened string `json:"shortened"`
}

type URLRes struct {
	URL string `json:"url"`
}
