package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type CommentRequest struct {
	NewsID   int    `json:"news_id"`
	ParentID *int   `json:"parent_id"`
	Content  string `json:"content"`
	Author   string `json:"author"`
}

// Добавление комментария
func (api *API) AddComment(w http.ResponseWriter, r *http.Request) {
	var req CommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Отправка запроса на CensorshipService
	censorReq := map[string]string{"content": req.Content}
	body, _ := json.Marshal(censorReq)

	resp, err := http.Post("http://localhost:8082/censor", "application/json", bytes.NewReader(body))
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Комментарий отклонен модератором", http.StatusBadRequest)
		return
	}

	// Теперь проксируем оригинальный запрос в CommentService
	commentBody, _ := json.Marshal(req)
	commentResp, err := http.Post("http://localhost:8081/comments", "application/json", bytes.NewReader(commentBody))
	if err != nil || commentResp.StatusCode != http.StatusCreated {
		http.Error(w, "Ошибка при сохранении комментария", http.StatusInternalServerError)
		return
	}

	// Возврат клиенту успешного ответа от CommentService
	w.WriteHeader(http.StatusCreated)
	io.Copy(w, commentResp.Body)
	commentResp.Body.Close()

}
