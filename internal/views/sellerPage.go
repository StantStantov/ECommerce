package views

import (
	"Stant/ECommerce/internal/domain/models"
	"Stant/ECommerce/internal/views/templates"
	"context"
	"fmt"
	"net/http"
	"reflect"
)

func RenderSellerPage(seller models.Seller, products []models.Product, user models.User,
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

	templates.Seller(seller, products, viewModel).Render(responseCtx, w)
}
