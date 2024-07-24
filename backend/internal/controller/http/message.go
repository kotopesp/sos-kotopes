package http

import (
	"gitflic.ru/spbu-se/sos-kotopes/internal/controller/http/model"
	"gitflic.ru/spbu-se/sos-kotopes/internal/core"
	"gitflic.ru/spbu-se/sos-kotopes/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) getAllMessages(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	sortType := ctx.Query("sort")
	if sortType != "asc" && sortType != "desc" && sortType != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid sort type"))
	}
	searchText := ctx.Query("query")
	messages, total, err := r.messageService.GetAllMessages(ctx.UserContext(), id, sortType, searchText)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	responce := struct {
		Total    int            `json:"total"`
		Messages []core.Message `json:"message"`
	}{
		Total:    total,
		Messages: messages,
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(responce))
}

func (r *Router) createMessage(ctx *fiber.Ctx) error {
	chatId, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	var message core.Message
	if err := ctx.BodyParser(&message); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	message.ChatID = chatId
	createdMessage, err := r.messageService.CreateMessage(ctx.UserContext(), message)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(createdMessage))
}

func (r *Router) updateMessage(ctx *fiber.Ctx) error {
	chatId, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	messageId, err := ctx.ParamsInt("message_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid message ID"))
	}
	type UpdateMessageRequest struct {
		Content string `json:"content"`
	}
	var request UpdateMessageRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	updatedMessage, err := r.messageService.UpdateMessage(ctx.UserContext(), chatId, messageId, request.Content)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(updatedMessage))
}

func (r *Router) deleteMessage(ctx *fiber.Ctx) error {
	chatId, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	messageId, err := ctx.ParamsInt("message_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid message ID"))
	}
	if err := r.messageService.DeleteMessage(ctx.UserContext(), chatId, messageId); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Message deleted"))
}
