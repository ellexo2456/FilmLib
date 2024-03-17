package http

import (
	"encoding/json"
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
	"net/http"
	"strconv"
)

type FilmsHandler struct {
	FilmsUsecase domain.FilmsUsecase
}

func NewFilmHandler(mux *http.ServeMux, fu domain.FilmsUsecase) {
	handler := &FilmsHandler{
		FilmsUsecase: fu,
	}

	mux.HandleFunc("POST /films", handler.AddFilm)
	mux.HandleFunc("GET /films", handler.GetFilms)
	mux.HandleFunc("GET /films/search", handler.Search)
	mux.HandleFunc("DELETE /films/{id}", handler.DeleteFilm)
	mux.HandleFunc("PUT /films", handler.ModifyFilm)

}

// AddFilm godoc
//
//	@Summary		Adds a new film.
//	@Description	Adds a new film with provided data.
//	@Tags			Films
//	@Param			body	body	domain.Film	true	"film to add"
//	@Produce		json
//	@Success		200	{object}	object{body=object{id=int}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		404	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/films [post]
func (h *FilmsHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "films/http", "AddFilm", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("AddFilm session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "films/http", "AddFilm", errors.New("forbidden"), "invalid role")
	}

	var film domain.Film

	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films/http", "AddFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("AddFilm film:\n", film)
	defer domain.CloseAndAlert(r.Body, "films/http", "AddFilm")

	id, err := h.FilmsUsecase.Add(film)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "films/http", "AddFilm", err, err.Error())
		return
	}

	logs.Logger.Debug("AddFilm film id:\n", id)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"id": id,
		},
		http.StatusOK,
	)
}

// GetFilms godoc
//
//	@Summary		Gets films.
//	@Description	Gets all films descending sorted by rating (by default). Only one sort can be applied at a time. If several are applied, the priority is as follows: title, releaseDate, rating (by default).
//	@Tags			Films
//	@Param			sortTitle		query	domain.SortDirection	false	"Direction of title sort. Sorting wont be applied if param isnt specified."
//	@Param			sortReleaseDate	query	domain.SortDirection	false	"Direction of release date sort. Sorting wont be applied if param isnt specified."
//	@Produce		json
//	@Success		200	{object}	object{body=object{films=[]domain.Film}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/films [get]
func (h *FilmsHandler) GetFilms(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	titleDir := (domain.SortDirection)(queryParams.Get(domain.TitleParam))
	releaseDateDir := (domain.SortDirection)(queryParams.Get(domain.ReleaseDateParam))

	films, err := h.FilmsUsecase.GetAll(titleDir, releaseDateDir)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "films/http", "GetFilms", err, err.Error())
		return
	}

	logs.Logger.Debug("GetFilms films:\n", films)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"films": films,
		},
		http.StatusOK,
	)
}

// Search godoc
//	@Summary		Searches films
//	@Description	Searches films by parts of its titles and parts of films names.
//	@Tags			Films
//	@Produce		json
//	@Param			searchStr	query		string	true	"The string to be searched for"
//	@Success		200			{object}	object{body=object{films=[]domain.Film}}
//	@Failure		400			{object}	object{err=string}
//	@Failure		404			{object}	object{err=string}
//	@Failure		500			{object}	object{err=string}
//	@Router			/api/v1/films/search [get]
func (h *FilmsHandler) Search(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	searchStr := queryParams.Get(domain.SearchParam)

	films, err := h.FilmsUsecase.Search(searchStr)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "films/http", "Search", err, err.Error())
		return
	}

	logs.Logger.Debug("films/http Search films:\n", films)
	domain.WriteResponse(
		w,
		map[string]interface{}{
			"films": films,
		},
		http.StatusOK,
	)
}

// DeleteFilm godoc
//
//	@Summary		Deletes a film.
//	@Description	Deletes a film by id with all its relations with actors.
//	@Tags			Films
//	@Param			id	path	int	true	"Film id"
//	@Produce		json
//	@Success		204
//	@Failure		400	{object}	object{err=string}
//	@Failure		403	{object}	object{err=string}
//	@Failure		404	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/films/{id} [delete]
func (h *FilmsHandler) DeleteFilm(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "films/http", "DeleteFilm", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("DeleteFilm session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "films/http", "DeleteFilm", errors.New("forbidden"), "invalid role")
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films/http", "DeleteFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("DeleteFilm id:\n", id)

	err = h.FilmsUsecase.Remove(id)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "films/http", "DeleteFilm", err, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ModifyFilm godoc
//
//	@Summary		Modify a film.
//	@Description	Modify a film by id and retrieves a new film.
//	@Tags			Films
//	@Param			body	body	domain.Film	true	"Film to modify"
//	@Produce		json
//	@Success		200	{object}	object{body=object{film=domain.Film}}
//	@Failure		400	{object}	object{err=string}
//	@Failure		403	{object}	object{err=string}
//	@Failure		500	{object}	object{err=string}
//	@Router			/api/v1/films [put]
func (h *FilmsHandler) ModifyFilm(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "films/http", "ModifyFilm", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("DeleteFilm session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "films/http", "ModifyFilm", errors.New("forbidden"), "invalid role")
	}

	var film domain.Film
	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "films/http", "ModifyFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("ModifyFilm new film:\n", film)
	defer domain.CloseAndAlert(r.Body, "films/http", "ModifyFilm")

	film, err = h.FilmsUsecase.Modify(film)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "films/http", "ModifyFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("ModifyFilm updated film:\n", film)

	domain.WriteResponse(
		w,
		map[string]interface{}{
			"film": film,
		},
		http.StatusOK,
	)
}
