package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
)

func (r *Router) GetFavouriteUsers(ctx *fiber.Ctx) error {
	var res []chat.User
	for i := 0; i < 10; i++ {
		currentUser, err := r.userService.GetUser(ctx.UserContext(), i)
		if err != nil {
			continue
		} else {
			var u chat.User
			u.ID = currentUser.ID
			u.Username = currentUser.Username
			res = append(res, u)
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (r *Router) AddUserToFavourites(ctx *fiber.Ctx) error {
	return nil
}

func (r *Router) DeleteUserFromFavourites(ctx *fiber.Ctx) error {
	return nil
}
