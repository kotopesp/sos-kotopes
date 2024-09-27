package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
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
		Chats []core.Chat `json:"chats"`
	}{
		Total: total,
		Chats: chats,
	}

	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) getChatWithUsersByID(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	chat, err := r.chatService.GetChatWithUsersByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(chat))
}

func (r *Router) createChat(ctx *fiber.Ctx) error {
	chatType := ctx.Query("type", "")
	var users struct {
		UserIds []int `json:"userIds"`
	}
	if err := ctx.BodyParser(&users); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	if len(users.UserIds) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("No users selected"))
	}

	existingChat, err := r.chatService.FindChatByUsers(ctx.UserContext(), users.UserIds)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	if existingChat.ID != -1 {
		return ctx.Status(fiber.StatusConflict).JSON(model.OKResponse(existingChat))
	}

	chat := core.Chat{
		ChatType: chatType,
	}

	createdChat, err := r.chatService.CreateChat(ctx.UserContext(), chat, users.UserIds)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(createdChat))
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
