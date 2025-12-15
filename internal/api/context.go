package api

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/vn-fin/xpb/xpb"
)

// GetUserID extracts the user ID from the fiber context
// This should be used in handlers that are protected by AuthMiddleware
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

// GetUserInfo extracts the full user info from the fiber context
// This should be used in handlers that are protected by AuthMiddleware
func GetUserInfo(c *fiber.Ctx) *pb.UserInfo {
	userInfo := c.Locals("userInfo")
	if userInfo == nil {
		return nil
	}
	if info, ok := userInfo.(*pb.UserInfo); ok {
		return info
	}
	return nil
}
