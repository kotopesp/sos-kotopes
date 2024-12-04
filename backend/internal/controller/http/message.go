package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model"
	"github.com/kotopesp/sos-kotopes/internal/controller/http/model/chat"
	"github.com/kotopesp/sos-kotopes/pkg/logger"
	"strings"
)

func (r *Router) getAllMessages(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	sortType := ctx.Query("sort")
	if sortType != "asc" && sortType != "desc" && sortType != "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid sort type"))
	}
	searchText := ctx.Query("query")
	messages, total, err := r.messageService.GetAllMessages(ctx.UserContext(), chatID, userID, sortType, searchText)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	response := struct {
		Total    int            `json:"total"`
		Messages []chat.Message `json:"message"`
	}{
		Total:    total,
		Messages: messages,
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(response))
}

func (r *Router) markMessagesAsRead(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	err = r.messageService.MarkMessagesAsRead(ctx.UserContext(), chatID, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.JSON(model.OKResponse("Messages marked as read"))
}

func (r *Router) getUnreadMessageCount(ctx *fiber.Ctx) error {
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}

	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	count, err := r.messageService.GetUnreadMessageCount(ctx.UserContext(), chatID, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}

	return ctx.JSON(model.OKResponse(count))
}

func (r *Router) createMessage(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}

	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}

	var message chat.Message
	if err := ctx.BodyParser(&message); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	message.ChatID = chatID
	message.UserID = userID

	// If contains multipart/form-data, then it audio message
	if containsMultipartFormData(ctx) {
		file, err := ctx.FormFile("audio_bytes")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input: audio not found"))
		}

		openedFile, err := file.Open()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse("Failed to read audio bytes"))
		}

		byteBuffer := make([]byte, file.Size)
		_, err = openedFile.Read(byteBuffer)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
		}
		message.AudioBytes = byteBuffer
	}

	createdMessage, err := r.messageService.CreateMessage(ctx.UserContext(), message)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusOK).JSON(model.OKResponse(createdMessage))
}

func containsMultipartFormData(ctx *fiber.Ctx) bool {
	headers := ctx.GetReqHeaders()

	for _, contentType := range headers[fiber.HeaderContentType] {
		if strings.HasPrefix(contentType, fiber.MIMEMultipartForm) {
			return true
		}
	}

	return false
}

func (r *Router) updateMessage(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	messageID, err := ctx.ParamsInt("message_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid message ID"))
	}
	type UpdateMessageRequest struct {
		Content string `json:"message_content"`
	}
	var request UpdateMessageRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid input"))
	}
	updatedMessage, err := r.messageService.UpdateMessage(ctx.UserContext(), chatID, userID, messageID, request.Content)
	if err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse(updatedMessage))
}

func (r *Router) deleteMessage(ctx *fiber.Ctx) error {
	userID, err := getIDFromToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse(err.Error()))
	}
	chatID, err := ctx.ParamsInt("chat_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid chat ID"))
	}
	messageID, err := ctx.ParamsInt("message_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse("Invalid message ID"))
	}
	if err := r.messageService.DeleteMessage(ctx.UserContext(), chatID, userID, messageID); err != nil {
		logger.Log().Error(ctx.UserContext(), err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse(err.Error()))
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.OKResponse("Message deleted"))
}
