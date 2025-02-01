package driver

import (
	"html/template"
	"io"
	"net/http"
	"strings"
	"github.com/labstack/echo/v4"
	"barbecue/core"
)

type HtmlDriver struct {
	Echo		echo.Context
}

func NewHtmlDriver(ctx echo.Context) core.Driver {
	return HtmlDriver { ctx }
}

const (
	T_INDEX = "index.html"
	T_TASK = "task.html"
	T_FORM = "form.html"
	T_ERROR = "error.html"
)

func (d HtmlDriver) redirect(uri string) error {
	if len(uri) == 0 || !strings.HasPrefix(uri, "/") { uri = "/" }
	return d.Echo.Redirect(http.StatusSeeOther, uri)
}

func (d HtmlDriver) Out(data interface{}) error {
	if data.(HtmlDriverData).redirect {
		return d.redirect(data.(HtmlDriverData).template)
	}
	return d.Echo.Render(data.(HtmlDriverData).code, data.(HtmlDriverData).template, data.(HtmlDriverData).data)
}

func (d HtmlDriver) Err(code int, msg string) error {
	return d.Echo.Render(code, T_ERROR, struct { Code int ; Message string } {
		Code: code,
		Message: http.StatusText(code),
	})
}

// Data

type HtmlDriverData struct {
	code		int
	redirect	bool
	template	string
	data		interface{}
}

func NewHtmlDriverData(code int, template string, data interface{}) HtmlDriverData {
	return HtmlDriverData { code, false, template, data }
}

func NewHtmlDriverRedirect(uri string) HtmlDriverData {
	return HtmlDriverData{ http.StatusSeeOther, true, uri, nil }
}

// Renderer

type HtmlDriverRenderer struct {
    Templates	*template.Template
}

func NewHtmlDriverRenderer() HtmlDriverRenderer {
	return HtmlDriverRenderer {
	    Templates: template.Must(template.ParseGlob("html/*.html")),
	}
}

func (r HtmlDriverRenderer) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	return r.Templates.ExecuteTemplate(w, name, data)
}
