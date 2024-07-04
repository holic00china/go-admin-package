package gf

import (
	"net/http"
	"os"

	// add gf adapter
	_ "go-admin-package/adapter/gf"
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

	"go-admin-package/engine"
	"go-admin-package/modules/config"
	"go-admin-package/modules/language"
	"go-admin-package/plugins/admin"
	"go-admin-package/plugins/admin/modules/table"
	"go-admin-package/template"
	"go-admin-package/template/chartjs"
	"go-admin-package/tests/tables"

	"github.com/gogf/gf/frame/g"
)

func internalHandler() http.Handler {
	s := g.Server(8103)

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators).AddDisplayFilterXssJsFilter()

	template.AddComp(chartjs.NewChart())

	adminPlugin.AddGenerator("user", tables.GetUserTable)

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin).
		Use(s); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	return s
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {

	s := g.Server(8103)

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
		AddPlugins(adminPlugin).Use(s); err != nil {
		panic(err)
	}

	template.AddComp(chartjs.NewChart())

	eng.HTML("GET", "/admin", tables.GetContent)

	return s
}
