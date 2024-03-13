package app

import (
	"log"
	"net/http"
)

func StartServer() {
	mux := http.NewServeMux()

	err := http.ListenAndServe(serverPort, mux)
	if err != nil {
		log.Fatal("can`t start sever")
		//logs.LogFatal(logs.Logger, "auth", "main", err, err.Error())
	}

	//logs.Logger.Info("auth http server stopped")
}
