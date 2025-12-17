package api

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/vn-fin/xpb/xpb"
)

func GetUserID(c *fiber.Ctx) string {
	userID := c.Locals("userId")
	if userID == nil {
		return ""
	}
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

func GetUserEmail(c *fiber.Ctx) string {
	userInfo := c.Locals("userInfo")
	if userInfo == nil {
		return ""
	}
	if info, ok := userInfo.(*pb.UserInfo); ok {
		return info.Email
	}
	return ""
}
