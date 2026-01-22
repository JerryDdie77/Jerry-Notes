package main

import (
	"jerry-notes/internal/api"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Создаём роутер и регистрируем все хендлеры в нём
	r := mux.NewRouter()
	api.RegisterRoutes(r)

	// Запускаем сервер на порту :8080
	http.ListenAndServe(":8080", r)
}
