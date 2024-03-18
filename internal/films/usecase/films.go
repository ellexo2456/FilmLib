package usecase

import (
	"cmp"
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
	"slices"
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
	if isEmpty(film) {
		return 0, domain.ErrBadRequest
	}

	id, err := u.filmsRepo.Insert(film)
	if err != nil {
		logs.LogError(logs.Logger, "films/usecase", "Add", err, err.Error())
		return 0, err
	}

	logs.Logger.Debug("films/usecase Add:", id)
	return id, nil
}

func (u *filmsUsecase) GetAll(titleDir, releaseDateDir domain.SortDirection) ([]domain.Film, error) {
	films, err := u.filmsRepo.SelectAll()
	if err != nil {
		logs.LogError(logs.Logger, "films/usecase", "GetAll", err, err.Error())
		return nil, err
	}

	logs.Logger.Debug("films/usecase GetAll films:\n", films)
	return sort(films, titleDir, releaseDateDir), nil
}

func (u *filmsUsecase) Search(searchStr string) ([]domain.Film, error) {
	if searchStr == "" {
		return nil, domain.ErrBadRequest
	}

	films, err := u.filmsRepo.Search(searchStr)
	if err != nil {
		logs.LogError(logs.Logger, "films/usecase", "Search", err, err.Error())
		return nil, err
	}
	logs.Logger.Debug("films/usecase Search films:", films)

	return films, nil
}

func (u *filmsUsecase) Remove(id int) error {
	if id <= 0 {
		return domain.ErrNotFound
	}

	err := u.filmsRepo.Delete(id)
	if err != nil {
		logs.LogError(logs.Logger, "films/usecase", "Remove", err, err.Error())
		return err
	}

	return nil
}

func (u *filmsUsecase) Modify(newFilm domain.Film) (domain.Film, error) {
	if newFilm.ID <= 0 {
		return domain.Film{}, domain.ErrNotFound
	}

	oldFilm, err := u.filmsRepo.SelectById(newFilm.ID)
	if err != nil {
		logs.LogError(logs.Logger, "films/usecase", "Modify", err, err.Error())
		return domain.Film{}, err
	}
	logs.Logger.Debug("films/usecase Modify old actor:\n", oldFilm)

	newFilm = getOldFields(newFilm, oldFilm)
	updatedActor, err := u.filmsRepo.Update(newFilm)
	if err != nil {
		logs.LogError(logs.Logger, "films/usecase", "Modify", err, err.Error())
		return domain.Film{}, err
	}
	logs.Logger.Debug("films/usecase Modify updated actor:\n", updatedActor)

	return updatedActor, nil
}

func getOldFields(newFilm, oldFilm domain.Film) domain.Film {
	if newFilm.Title == "" {
		newFilm.Title = oldFilm.Title
	}
	if newFilm.Description == "" {
		newFilm.Description = oldFilm.Description
	}
	if !newFilm.ReleaseDate.Valid {
		newFilm.ReleaseDate = oldFilm.ReleaseDate
	}
	if newFilm.Rating == 0 {
		newFilm.Rating = oldFilm.Rating
	}

	return newFilm
}

func sort(films []domain.Film, titleDir, releaseDateDir domain.SortDirection) []domain.Film {
	if titleDir == domain.Asc {
		slices.SortFunc(films, ascTitle())
		return films
	}

	if titleDir == domain.Desc {
		slices.SortFunc(films, descTitle())
		return films
	}

	if releaseDateDir == domain.Asc {
		slices.SortFunc(films, ascReleaseDate())
		return films
	}

	if releaseDateDir == domain.Desc {
		slices.SortFunc(films, descReleaseDate())
		return films
	}

	slices.SortFunc(films, descRating())
	return films
}

func ascTitle() func(a, b domain.Film) int {
	return func(a, b domain.Film) int {
		return cmp.Compare(a.Title, b.Title)
	}
}

func descTitle() func(a, b domain.Film) int {
	return func(a, b domain.Film) int {
		return cmp.Compare(b.Title, a.Title)
	}
}

func ascReleaseDate() func(a, b domain.Film) int {
	return func(a, b domain.Film) int {
		return a.ReleaseDate.Time.Compare(b.ReleaseDate.Time)
	}
}

func descReleaseDate() func(a, b domain.Film) int {
	return func(a, b domain.Film) int {
		return b.ReleaseDate.Time.Compare(a.ReleaseDate.Time)
	}
}

func descRating() func(a, b domain.Film) int {
	return func(a, b domain.Film) int {
		return cmp.Compare(b.Rating, a.Rating)
	}
}

func isEmpty(film domain.Film) bool {
	if film.Rating == 0 {
		return true
	}
	if !film.ReleaseDate.Valid || len(film.Actors) == 0 {
		return true
	}
	if film.Description == "" || film.Title == "" {
		return true
	}

	return false
}
