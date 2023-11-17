package api

import "github.com/olkonon/shortener/internal/app/common"

type AddURLRequest struct {
	URL string `json:"url"`
}

type AddURLResponse struct {
	Result string `json:"result"`
}

type BatchAddURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

func (br *BatchAddURLRequest) IsValid() bool {
	return common.IsValidURL(br.OriginalURL)
}

type BatchAddURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
