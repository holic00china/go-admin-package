package iris

import (
	// add iris adapter
	_ "go-admin-package/adapter/iris"
	"go-admin-package/modules/config"
	"go-admin-package/modules/language"
	"go-admin-package/plugins/admin/modules/table"

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

	"go-admin-package/template"
	"go-admin-package/template/chartjs"

	"net/http"
	"os"

	"go-admin-package/engine"
	"go-admin-package/plugins/admin"
	"go-admin-package/plugins/example"
	"go-admin-package/tests/tables"

	"github.com/kataras/iris/v12"
)

func internalHandler() http.Handler {
	app := iris.New()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)
	examplePlugin := example.NewExample()
	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	if err := app.Build(); err != nil {
		panic(err)
	}

	return app.Router
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	app := iris.New()

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(gens)

	examplePlugin := example.NewExample()
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
		AddPlugins(adminPlugin, examplePlugin).Use(app); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	if err := app.Build(); err != nil {
		panic(err)
	}

	return app.Router
}
