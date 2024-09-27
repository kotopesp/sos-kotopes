package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func (r *Router) getAllMembers(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	members, total, err := r.chatMemberService.GetAllMembers(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	response := struct {
		Total   int               `json:"total"`
		Members []core.ChatMember `json:"member"`
	}{
		Total:   total,
		Members: members,
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) addMemberToChat(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	userID, err := ctx.ParamsInt("user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	member := core.ChatMember{
		ChatID: chatID,
		UserID: userID,
	}
	// if err := ctx.BodyParser(&member); err != nil {
	// 	return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	// }
	// TODO: получать user_id
	createdMember, err := r.chatMemberService.AddMemberToChat(ctx.UserContext(), member)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(createdMember))
}

func (r *Router) updateMemberInfo(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	userID, err := ctx.ParamsInt("user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}
	updatedMember, err := r.chatMemberService.UpdateMemberInfo(ctx.UserContext(), chatID, userID)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(updatedMember))
}

func (r *Router) deleteMemberFromChat(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	userID, err := ctx.ParamsInt("user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}
	if err := r.chatMemberService.DeleteMemberFromChat(ctx.UserContext(), chatID, userID); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Member kicked from chat"))
}