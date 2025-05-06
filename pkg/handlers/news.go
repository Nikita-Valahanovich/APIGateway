package handlers

import (
	"APIGateway/pkg/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type API struct {
	r *mux.Router
}

func New() *API {
	a := API{r: mux.NewRouter()}
	a.endpoints()
	return &a
}

func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	api.r.Use(middleware.RequestIDMiddleware)
	api.r.Use(middleware.LoggingMiddleware)
	api.r.HandleFunc("/news", api.NewsList).Methods(http.MethodGet)
	api.r.HandleFunc("/news/{n}", api.posts).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comments", api.AddComment).Methods(http.MethodPost, http.MethodOptions)
}

// Получение списка новостей
func (api *API) NewsList(w http.ResponseWriter, r *http.Request) {
	// Чтение параметров запроса
	queryParams := r.URL.RawQuery

	// Чтение и проброс X-Request-ID
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = generateRequestID()
	}

	// Собираем URL с параметрами
	newsServiceURL := "http://localhost:80/news"
	if queryParams != "" {
		newsServiceURL += "?" + queryParams
	}

	req, err := http.NewRequest("GET", newsServiceURL, nil)
	if err != nil {
		http.Error(w, "Ошибка создания запроса", http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-Request-ID", requestID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Сервис новостей недоступен", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func generateRequestID() string {
	return uuid.New().String()
}

func (api *API) posts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n := vars["n"]

	// Прокси-запрос к сервису новостей
	resp, err := http.Get("http://localhost:80/news/" + n)
	if err != nil {
		http.Error(w, "Ошибка запроса к сервису новостей", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем заголовки и тело ответа
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
