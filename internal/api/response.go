package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/vn-fin/oms/internal/models"
)

type ResponseBuilder struct {
	success bool
	status  int
	message string
	data    interface{}
	page    *models.Pagination
}

// Create a new builder with default values
func Response() *ResponseBuilder {
	return &ResponseBuilder{
		success: true,
		status:  fiber.StatusOK,
		message: "",
		data:    nil,
		page:    nil,
	}
}

// --- setters ---

func (r *ResponseBuilder) Success(v bool) *ResponseBuilder {
	r.success = v
	return r
}

func (r *ResponseBuilder) Status(code int) *ResponseBuilder {
	r.status = code
	return r
}

func (r *ResponseBuilder) Message(msg string) *ResponseBuilder {
	r.message = msg
	return r
}

func (r *ResponseBuilder) Data(d interface{}) *ResponseBuilder {
	r.data = d
	return r
}

func (r *ResponseBuilder) Page(p *models.Pagination) *ResponseBuilder {
	r.page = p
	return r
}

// --- helper methods like BadRequest, Unauthorized, etc. ---

func (r *ResponseBuilder) BadRequest(msg string) *ResponseBuilder {
	log.Error().Msg(msg)
	r.success = false
	r.status = fiber.StatusBadRequest
	r.message = msg
	return r
}

func (r *ResponseBuilder) Unauthorized(msg string) *ResponseBuilder {
	log.Error().Msg(msg)
	r.success = false
	r.status = fiber.StatusUnauthorized
	r.message = msg
	return r
}

func (r *ResponseBuilder) Forbidden(msg string) *ResponseBuilder {
	log.Error().Msg(msg)
	r.success = false
	r.status = fiber.StatusForbidden
	r.message = msg
	return r
}

func (r *ResponseBuilder) NotFound(msg string) *ResponseBuilder {
	log.Error().Msg(msg)
	r.success = false
	r.status = fiber.StatusNotFound
	r.message = msg
	return r
}

func (r *ResponseBuilder) InternalError(err error) *ResponseBuilder {
	log.Error().Err(err).Msg("internal server error")
	r.success = false
	r.status = fiber.StatusInternalServerError
	if err != nil {
		r.message = err.Error()
	} else {
		r.message = "internal server error"
	}
	return r
}

// Final send
func (r *ResponseBuilder) Send(c *fiber.Ctx) error {
	return c.Status(r.status).JSON(models.DefaultResponseModel{
		Success:    r.success,
		StatusCode: r.status,
		Message:    r.message,
		Data:       r.data,
		Page:       r.page,
	})
}
