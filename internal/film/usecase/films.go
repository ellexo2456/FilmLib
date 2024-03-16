package usecase

import (
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
)

type filmUsecase struct {
	filmRepo domain.FilmRepository
}

func NewFilmUsecase(fr domain.FilmRepository) domain.FilmUsecase {
	return &filmUsecase{
		filmRepo: fr,
	}
}

func (u *filmUsecase) Add(film domain.Film) (int, error) {
	id, err := u.filmRepo.Insert(film)
	if err != nil {
		logs.LogError(logs.Logger, "film/usecase", "Add", err, err.Error())
		return 0, err
	}

	logs.Logger.Debug("film/usecase Add:", id)
	return id, nil
}
