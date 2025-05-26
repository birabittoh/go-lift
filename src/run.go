package src

import (
	"log"
	"net/http"
	"os"

	"github.com/birabittoh/go-lift/src/api"
	"github.com/birabittoh/go-lift/src/database"
	"github.com/joho/godotenv"
)

func getEnv(key, def string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return def
}

func Run() (err error) {
	godotenv.Load()

	db, err := database.InitializeDB()
	if err != nil {
		return
	}

	listenAddress := getEnv("APP_LISTEN_ADDRESS", ":3000")

	mux := api.GetServeMux(db)

	log.Println("Listening at", listenAddress)
	err = http.ListenAndServe(listenAddress, mux)
	return
}
