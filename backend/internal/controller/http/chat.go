package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
)

func IsValidType(sortType string) bool {
	switch sortType {
	case
		"seeker",
		"keeper",
		"vet",
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

	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	chats, total, err := r.chatService.GetAllChats(ctx.UserContext(), sortType, userID)
	if err != nil {
		logger.Log().Debug(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse(err.Error()))
	}

	response := struct {
		Total int         `json:"total"`
		Chats []chat.Chat `json:"chats"`
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
	data, err := r.chatService.GetChatWithUsersByID(ctx.UserContext(), id)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(data))
}

func (r *Router) createChat(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
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

	isChosen := false
	for _, user := range users.UserIds {
		if user == userID {
			isChosen = true
		}
	}
	if !isChosen {
		users.UserIds = append(users.UserIds, userID)
	}

	existingChat, err := r.chatService.FindChatByUsers(ctx.UserContext(), users.UserIds)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	if existingChat.ID != -1 {
		return ctx.Status(fiber.StatusConflict).JSON(model.OKResponse(existingChat))
	}

	data := chat.Chat{
		ChatType: chatType,
	}

	createdChat, err := r.chatService.CreateChat(ctx.UserContext(), data, users.UserIds)
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
