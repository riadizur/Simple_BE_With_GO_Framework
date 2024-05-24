package main

import (
	"log"
	"net/http"
	"os"

	"kubota-gasoline-api/internal/database"
	"kubota-gasoline-api/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initLog() {
	logFile, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(0) // disable default timestamp flags
	log.SetPrefix("")
}

func logFormatter() {
	log.SetFlags(0) // disable default flags
	log.SetOutput(log.Writer())
	log.SetPrefix("")
	log.SetFlags(log.LstdFlags)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	initLog()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in .env file")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is not set in .env file")
	}

	db, err := database.Connect(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	handlers.RegisterShiftHandlers(router, db)
	handlers.RegisterWebSocketHandlers(router)

	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
