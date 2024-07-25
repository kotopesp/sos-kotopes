package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func IsValidType(sortType string) bool {
	switch sortType {
	case
		"seeker",
		"keeper",
		"":
		return true
	}
	return false
}

func (r *Router) getAllChats(ctx *fiber.Ctx) error {
	sortType := ctx.Query("chat_type")
	if !IsValidType(sortType) {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat type"))
	}

	chats, total, err := r.chatService.GetAllChats(ctx.UserContext(), sortType)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	response := struct {
		Total int         `json:"total"`
		Chats []core.Chat `json:"chat"`
	}{
		Total: total,
		Chats: chats,
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) getChatByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	chat, err := r.chatService.GetChatByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(chat))
}

func (r *Router) createChat(ctx *fiber.Ctx) error {
	var chat core.Chat
	if err := ctx.BodyParser(&chat); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	createdChat, err := r.chatService.CreateChat(ctx.UserContext(), chat)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(createdChat))
}

func (r *Router) deleteChat(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	if err := r.chatService.DeleteChat(ctx.UserContext(), id); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Chat deleted"))
}
