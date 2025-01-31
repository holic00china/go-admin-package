package controller

import (
	"bytes"
	template2 "html/template"
	"net/http"
	"net/url"

	"go-admin-package/context"
	"go-admin-package/modules/auth"
	"go-admin-package/modules/config"
	"go-admin-package/modules/db"
	"go-admin-package/modules/logger"
	"go-admin-package/modules/system"
	"go-admin-package/plugins/admin/models"
	"go-admin-package/plugins/admin/modules/captcha"
	"go-admin-package/plugins/admin/modules/response"
	"go-admin-package/template"
	"go-admin-package/template/types"
)

// Auth check the input password and username for authentication.
func (h *Handler) Auth(ctx *context.Context) {

	var (
		user     models.UserModel
		ok       bool
		errMsg   = "fail"
		s, exist = h.services.GetOrNot(auth.ServiceKey)
	)

	if capDriver, ok := h.captchaConfig["driver"]; ok {
		capt, ok := captcha.Get(capDriver)

		if ok {
			if !capt.Validate(ctx.FormValue("token")) {
				response.BadRequest(ctx, "wrong captcha")
				return
			}
		}
	}

	if !exist {
		password := ctx.FormValue("password")
		username := ctx.FormValue("username")

		if password == "" || username == "" {
			response.BadRequest(ctx, "wrong password or username")
			return
		}
		user, ok = auth.Check(password, username, h.conn)
	} else {
		user, ok, errMsg = auth.GetService(s).P(ctx)
	}

	if !ok {
		response.BadRequest(ctx, errMsg)
		return
	}

	err := auth.SetCookie(ctx, user, h.conn)

	if err != nil {
		response.Error(ctx, err.Error())
		return
	}

	if ref := ctx.Referer(); ref != "" {
		if u, err := url.Parse(ref); err == nil {
			v := u.Query()
			if r := v.Get("ref"); r != "" {
				rr, _ := url.QueryUnescape(r)
				response.OkWithData(ctx, map[string]interface{}{
					"url": rr,
				})
				return
			}
		}
	}

	response.OkWithData(ctx, map[string]interface{}{
		"url": h.config.GetIndexURL(),
	})
}

// Logout delete the cookie.
func (h *Handler) Logout(ctx *context.Context) {
	err := auth.DelCookie(ctx, db.GetConnection(h.services))
	if err != nil {
		logger.ErrorCtx(ctx, "logout error %+v", err)
	}
	ctx.AddHeader("Location", h.config.Url(config.GetLoginUrl()))
	ctx.SetStatusCode(302)
}

// ShowLogin show the login page.
func (h *Handler) ShowLogin(ctx *context.Context) {

	tmpl, name := template.GetComp("login").GetTemplate()
	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, name, struct {
		UrlPrefix string
		Title     string
		Logo      template2.HTML
		CdnUrl    string
		System    types.SystemInfo
	}{
		UrlPrefix: h.config.AssertPrefix(),
		Title:     h.config.LoginTitle,
		Logo:      h.config.LoginLogo,
		System: types.SystemInfo{
			Version: system.Version(),
		},
		CdnUrl: h.config.AssetUrl,
	}); err == nil {
		ctx.HTML(http.StatusOK, buf.String())
	} else {
		logger.ErrorCtx(ctx, "ShowLogin error %+v", err)
		ctx.HTML(http.StatusOK, "parse template error (；′⌒`)")
	}
}
