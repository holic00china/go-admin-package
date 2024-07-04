package example

import (
	"go-admin-package/context"
	"go-admin-package/modules/auth"
	"go-admin-package/modules/db"
	"go-admin-package/modules/service"
)

func (e *Example) initRouter(prefix string, srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(prefix)
	route.GET("/example", auth.Middleware(db.GetConnection(srv)), e.TestHandler)

	return app
}
