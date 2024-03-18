package app

import (
	"context"
	"github.com/ellexo2456/FilmLib/internal/middleware"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"

	auth_http "github.com/ellexo2456/FilmLib/internal/auth/delivery/http"
	auth_postgres "github.com/ellexo2456/FilmLib/internal/auth/repository/postgresql"
	auth_redis "github.com/ellexo2456/FilmLib/internal/auth/repository/redis"
	auth_usecase "github.com/ellexo2456/FilmLib/internal/auth/usecase"

	films_http "github.com/ellexo2456/FilmLib/internal/films/delivery/http"
	films_postgres "github.com/ellexo2456/FilmLib/internal/films/repository/postgresql"
	films_usecase "github.com/ellexo2456/FilmLib/internal/films/usecase"

	actors_http "github.com/ellexo2456/FilmLib/internal/actors/delivery/http"
	actors_postgres "github.com/ellexo2456/FilmLib/internal/actors/repository/postgresql"
	actors_usecase "github.com/ellexo2456/FilmLib/internal/actors/usecase"

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

	sr := auth_redis.NewSessionRedisRepository(rc)
	ar := auth_postgres.NewAuthPostgresqlRepository(pc, ctx)
	acr := actors_postgres.NewActorsPostgresqlRepository(pc, ctx)
	fr := films_postgres.NewFilmsPostgresqlRepository(pc, ctx)

	au := auth_usecase.NewAuthUsecase(ar, sr)
	acu := actors_usecase.NewActorsUsecase(acr)
	fu := films_usecase.NewFilmsUsecase(fr)

	authMux := http.NewServeMux()
	apiMux := http.NewServeMux()

	auth_http.NewAuthHandler(authMux, au)
	actors_http.NewActorsHandler(apiMux, acu)
	films_http.NewFilmsHandler(apiMux, fu)
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
