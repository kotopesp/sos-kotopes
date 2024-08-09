package bot

import (
	"github.com/kotopesp/tgbot/internal/core"
	th "github.com/mymmrac/telego/telegohandler"
)

type Router struct {
	app         *th.BotHandler
	authService core.AuthService
}

func NewRouter(
	app *th.BotHandler,
	authService core.AuthService,
) {
	router := &Router{
		app:         app,
		authService: authService,
	}

	router.initRoutes()
}

func (r *Router) initRoutes() {
	r.app.Handle(r.start, th.TextEqual("/start"))
	r.app.Handle(r.back, th.CallbackDataEqual("back"))
	r.app.Handle(r.authorize, th.CallbackDataEqual("authorize"))
}
