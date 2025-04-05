package views

import (
	"Stant/ECommerce/internal/views/templates"
	"context"
	"net/http"
)

func RenderLoginPage(w http.ResponseWriter, ctx context.Context) {
	templates.Login().Render(ctx, w)
}
