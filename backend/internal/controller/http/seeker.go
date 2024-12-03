package http

import (
	"github.com/gofiber/fiber/v2"
)

// @Summary		Get seeker
// @Tags			seeker
// @Description	Get seeker by id
// @ID				get-seeker
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	model.Response{data=seeker.ResponseSeekers}
// @Failure		400	{object}	model.Response
// @Failure		500	{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/seekers/{id}  [get]
func (r *Router) getSeeker(ctx *fiber.Ctx) error { return nil }

// @Summary		Create a seeker
// @Tags			seeker
// @Description	Create a seeker
// @ID				create-seeker
// @Accept			json
// @Produce		json
// @Param			request	body		seeker.ResponseSeekers	true	"Seeker"
// @Success		200		{object}	model.Response{data=seeker.ResponseSeekers}
// @Failure		400		{object}	model.Response
// @Failure		500		{object}	model.Response
// @Security		ApiKeyAuthBasic
// @Router			/seekers [post]
func (r *Router) createSeeker(ctx *fiber.Ctx) error { return nil }

func (r *Router) updateSeeker(ctx *fiber.Ctx) error { return nil }
