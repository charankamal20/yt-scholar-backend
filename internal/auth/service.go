package auth

import (
	"context"
	"database/sql"

	authStore "github.com/charankamal20/youtube-scholar-backend/database/repository/auth"
)

type AuthService interface {
	registerUser(c context.Context, user GoogleOAuthUser) *HTTPError
	loginUser(c context.Context, user GoogleOAuthUser) *HTTPError

	getUserInfo(c context.Context, userID string) (authStore.User, *HTTPError)
}

type Service struct {
	store authStore.Querier
}

func newAuthService(store authStore.Querier) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) loginUser(c context.Context, user GoogleOAuthUser) *HTTPError {
	dbUser, err := s.store.GetUserById(c, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}

		return ErrServer(err)
	}

	if dbUser.UserID == "" || dbUser.Email != user.Email {
		return ErrUserNotFound
	}

	return nil
}

func (s *Service) registerUser(c context.Context, user GoogleOAuthUser) *HTTPError {
	err := s.store.CreateUser(c, authStore.CreateUserParams{
		Email:      user.Email,
		UserID:     user.ID,
		Name:       user.Name,
		ProfilePic: user.Picture,
	})
	if err != nil {
		return ErrServer(err)
	}
	return nil
}

func (s *Service) getUserInfo(c context.Context, userID string) (authStore.User, *HTTPError) {
	dbUser, err := s.store.GetUserById(c, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return authStore.User{}, ErrUserNotFound
		}

		return authStore.User{}, ErrServer(err)
	}

	return dbUser, nil
}
