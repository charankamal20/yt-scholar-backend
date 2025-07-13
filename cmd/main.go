// @title           Youtube Scholar API
// @version         1.0
// @description     Backend for youtube scholar.
// @host            localhost:8080
// @BasePath        /

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	authStore "github.com/charankamal20/youtube-scholar-backend/database/repository/auth"
	playlistStore "github.com/charankamal20/youtube-scholar-backend/database/repository/playlist"
	"github.com/charankamal20/youtube-scholar-backend/internal/auth"
	"github.com/charankamal20/youtube-scholar-backend/internal/common"
	"github.com/charankamal20/youtube-scholar-backend/internal/playlist"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("emv not found")
	}

}

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		panic("DB_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	defer db.Close()

	server := common.NewServer()

	authStr := authStore.New(db)
	auth.NewAuth(server, authStr)

	playlistStore := playlistStore.New(db)
	playlist.NewPlaylistServer(server, playlistStore)

	server.Run()
}
