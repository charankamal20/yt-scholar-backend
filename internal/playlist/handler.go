package playlist

import (
	"errors"
	"net/http"

	playlistStore "github.com/charankamal20/youtube-scholar-backend/database/repository/playlist"
	"github.com/charankamal20/youtube-scholar-backend/internal/common"
	"github.com/gin-gonic/gin"
)

// addPlaylistHandler creates a new playlist for the authenticated user
// @Summary Create a new playlist
// @Description Creates a new playlist for the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Param playlist body playlistStore.AddNewPlaylistParams true "Playlist data"
// @Success 201 {object} map[string]string "Playlist created successfully"
// @Failure 400 {object} HTTPError "Bad request or validation error"
// @Failure 401 {object} HTTPError "Unauthorized"
// @Failure 500 {object} HTTPError "Internal server error"
// @Router /playlists [post]
func (s *PlaylistServer) addPlaylistHandler(c *gin.Context) {
	userId, exists := common.GetAuthUserID(c)
	if !exists {
		c.AbortWithError(http.StatusUnauthorized, errors.New("user id not available"))
		return
	}

	var req playlistStore.AddNewPlaylistParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Set the user ID from the authenticated user
	req.UserID = userId

	err := s.service.addNewPlaylist(c.Request.Context(), req)
	if err != nil {
		c.JSON(err.Code, err.GetGinError())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Playlist created successfully"})
}

// deleteUserPlaylistHandler deletes a playlist for the authenticated user
// @Summary Delete a user's playlist
// @Description Deletes a specific playlist belonging to the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Param playlistId path string true "Playlist ID"
// @Success 200 {object} map[string]string "Playlist deleted successfully"
// @Failure 400 {object} HTTPError "Bad request"
// @Failure 401 {object} HTTPError "Unauthorized"
// @Failure 404 {object} HTTPError "Playlist not found"
// @Failure 500 {object} HTTPError "Internal server error"
// @Router /playlists/{playlistId} [delete]
func (s *PlaylistServer) deleteUserPlaylistHandler(c *gin.Context) {
	userId, exists := common.GetAuthUserID(c)
	if !exists {
		c.AbortWithError(http.StatusUnauthorized, errors.New("user id not available"))
		return
	}

	playlistId := c.Param("playlistId")
	if playlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	arg := playlistStore.DeletePlaylistForUserParams{
		UserID:     userId,
		PlaylistID: playlistId,
	}

	err := s.service.deletePlaylistForUser(c.Request.Context(), arg)
	if err != nil {
		c.JSON(err.Code, err.GetGinError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted successfully"})
}

// getAllUserPlaylists retrieves all playlists for the authenticated user
// @Summary Get all user playlists
// @Description Retrieves all playlists belonging to the authenticated user
// @Tags playlists
// @Accept json
// @Produce json
// @Success 200 {array} playlistStore.Playlist "List of user playlists"
// @Failure 401 {object} HTTPError "Unauthorized"
// @Failure 500 {object} HTTPError "Internal server error"
// @Router /playlists [get]
func (s *PlaylistServer) getAllUserPlaylists(c *gin.Context) {
	userId, exists := common.GetAuthUserID(c)
	if !exists {
		c.AbortWithError(http.StatusUnauthorized, errors.New("user id not available"))
		return
	}

	playlists, err := s.service.getAllUserPlaylists(c.Request.Context(), userId)
	if err != nil {
		c.JSON(err.Code, err.GetGinError())
		return
	}

	c.JSON(http.StatusOK, playlists)
}

// getPlaylistById retrieves a playlist by its ID
// @Summary Get playlist by ID
// @Description Retrieves a specific playlist by its ID
// @Tags playlists
// @Accept json
// @Produce json
// @Param playlistId path string true "Playlist ID"
// @Success 200 {object} playlistStore.Playlist "Playlist data"
// @Failure 400 {object} HTTPError "Bad request"
// @Failure 404 {object} HTTPError "Playlist not found"
// @Failure 500 {object} HTTPError "Internal server error"
// @Router /playlists/{playlistId} [get]
func (s *PlaylistServer) getPlaylistById(c *gin.Context) {
	playlistId := c.Param("playlistId")
	if playlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	playlist, err := s.service.getPlaylistById(c.Request.Context(), playlistId)
	if err != nil {
		c.JSON(err.Code, err.GetGinError())
		return
	}

	c.JSON(http.StatusOK, playlist)
}

// getPlaylistForUser retrieves a specific playlist for a user
// @Summary Get playlist for user
// @Description Retrieves a specific playlist belonging to a user
// @Tags playlists
// @Accept json
// @Produce json
// @Param playlistId path string true "Playlist ID"
// @Param userId query string false "User ID (optional, defaults to authenticated user)"
// @Success 200 {object} playlistStore.Playlist "Playlist data"
// @Failure 400 {object} HTTPError "Bad request"
// @Failure 401 {object} HTTPError "Unauthorized"
// @Failure 404 {object} HTTPError "Playlist not found"
// @Failure 500 {object} HTTPError "Internal server error"
// @Router /playlists/{playlistId}/user [get]
func (s *PlaylistServer) getPlaylistForUser(c *gin.Context) {
	playlistId := c.Param("playlistId")
	if playlistId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Playlist ID is required"})
		return
	}

	// Get user ID from query parameter or use authenticated user
	userId := c.Query("userId")
	if userId == "" {
		authUserId, exists := common.GetAuthUserID(c)
		if !exists {
			c.AbortWithError(http.StatusUnauthorized, errors.New("user id not available"))
			return
		}
		userId = authUserId
	}

	arg := playlistStore.GetPlaylistForUserParams{
		UserID:     userId,
		PlaylistID: playlistId,
	}

	playlist, err := s.service.getPlaylistForUser(c.Request.Context(), arg)
	if err != nil {
		c.JSON(err.Code, err.GetGinError())
		return
	}

	c.JSON(http.StatusOK, playlist)
}
