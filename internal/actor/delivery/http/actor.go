package http

import (
	"encoding/json"
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	ActorUsecase domain.ActorUsecase
}

func NewActorHandler(mux *http.ServeMux, au domain.ActorUsecase) {
	handler := &ActorHandler{
		ActorUsecase: au,
	}

	mux.HandleFunc("POST /actor", handler.AddActor)
	mux.HandleFunc("DELETE /actor/{id}", handler.DeleteActor)
	mux.HandleFunc("PUT /actor", handler.ModifyActor)
	mux.HandleFunc("GET /actor", handler.GetActors)

}

// AddActor godoc
//
//	@Summary		Adds a new actor.
//	@Description	Adds a new actor with the provided data.
//	@Tags			Actor
//	@Param			body	body		domain.Actor	true	"actor to add"
//	@Produce		json
//	@Success		200		{json}	object{body=object{id=int}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		403		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//	@Router			/api/v1/actor [post]
func (h *ActorHandler) AddActor(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "actor/http", "AddActor", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("AddActor session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "actor/http", "AddActor", errors.New("forbidden"), "invalid role")
	}

	var actor domain.Actor

	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "actor/http", "AddActor", err, err.Error())
		return
	}
	logs.Logger.Debug("AddActor actor:\n", actor)
	defer domain.CloseAndAlert(r.Body, "actor/http", "AddActor")

	id, err := h.ActorUsecase.Add(actor)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actor/http", "AddActor", err, err.Error())
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
//	@Description	Deletes an actor by id.
//	@Tags			Actor
//	@Param			id	path	int	true	"Actor id"
//	@Produce		json
//	@Success		204
//	@Failure		400		{json}	object{err=string}
//	@Failure		403		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//	@Router			/api/v1/actor/{id} [delete]
func (h *ActorHandler) DeleteActor(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "actor/http", "DeleteActor", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("DeleteActor session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "actor/http", "DeleteActor", errors.New("forbidden"), "invalid role")
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "actor/http", "DeleteActor", err, err.Error())
		return
	}
	logs.Logger.Debug("DeleteActor id:\n", id)

	err = h.ActorUsecase.Remove(id)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actor/http", "DeleteActor", err, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ModifyActor godoc
//
//	@Summary		Modify an actor.
//	@Description	Modify an actor by id and retrieves new actor.
//	@Tags			Actor
//	@Param			body	body		domain.Actor	true	"actor to add"
//	@Produce		json
//	@Success		200		{json}	object{body=object{actor=domain.Actor}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		403		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//	@Router			/api/v1/actor [put]
func (h *ActorHandler) ModifyActor(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "actor/http", "ModifyActor", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("DeleteActor session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "actor/http", "ModifyActor", errors.New("forbidden"), "invalid role")
	}

	var actor domain.Actor
	err := json.NewDecoder(r.Body).Decode(&actor)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "actor/http", "ModifyActor", err, err.Error())
		return
	}
	logs.Logger.Debug("ModifyActor new actor:\n", actor)
	defer domain.CloseAndAlert(r.Body, "actor/http", "ModifyActor")

	actor, err = h.ActorUsecase.Modify(actor)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actor/http", "ModifyActor", err, err.Error())
		return
	}
	logs.Logger.Debug("ModifyActor updated actor:\n", actor)

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
//	@Summary		Gets an actors.
//	@Description	Gets all actors with related films.
//	@Tags			Actor
//	@Produce		json
//	@Success		200		{json}	object{body=object{actors=[]domain.Actor}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//	@Router			/api/v1/actor [get]
func (h *ActorHandler) GetActors(w http.ResponseWriter, r *http.Request) {
	actors, err := h.ActorUsecase.GetAll()
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "actor/http", "GetActors", err, err.Error())
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
