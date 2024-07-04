package beego

import (
	// add beego adapter
	_ "go-admin-package/adapter/beego"
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

	"net/http"
	"os"

	"go-admin-package/engine"
	"go-admin-package/plugins/admin"
	"go-admin-package/plugins/example"
	"go-admin-package/template"
	"go-admin-package/template/chartjs"
	"go-admin-package/tests/tables"

	"github.com/astaxie/beego"
)

func internalHandler() http.Handler {

	app := beego.NewApp()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	examplePlugin := example.NewExample()

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 9087

	return app.Handlers
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {

	app := beego.NewApp()

	eng := engine.Default()
	adminPlugin := admin.NewAdmin(gens)

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
		AddPlugins(adminPlugin).Use(app); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 9087

	return app.Handlers
}
