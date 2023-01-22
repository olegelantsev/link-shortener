package shortener

type CreateUrlRequest struct {
	Url   string `json:"URL"`
	Title string `json:"Title"`
}

type CreateUrlResponse struct {
	Url string `json:"URL"`
}
