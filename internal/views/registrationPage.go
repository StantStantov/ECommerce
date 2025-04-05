package views

import (
	"Stant/ECommerce/internal/views/templates"
	"context"
	"net/http"
)

func RenderRegistrationPage(w http.ResponseWriter, ctx context.Context) {
	templates.Registration().Render(ctx, w)
}
