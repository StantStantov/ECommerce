package views

import (
	"Stant/ECommerce/internal/domain"
	templates "Stant/ECommerce/web"
	"context"
	"fmt"
	"net/http"
	"reflect"
)

func RenderIndexPage(categories []domain.Category, user domain.User, w http.ResponseWriter, ctx context.Context) {
	var viewModel templates.UserViewModel
	if !reflect.ValueOf(user).IsZero() {
		viewModel = templates.UserViewModel{
			IsLogged: true,
			Name:     fmt.Sprintf("%s %s", user.FirstName(), user.SecondName()),
		}
	} else {
		viewModel = templates.UserViewModel{
			IsLogged: false,
			Name:     "Anonymous",
		}
	}

	templates.Index(categories, viewModel).Render(ctx, w)
}
