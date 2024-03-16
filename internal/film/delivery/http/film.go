package http

import (
	"encoding/json"
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
	"net/http"
)

type FilmHandler struct {
	FilmUsecase domain.FilmUsecase
}

func NewFilmHandler(mux *http.ServeMux, fu domain.FilmUsecase) {
	handler := &FilmHandler{
		FilmUsecase: fu,
	}

	mux.HandleFunc("POST /film", handler.AddFilm)

}

// AddFilm godoc
//
//	@Summary		Adds new film.
//	@Description	Adds new film with provided data.
//	@Tags			Film
//	@Param			body	body		domain.Film	true	"film to add"
//	@Produce		json
//	@Success		200		{json}	object{body=object{id=int}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//	@Router			/api/v1/film [post]
func (h *FilmHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
	sc, ok := r.Context().Value(domain.SessionContextKey).(domain.SessionContext)
	if !ok {
		domain.WriteError(w, "can`t find user", http.StatusInternalServerError)
		logs.LogError(logs.Logger, "film/http", "AddFilm", errors.New("can`t find user"), "can`t find user")
	}
	logs.Logger.Debug("AddFilm session context\n: ", sc)

	if sc.Role != domain.Moder {
		domain.WriteError(w, "forbidden", http.StatusForbidden)
		logs.LogError(logs.Logger, "film/http", "AddFilm", errors.New("forbidden"), "invalid role")
	}

	var film domain.Film

	err := json.NewDecoder(r.Body).Decode(&film)
	if err != nil {
		domain.WriteError(w, err.Error(), http.StatusBadRequest)
		logs.LogError(logs.Logger, "film/http", "AddFilm", err, err.Error())
		return
	}
	logs.Logger.Debug("AddFilm film:\n", film)
	defer domain.CloseAndAlert(r.Body, "film/http", "AddFilm")

	id, err := h.FilmUsecase.Add(film)
	if err != nil {
		domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
		logs.LogError(logs.Logger, "film/http", "AddFilm", err, err.Error())
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
