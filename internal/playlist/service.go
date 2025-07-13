package playlist

import (
	"context"

	errors "github.com/charankamal20/youtube-scholar-backend/internal/common"

	playlistStore "github.com/charankamal20/youtube-scholar-backend/database/repository/playlist"
)

type PlaylistService interface {
	addNewPlaylist(ctx context.Context, arg playlistStore.AddNewPlaylistParams) *errors.HTTPError
	deletePlaylistForUser(ctx context.Context, arg playlistStore.DeletePlaylistForUserParams) *errors.HTTPError
	getAllUserPlaylists(ctx context.Context, userID string) ([]playlistStore.Playlist, *errors.HTTPError)
	getPlaylistById(ctx context.Context, playlistID string) (playlistStore.Playlist, *errors.HTTPError)
	getPlaylistForUser(ctx context.Context, arg playlistStore.GetPlaylistForUserParams) (playlistStore.Playlist, *errors.HTTPError)
}

type Service struct {
	playlistStore.Querier
}

func newPlaylistService(store playlistStore.Querier) *Service {
	return &Service{
		store,
	}
}

func (s *Service) addNewPlaylist(ctx context.Context, arg playlistStore.AddNewPlaylistParams) *errors.HTTPError {
	err := s.Querier.AddNewPlaylist(ctx, arg)
	if err != nil {
		return errors.ErrServer(err)
	}
	return nil
}

func (s *Service) deletePlaylistForUser(ctx context.Context, arg playlistStore.DeletePlaylistForUserParams) *errors.HTTPError {
	err := s.Querier.DeletePlaylistForUser(ctx, arg)
	if err != nil {
		return errors.ErrServer(err)
	}
	return nil
}

func (s *Service) getAllUserPlaylists(ctx context.Context, userID string) ([]playlistStore.Playlist, *errors.HTTPError) {
	playlists, err := s.Querier.GetAllUserPlaylists(ctx, userID)
	if err != nil {
		return nil, errors.ErrServer(err)
	}
	return playlists, nil
}

func (s *Service) getPlaylistById(ctx context.Context, playlistID string) (playlistStore.Playlist, *errors.HTTPError) {
	playlist, err := s.Querier.GetPlaylistById(ctx, playlistID)
	if err != nil {
		return playlistStore.Playlist{}, errors.ErrServer(err)
	}
	return playlist, nil
}

func (s *Service) getPlaylistForUser(ctx context.Context, arg playlistStore.GetPlaylistForUserParams) (playlistStore.Playlist, *errors.HTTPError) {
	playlist, err := s.Querier.GetPlaylistForUser(ctx, arg)
	if err != nil {
		return playlistStore.Playlist{}, errors.ErrServer(err)
	}
	return playlist, nil
}
