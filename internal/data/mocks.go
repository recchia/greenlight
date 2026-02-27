package data

import (
	"time"
)

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {
	return nil
}
func (m MockMovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	return &Movie{
		ID:      id,
		Title:   "Test Movie",
		Year:    2024,
		Runtime: 120,
		Genres:  []string{"action", "comedy"},
		Version: 1,
	}, nil
}
func (m MockMovieModel) Update(movie *Movie) error {
	return nil
}
func (m MockMovieModel) Delete(id int64) error {
	return nil
}
func (m MockMovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error) {
	return nil, Metadata{}, nil
}

type MockPermissionModel struct{}

func (m MockPermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	return Permissions{"movies:read", "movies:write"}, nil
}
func (m MockPermissionModel) AddForUser(userID int64, codes ...string) error {
	return nil
}

type MockTokenModel struct{}

func (m MockTokenModel) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	return &Token{
		Plaintext: "token26charslong1234567890",
		UserID:    userId,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}, nil
}
func (m MockTokenModel) Insert(token *Token) error {
	return nil
}
func (m MockTokenModel) DeleteAllForUser(userId int64, scope string) error {
	return nil
}

type MockUserModel struct{}

func (m MockUserModel) Insert(user *User) error {
	user.ID = 1
	user.CreatedAt = time.Now()
	user.Activated = false
	return nil
}

func (m MockUserModel) GetByEmail(email string) (*User, error) {
	if email == "test@example.com" {
		return &User{
			ID:        1,
			CreatedAt: time.Now(),
			Name:      "Test User",
			Email:     "test@example.com",
			Activated: true,
		}, nil
	}
	return nil, ErrRecordNotFound
}
func (m MockUserModel) Update(user *User) error {
	return nil
}
func (m MockUserModel) GetForToken(tokenScope string, tokenPlaintext string) (*User, error) {
	return &User{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "Test User",
		Email:     "test@example.com",
		Activated: true,
	}, nil
}

func NewMockModels() Models {
	return Models{
		Movies:      MockMovieModel{},
		Permissions: MockPermissionModel{},
		Tokens:      MockTokenModel{},
		Users:       MockUserModel{},
	}
}
