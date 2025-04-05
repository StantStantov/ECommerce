package views

import (
	"Stant/ECommerce/internal/domain"
	"Stant/ECommerce/internal/views/templates"
	"context"
	"fmt"
	"net/http"
	"reflect"
)

func RenderCategoryPage(category domain.Category, products []domain.Product, user domain.User,
	w http.ResponseWriter, responseCtx context.Context,
) {
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

	templates.Category(category, products, viewModel).Render(responseCtx, w)
}
