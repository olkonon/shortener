package api

type AddURLRequest struct {
	URL string `json:"url"`
}

type AddURLResponse struct {
	Result string `json:"result"`
}
