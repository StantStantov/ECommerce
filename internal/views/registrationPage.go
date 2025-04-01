package views

import (
	templates "Stant/ECommerce/web"
	"context"
	"net/http"
)

func RenderRegistrationPage(w http.ResponseWriter, ctx context.Context) {
	templates.Registration().Render(ctx, w)
}
