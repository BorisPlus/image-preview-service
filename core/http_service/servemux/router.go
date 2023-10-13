package servemux

import (
	"net/http"

	"github.com/BorisPlus/previewer/core/http_service/middleware"
	"github.com/BorisPlus/previewer/core/http_service/servemux/handlers/common"
	"github.com/BorisPlus/previewer/core/http_service/servemux/handlers/transformations"
	"github.com/BorisPlus/previewer/core/interfaces"

	rx "github.com/BorisPlus/regexphandlers"
)

type DefaultHandler struct{}

func (h DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/default", http.StatusTemporaryRedirect)
}

var (
	none      = rx.Params{}
	narrow_to = rx.Params{"height", "width", "url"}
)

func Handlers(logger interfaces.Logger, app interfaces.ImagePreviewProvider) rx.RegexpHandlers {
	return rx.NewRegexpHandlers(
		middleware.Instance().Listen(DefaultHandler{}),
		*rx.NewRegexpHandler(
			`/favicon.ico`,
			none,
			common.FaviconHandler{},
		),
		*rx.NewRegexpHandler(
			`/default`,
			none,
			middleware.Instance().Listen(common.DefaultHandler{Log: logger}),
		),
		*rx.NewRegexpHandler(
			`/{numeric}/{numeric}/{any}`,
			narrow_to,
			transformations.Narrow{Log: logger, ImageProvider: app},
		),
	)
}
