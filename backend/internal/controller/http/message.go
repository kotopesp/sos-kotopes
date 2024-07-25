package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/core"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
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
	response := struct {
		Total    int            `json:"total"`
		Messages []core.Message `json:"message"`
	}{
		Total:    total,
		Messages: messages,
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) createMessage(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	var message core.Message
	if err := ctx.BodyParser(&message); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	message.ChatID = chatID
	createdMessage, err := r.messageService.CreateMessage(ctx.UserContext(), message)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(createdMessage))
}

func (r *Router) updateMessage(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	messageID, err := ctx.ParamsInt("message_id")
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
	updatedMessage, err := r.messageService.UpdateMessage(ctx.UserContext(), chatID, messageID, request.Content)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(updatedMessage))
}

func (r *Router) deleteMessage(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	messageID, err := ctx.ParamsInt("message_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid message ID"))
	}
	if err := r.messageService.DeleteMessage(ctx.UserContext(), chatID, messageID); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Message deleted"))
}
