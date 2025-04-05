package views

import (
	"Stant/ECommerce/internal/domain/models"
	"Stant/ECommerce/internal/views/templates"
	"context"
	"fmt"
	"net/http"
	"reflect"
)

func RenderProductsPage(query string, products []models.Product, user models.User, w http.ResponseWriter, responseCtx context.Context) {
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

	templates.Products(query, products, viewModel).Render(responseCtx, w)
}
