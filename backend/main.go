package main

import (
	"log"
	"net/http"

	cors "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/leoldding/odd-one-out/handlers"
	"github.com/leoldding/odd-one-out/pubsub"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	publisher := pubsub.NewPublisher()

	handlers.RegisterRoomHandlers(router, publisher)
	handlers.RegisterGameHandlers(router, publisher)

	headersOk := cors.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := cors.AllowedOrigins([]string{"http://localhost:5173"})
	methodsOk := cors.AllowedMethods([]string{"GET", "POST"})

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", cors.CORS(originsOk, headersOk, methodsOk)(router)))
}
