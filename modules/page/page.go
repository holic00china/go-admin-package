// Copyright 2019 GoAdmin Core Team. All rights reserved.
// Use of this source code is governed by a Apache-2.0 style
// license that can be found in the LICENSE file.

package page

import (
	"bytes"

	"go-admin-package/context"
	"go-admin-package/modules/config"
	"go-admin-package/modules/db"
	"go-admin-package/modules/logger"
	"go-admin-package/modules/menu"
	"go-admin-package/plugins/admin/models"
	"go-admin-package/template"
	"go-admin-package/template/types"
)

// SetPageContent set and return the panel of page content.
func SetPageContent(ctx *context.Context, user models.UserModel, c func(ctx interface{}) (types.Panel, error), conn db.Connection) {

	panel, err := c(ctx)

	if err != nil {
		logger.ErrorCtx(ctx, "SetPageContent %+v", err)
		panel = template.WarningPanel(ctx, err.Error())
	}

	tmpl, tmplName := template.Get(ctx, config.GetTheme()).GetTemplate(ctx.IsPjax())

	ctx.AddHeader("Content-Type", "text/html; charset=utf-8")

	buf := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(buf, tmplName, types.NewPage(ctx, &types.NewPageParam{
		User:         user,
		Menu:         menu.GetGlobalMenu(user, conn, ctx.Lang()).SetActiveClass(config.URLRemovePrefix(ctx.Path())),
		Panel:        panel.GetContent(config.IsProductionEnvironment()),
		Assets:       template.GetComponentAssetImportHTML(ctx),
		TmplHeadHTML: template.Default(ctx).GetHeadHTML(),
		TmplFootJS:   template.Default(ctx).GetFootJS(),
		Iframe:       ctx.IsIframe(),
	}))
	if err != nil {
		logger.ErrorCtx(ctx, "SetPageContent %+v", err)
	}
	ctx.WriteString(buf.String())
}
