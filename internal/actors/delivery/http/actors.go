package http

import (
	"encoding/json"
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
	"net/http"
	"strconv"
)

type ActorsHandler struct {
	ActorsUsecase domain.ActorsUsecase
}

func NewActorsHandler(mux *http.ServeMux, au domain.ActorsUsecase) {
	handler := &ActorsHandler{
		ActorsUsecase: au,
	}

	mux.HandleFunc("POST /actors", handler.AddActor)
	mux.HandleFunc("DELETE /actors/{id}", handler.DeleteActor)
	mux.HandleFunc("PUT /actors", handler.ModifyActor)
	mux.HandleFunc("GET /actors", handler.GetActors)

}

// AddActor godoc
//
//	@Summary		Adds a new actor.
//	@Description	Adds a new actor with the provided data.
//	@Tags			Actors
//	@Param			body	body	domain.ActorSWG	true	"actor to add"
//	@Produce		json
//	@Success		200	{object}	object{body=object{id=int}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		403	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/actors [post]
func (h *ActorsHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "actors/http", "AddActor", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("AddActor session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "actors/http", "AddActor", errors.New("forbidden"), "invalid role")
	}

	var actor domain.Actor

	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "actors/http", "AddActor", err, err.Error())
		return
	}
	logs.Logger.Debug("AddActor actor:\n", actor)
	defer domain.CloseAndAlert(r.Body, "actors/http", "AddActor")

	id, err := h.ActorsUsecase.Add(actor)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actors/http", "AddActor", err, err.Error())
		return
	}
	logs.Logger.Debug("AddActor actor id:\n", id)

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"id": id,
		},
		http.StatusOK,
	)
}

// DeleteActor godoc
//
//	@Summary		Deletes an actor.
//	@Description	Deletes an actor by id with all its relations with films.
//	@Tags			Actors
//	@Param			id	path	int	true	"Actor id"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	object{err=string}
//	@Failure		403	{object}	object{err=string}
//	@Failure		404	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/actors/{id} [delete]
func (h *ActorsHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "actors/http", "DeleteFilm", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("DeleteFilm session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "actors/http", "DeleteFilm", errors.New("forbidden"), "invalid role")
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "actors/http", "DeleteFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("DeleteFilm id:\n", id)

	err = h.ActorsUsecase.Remove(id)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actors/http", "DeleteFilm", err, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ModifyActor godoc
//
//	@Summary		Modify an actor.
//	@Description	Modify an actor by id and retrieves a new actor.
//	@Tags			Actors
//	@Param			body	body	domain.Actor	true	"Actor to modify"
//	@Produce		json
//	@Success		200	{object}	object{body=object{actors=domain.Actor}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		403	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/actors [put]
func (h *ActorsHandler) ModifyActor(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "actors/http", "ModifyFilm", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("DeleteFilm session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "actors/http", "ModifyFilm", errors.New("forbidden"), "invalid role")
	}

	var actor domain.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "actors/http", "ModifyFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("ModifyFilm new actor:\n", actor)
	defer domain.CloseAndAlert(r.Body, "actors/http", "ModifyFilm")

	actor, err = h.ActorsUsecase.Modify(actor)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actors/http", "ModifyFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("ModifyFilm updated actor:\n", actor)

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"actor": actor,
		},
		http.StatusOK,
	)
}

// GetActors godoc
//
//	@Summary		Gets actors.
//	@Description	Gets all actors with related films.
//	@Tags			Actors
//	@Produce		json
//	@Success		200	{object}	object{body=object{actors=[]domain.Actor}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/actors [get]
func (h *ActorsHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	actors, err := h.ActorsUsecase.GetAll()
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actors/http", "GetActors", err, err.Error())
		return
	}

	logs.Logger.Debug("GetActors actors:\n", actors)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"actors": actors,
		},
		http.StatusOK,
	)
}
