package usecase

import (
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
)

type actorsUsecase struct {
	actorsRepo domain.ActorsRepository
}

func NewActorsUsecase(ar domain.ActorsRepository) domain.ActorsUsecase {
	return &actorsUsecase{
		actorsRepo: ar,
	}
}

func (u *actorsUsecase) Add(actor domain.Actor) (int, error) {
	if actor.Name == "" {
		return 0, domain.ErrBadRequest
	}

	id, err := u.actorsRepo.Insert(actor)
	if err != nil {
		logs.LogError(logs.Logger, "actors/usecase", "Add", err, err.Error())
		return 0, err
	}

	logs.Logger.Debug("actors/usecase Add:\n", id)
	return id, nil
}

func (u *actorsUsecase) Remove(id int) error {
	if id <= 0 {
		return domain.ErrNotFound
	}

	err := u.actorsRepo.Delete(id)
	if err != nil {
		logs.LogError(logs.Logger, "actors/usecase", "Remove", err, err.Error())
		return err
	}

	return nil
}

func (u *actorsUsecase) Modify(newActor domain.Actor) (domain.Actor, error) {
	if newActor.ID <= 0 {
		return domain.Actor{}, domain.ErrNotFound
	}

	oldActor, err := u.actorsRepo.SelectById(newActor.ID)
	if err != nil {
		logs.LogError(logs.Logger, "actors/usecase", "Modify", err, err.Error())
		return domain.Actor{}, err
	}
	logs.Logger.Debug("actors/usecase Modify old actor:\n", oldActor)

	newActor = getOldFields(newActor, oldActor)
	updatedActor, err := u.actorsRepo.Update(newActor)
	if err != nil {
		logs.LogError(logs.Logger, "actors/usecase", "Modify", err, err.Error())
		return domain.Actor{}, err
	}
	logs.Logger.Debug("actors/usecase Modify updated actor:\n", updatedActor)

	return updatedActor, nil
}

func (u *actorsUsecase) GetAll() ([]domain.Actor, error) {
	actors, err := u.actorsRepo.SelectAll()
	if err != nil {
		logs.LogError(logs.Logger, "actors/usecase", "GetAll", err, err.Error())
		return nil, err
	}

	logs.Logger.Debug("actors/usecase GetAll actors:\n", actors)

	return actors, nil

}

func getOldFields(newActor, oldActor domain.Actor) domain.Actor {
	if newActor.Name == "" {
		newActor.Name = oldActor.Name
	}
	if newActor.Sex == "" {
		newActor.Sex = oldActor.Sex
	}
	if !newActor.Birthdate.Valid {
		newActor.Birthdate = oldActor.Birthdate
	}
	return newActor
}
