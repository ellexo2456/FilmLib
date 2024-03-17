package usecase

import (
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
)

type actorUsecase struct {
	actorRepo domain.ActorRepository
}

func NewFilmUsecase(ar domain.ActorRepository) domain.ActorUsecase {
	return &actorUsecase{
		actorRepo: ar,
	}
}

func (u *actorUsecase) Add(actor domain.Actor) (int, error) {
	if actor.Name == "" {
		return 0, domain.ErrBadRequest
	}

	id, err := u.actorRepo.Insert(actor)
	if err != nil {
		logs.LogError(logs.Logger, "actor/usecase", "Add", err, err.Error())
		return 0, err
	}

	logs.Logger.Debug("actor/usecase Add:\n", id)
	return id, nil
}

func (u *actorUsecase) Remove(id int) error {
	if id <= 0 {
		return domain.ErrNotFound
	}

	err := u.actorRepo.Delete(id)
	if err != nil {
		logs.LogError(logs.Logger, "actor/usecase", "Remove", err, err.Error())
		return err
	}

	return nil
}

func (u *actorUsecase) Modify(newActor domain.Actor) (domain.Actor, error) {
	if newActor.ID <= 0 {
		return domain.Actor{}, domain.ErrNotFound
	}

	oldActor, err := u.actorRepo.SelectById(newActor.ID)
	if err != nil {
		logs.LogError(logs.Logger, "actor/usecase", "Modify", err, err.Error())
		return domain.Actor{}, err
	}
	logs.Logger.Debug("actor/usecase Modify old actor:\n", oldActor)

	newActor = getOldFields(newActor, oldActor)
	updatedActor, err := u.actorRepo.Update(newActor)
	if err != nil {
		logs.LogError(logs.Logger, "actor/usecase", "Modify", err, err.Error())
		return domain.Actor{}, err
	}
	logs.Logger.Debug("actor/usecase Modify updated actor:\n", updatedActor)

	return updatedActor, nil
}

func (u *actorUsecase) GetAll() ([]domain.Actor, error) {
	actors, err := u.actorRepo.SelectAll()
	if err != nil {
		logs.LogError(logs.Logger, "actor/usecase", "GetAll", err, err.Error())
		return nil, err
	}

	logs.Logger.Debug("actor/usecase GetAll actors:\n", actors)

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
