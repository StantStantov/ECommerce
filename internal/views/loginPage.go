package views

import (
	templates "Stant/ECommerce/web"
	"context"
	"net/http"
)

func RenderLoginPage(w http.ResponseWriter, ctx context.Context) {
	templates.Login().Render(ctx, w)
}
