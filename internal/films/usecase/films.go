package usecase

import (
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
)

type filmsUsecase struct {
	filmsRepo domain.FilmsRepository
}

func NewFilmsUsecase(fr domain.FilmsRepository) domain.FilmsUsecase {
	return &filmsUsecase{
		filmsRepo: fr,
	}
}

func (u *filmsUsecase) Add(film domain.Film) (int, error) {
	id, err := u.filmsRepo.Insert(film)
	if err != nil {
		logs.LogError(logs.Logger, "films/usecase", "Add", err, err.Error())
		return 0, err
	}

	logs.Logger.Debug("films/usecase Add:", id)
	return id, nil
}
