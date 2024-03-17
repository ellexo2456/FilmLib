package app

import (
	"context"
	"github.com/ellexo2456/FilmLib/internal/middleware"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"

	authhttp "github.com/ellexo2456/FilmLib/internal/auth/delivery/http"
	authpostgres "github.com/ellexo2456/FilmLib/internal/auth/repository/postgresql"
	authredis "github.com/ellexo2456/FilmLib/internal/auth/repository/redis"
	authusecase "github.com/ellexo2456/FilmLib/internal/auth/usecase"

	filmshttp "github.com/ellexo2456/FilmLib/internal/films/delivery/http"
	filmspostgres "github.com/ellexo2456/FilmLib/internal/films/repository/postgresql"
	filmsusecase "github.com/ellexo2456/FilmLib/internal/films/usecase"

	actorshttp "github.com/ellexo2456/FilmLib/internal/actors/delivery/http"
	actorspostgres "github.com/ellexo2456/FilmLib/internal/actors/repository/postgresql"
	actorsusecase "github.com/ellexo2456/FilmLib/internal/actors/usecase"

	_ "github.com/ellexo2456/FilmLib/docs"
	"github.com/ellexo2456/FilmLib/internal/connectors/postgres"
	"github.com/ellexo2456/FilmLib/internal/connectors/redis"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
)

// 61 73 76 31
func StartServer() {
	err := godotenv.Load()

	mux := http.NewServeMux()

	dbParams := postgres.GetDbParams()
	ctx := context.Background()
	pc := postgres.Connect(ctx, dbParams)
	defer pc.Close()
	rc := redis.Connect()
	defer rc.Close()

	sr := authredis.NewSessionRedisRepository(rc)
	ar := authpostgres.NewAuthPostgresqlRepository(pc, ctx)
	acr := actorspostgres.NewActorsPostgresqlRepository(pc, ctx)
	fr := filmspostgres.NewFilmsPostgresqlRepository(pc, ctx)

	au := authusecase.NewAuthUsecase(ar, sr)
	acu := actorsusecase.NewActorsUsecase(acr)
	fu := filmsusecase.NewFilmsUsecase(fr)

	authMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	authhttp.NewAuthHandler(authMux, au)
	actorshttp.NewActorsHandler(apiMux, acu)
	filmshttp.NewFilmHandler(apiMux, fu)
	mux.HandleFunc("/swagger/*", httpSwagger.WrapHandler)

	mw := middleware.NewAuth(au)

	mux.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", authMux))
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", mw.IsAuth(apiMux)))

	port := ":" + os.Getenv("HTTP_SERVER_PORT")
	logs.Logger.Info("start listening on port" + port)
	err = http.ListenAndServe(port, mux)
	if err != nil {
		logs.LogFatal(logs.Logger, "app", "main", err, err.Error())
	}

	logs.Logger.Info("server stopped")
}
