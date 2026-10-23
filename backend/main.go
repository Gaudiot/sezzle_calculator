package main

import (
	"fmt"
	"net/http"

	math "github.com/gaudiot/sezzle-challenge/src"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	mathRouter := math.NewMathRouter()
	router.Route("/math", mathRouter.Router)

	fmt.Println("Server is running on port :8000")
	http.ListenAndServe(":8000", router)
}
