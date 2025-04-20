package helper

import (
	"app_monitor/pkg/global/model"

	utils_v1 "github.com/FDSAP-Git-Org/hephaestus/utils/v1"
	"github.com/gofiber/fiber/v3"
)

func JSONResponse(c fiber.Ctx, retCode, retMessage string, httpStatusCode int) error {
	return c.Status(httpStatusCode).JSON(model.Response{
		ResponseTime: utils_v1.GetResponseTime(c),
		Device:       string(c.RequestCtx().UserAgent()),
		RetCode:      retCode,
		Message:      retMessage,
	})
}

func JSONResponseWithData(c fiber.Ctx, retCode, retMessage string, data any, httpStatusCode int) error {
	return c.Status(httpStatusCode).JSON(model.Response{
		ResponseTime: utils_v1.GetResponseTime(c),
		Device:       string(c.RequestCtx().UserAgent()),
		RetCode:      retCode,
		Message:      retMessage,
		Data:         data,
	})
}

func JSONResponseWithError(c fiber.Ctx, retCode, retMessage string, err error, httpStatusCode int) error {
	return c.Status(httpStatusCode).JSON(model.Response{
		ResponseTime: utils_v1.GetResponseTime(c),
		Device:       string(c.RequestCtx().UserAgent()),
		RetCode:      retCode,
		Message:      retMessage,
		Error:        err.Error(),
	})
}
