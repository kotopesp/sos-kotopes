package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model/user"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) GetFavouriteUsers(ctx *fiber.Ctx) error {
	idItem := getPayloadItem(ctx, "id")
	idFloat, ok := idItem.(float64)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("error while reading id from token"))
	}
	id := int(idFloat)
	var params user.GetFavourites
	if err := ctx.BodyParser(&params); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	FavouriteUsers, err := r.userFavouriteService.GetFavouriteUsers(ctx.UserContext())

}

func (r *Router) AddUserToFavourites(ctx fiber.Ctx) error {

}

func (r *Router) DeleteUserFromFavourites(ctx fiber.Ctx) error {

}
