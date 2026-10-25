package main

import (
	"fmt"
	"net/http"

	"github.com/gaudiot/sezzle-challenge/httpapi"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	calculatorHandler := httpapi.NewCalculatorHandler()
	router.Route("/math", calculatorHandler.Router)

	fmt.Println("Server is running on port :8000")
	http.ListenAndServe(":8000", router)
}
