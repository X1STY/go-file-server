package main

import (
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go-file-server/internal/app"
)

func main() {
	godotenv.Load()
	app.StartApplication()
}
