package gocartel

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func TestClient() BigCartelClient {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	opts := ClientOpts{
		BaseURL:    "https://api.bigcartel.com/v1",
		UserAgent:  os.Getenv("USER_AGENT"),
		BasicAuth:  os.Getenv("BASIC_AUTH"),
		HTTPClient: &http.Client{},
	}
	return NewClient(opts)
}

func InternalStoreID() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	return os.Getenv("STORE_ID")
}
