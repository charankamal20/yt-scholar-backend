// @title           Playlists Service API
// @version         1.0
// @description     Handles Playlist logic.
// @host            localhost:8080
// @BasePath        /playlist

package playlist

import (
	playlistStore "github.com/charankamal20/youtube-scholar-backend/database/repository/playlist"
	_ "github.com/charankamal20/youtube-scholar-backend/docs"
	"github.com/charankamal20/youtube-scholar-backend/internal/common"
)

type PlaylistServer struct {
	*common.Server
	service PlaylistService
}

func NewPlaylistServer(srv *common.Server, store playlistStore.Querier) *PlaylistServer {
	playlistServer := &PlaylistServer{
		Server:  srv,
		service: newPlaylistService(store),
	}

	playlistServer.registerRoutes()
	return playlistServer
}

func (a *PlaylistServer) registerRoutes() {
	playlistServer := a.Private.Group("/playlist")

	playlistServer.GET("", a.getAllUserPlaylists)
	playlistServer.POST("", a.addPlaylistHandler)

	playlistServer.GET("/:playlistId", a.getPlaylistForUser)
	playlistServer.DELETE("/:playlistId", a.deleteUserPlaylistHandler)
}
