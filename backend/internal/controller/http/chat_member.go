package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
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
	responce := struct {
		Total   int               `json:"total"`
		Members []core.ChatMember `json:"member"`
	}{
		Total:   total,
		Members: members,
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responce))
}

func (r *Router) addMemberToChat(ctx *fiber.Ctx) error {
	chatId, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	var member core.ChatMember
	if err := ctx.BodyParser(&member); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	member.ChatID = chatId
	// TODO: получать user_id
	createdMember, err := r.chatMemberService.AddMemberToChat(ctx.UserContext(), member)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(createdMember))
}

func (r *Router) updateMemberInfo(ctx *fiber.Ctx) error {
	chatId, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	userId, err := ctx.ParamsInt("user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}
	updatedMember, err := r.chatMemberService.UpdateMemberInfo(ctx.UserContext(), chatId, userId)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(updatedMember))
}

func (r *Router) deleteMemberFromChat(ctx *fiber.Ctx) error {
	chatId, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	userId, err := ctx.ParamsInt("user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid user ID"))
	}
	if err := r.chatMemberService.DeleteMemberFromChat(ctx.UserContext(), chatId, userId); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Member kicked from chat"))
}
