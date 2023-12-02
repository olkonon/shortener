package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/olkonon/shortener/internal/app/api"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const (
	ContentEncodingHeader      = "Content-Encoding"
	ContentTypeHeader          = "Content-Type"
	ContentTypeApplicationJSON = "application/json"
)

func New(config Config) *Handler {
	buf := make([]byte, 16)
	_, err := rand.Read(buf) // записываем байты в массив b
	if err != nil {
		log.Fatal(err)
	}
	return &Handler{
		store:     config.Store,
		baseURL:   config.BaseURL,
		dsn:       config.DSN,
		secretKey: buf,
	}
}

type Config struct {
	BaseURL string
	DSN     string
	Store   storage.Storage
}

type Handler struct {
	store     storage.Storage
	dsn       string
	secretKey []byte
	baseURL   string
}

func (h *Handler) GET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	longURL, err := h.store.GetURLByID(r.Context(), vars["id"])
	if err != nil {
		if errors.Is(err, storage.ErrDeletedURL) {
			w.WriteHeader(http.StatusGone)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func (h *Handler) UserGET(w http.ResponseWriter, r *http.Request) {
	urlList, err := h.store.GetByUser(r.Context(), mux.Vars(r)[common.MuxUserVarName])
	if errors.Is(err, storage.ErrUserURLListEmpty) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make([]api.UserGetResponse, len(urlList))
	for i, val := range urlList {
		response[i].OriginalURL = val.OriginalURL
		response[i].ShortURL = fmt.Sprintf("%s/%s", h.baseURL, val.ShortID)
	}

	buf, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("JSON serialization error:", err)
		return
	}

	w.Header().Set(ContentTypeHeader, ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)
	if _, tmpErr := w.Write(buf); tmpErr != nil {
		log.Error(tmpErr)
	}
}

func (h *Handler) POST(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	longURL := string(b)

	//Проверка, что переданный URl корректный
	if !common.IsValidURL(longURL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.store.GenIDByURL(r.Context(), longURL, mux.Vars(r)[common.MuxUserVarName])
	if errors.Is(err, storage.ErrDuplicateURL) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf("%s/%s", h.baseURL, id)))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, tmpErr := w.Write([]byte(fmt.Sprintf("%s/%s", h.baseURL, id))); tmpErr != nil {
		log.Error(tmpErr)
	}
}

func (h *Handler) PostJSON(w http.ResponseWriter, r *http.Request) {
	successStatusCode := http.StatusCreated
	if r.Header.Get(ContentTypeHeader) != ContentTypeApplicationJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := api.AddURLRequest{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("JSON deserialization error:", err)
		return
	}

	//Проверка, что переданный URl корректный
	if !common.IsValidURL(data.URL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.store.GenIDByURL(r.Context(), data.URL, mux.Vars(r)[common.MuxUserVarName])
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateURL) {
			successStatusCode = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	response := api.AddURLResponse{
		Result: fmt.Sprintf("%s/%s", h.baseURL, id),
	}

	buf, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("JSON serialization error:", err)
		return
	}

	w.Header().Set(ContentTypeHeader, ContentTypeApplicationJSON)
	w.WriteHeader(successStatusCode)
	if _, tmpErr := w.Write(buf); tmpErr != nil {
		log.Error(tmpErr)
	}
}

func (h *Handler) BatchPostJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(ContentTypeHeader) != ContentTypeApplicationJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make([]api.BatchAddURLRequest, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("JSON deserialization error:", err)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Empty request")
		return
	}

	batchUpdate := make([]storage.BatchSaveRequest, len(data))
	for i, val := range data {
		//Проверка, что переданный URl корректный
		if !val.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		batchUpdate[i].OriginalURL = val.OriginalURL
		batchUpdate[i].CorrelationID = val.CorrelationID
	}

	batchResponse, err := h.store.BatchSave(r.Context(), batchUpdate, mux.Vars(r)[common.MuxUserVarName])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make([]api.BatchAddURLResponse, len(batchResponse))
	for i, val := range batchResponse {
		response[i].CorrelationID = val.CorrelationID
		response[i].ShortURL = fmt.Sprintf("%s/%s", h.baseURL, val.ShortID)
	}

	buf, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("JSON serialization error:", err)
		return
	}

	w.Header().Set(ContentTypeHeader, ContentTypeApplicationJSON)
	w.WriteHeader(http.StatusCreated)
	if _, tmpErr := w.Write(buf); tmpErr != nil {
		log.Error(tmpErr)
	}
}

func (h *Handler) BatchDeleteJSON(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(ContentTypeHeader) != ContentTypeApplicationJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make([]string, 0)
	err = json.Unmarshal(b, &data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Error("JSON deserialization error:", err)
		return
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		log.Error("Empty request")
		return
	}

	h.store.BatchDelete(r.Context(), data, mux.Vars(r)[common.MuxUserVarName])
	w.WriteHeader(http.StatusAccepted)
}

func (h Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	if h.dsn == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("postgres", h.dsn)
	if err != nil {
		log.Error("DB connect error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		log.Error("DB Ping error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
