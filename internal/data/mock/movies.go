package mock

import (
	"time"

	"github.com/IfedayoAwe/greenlight/internal/data"
)

var mockMovie = &data.Movie{
	ID:        1,
	UserID:    1,
	CreatedAt: time.Now(),
	Title:     "Test Movie",
	Year:      2003,
	Runtime:   2000,
	Genres:    []string{"Comedy", "Drama"},
	Version:   1,
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *data.Movie) error {
	return nil
}

func (m MockMovieModel) Get(id int64) (*data.Movie, error) {
	switch id {
	case 1:
		return mockMovie, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m MockMovieModel) Update(movie *data.Movie) error {
	switch movie.ID {
	case 1:
		return nil
	default:
		return data.ErrRecordNotFound
	}
}

func (m MockMovieModel) Delete(id int64) error {
	switch id {
	case 1:
		return nil
	default:
		return data.ErrRecordNotFound
	}
}

func (m MockMovieModel) GetAll(title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
	movies := []*data.Movie{}
	metadata := data.Metadata{}
	return movies, metadata, nil
}
