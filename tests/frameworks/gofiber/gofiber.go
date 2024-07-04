package gofiber

import (
	// add fasthttp adapter
	ada "go-admin-package/adapter/gofiber"
	// add mysql driver
	_ "go-admin-package/modules/db/drivers/mysql"
	// add postgresql driver
	_ "go-admin-package/modules/db/drivers/postgres"
	// add sqlite driver
	_ "go-admin-package/modules/db/drivers/sqlite"
	// add mssql driver
	_ "go-admin-package/modules/db/drivers/mssql"
	// add adminlte ui theme
	"github.com/GoAdminGroup/themes/adminlte"

	"os"

	"go-admin-package/engine"
	"go-admin-package/modules/config"
	"go-admin-package/modules/language"
	"go-admin-package/plugins/admin"
	"go-admin-package/plugins/admin/modules/table"
	"go-admin-package/template"
	"go-admin-package/template/chartjs"
	"go-admin-package/tests/tables"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func internalHandler() fasthttp.RequestHandler {
	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
	})

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators).AddDisplayFilterXssJsFilter()
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return app.Handler()
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) fasthttp.RequestHandler {
	app := fiber.New(fiber.Config{
		ServerHeader: "Fiber",
	})

	eng := engine.Default()

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(&config.Config{
		Databases: dbs,
		UrlPrefix: "admin",
		Store: config.Store{
			Path:   "./uploads",
			Prefix: "uploads",
		},
		Language:    language.EN,
		IndexUrl:    "/",
		Debug:       true,
		ColorScheme: adminlte.ColorschemeSkinBlack,
	}).
		AddAdapter(new(ada.Gofiber)).
		AddGenerators(gens).
		Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return app.Handler()
}
