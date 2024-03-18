package main

import "github.com/ellexo2456/FilmLib/internal/app"

//	@title			FilmLib API
//	@version		1.0
//	@description	API of the FilmLib project

//	@contact.name	Alex Chinaev
//	@contact.url	https://vk.com/l.chinaev
//	@contact.email	ax.chinaev@yandex.ru

//	@license.name	AS IS (NO WARRANTY)

// @host		localhost:3000
// @schemes	http
// @BasePath	/
func main() {
	app.StartServer()
}
