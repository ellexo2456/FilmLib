package http

import (
	"encoding/json"
	"net/http"

	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
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
//	@Produce		json`
//	@Success		200		{json}	object{body=object{id=int}}
//	@Failure		400		{json}	object{err=string}
//	@Failure		404		{json}	object{err=string}
//	@Failure		500		{json}	object{err=string}
//
//	@Router			/api/v1/film [post]
func (h *FilmHandler) AddFilm(w http.ResponseWriter, r *http.Request) {
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
