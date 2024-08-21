package httphandler

import "net/http"

// Ping
// @Summary Пинг-сервис
// @Description Простая проверка доступности сервиса. Возвращает строку "pong".
// @Tags health
// @Produce plain
// @Success 200 {string} string "pong"
// @Router /ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
