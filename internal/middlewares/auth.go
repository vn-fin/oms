package middlewares

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vn-fin/oms/internal/api"
	"github.com/vn-fin/oms/internal/remote"
	pb "github.com/vn-fin/xpb/xpb"
)

// AuthMiddleware validates the JWT token via gRPC auth service
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authorization := c.Get("Authorization", "")
		if authorization == "" {
			return api.Response().Unauthorized("authorization header is missing").Send(c)
		}

		// Remove "Bearer " prefix if present
		token := strings.TrimPrefix(authorization, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			return api.Response().Unauthorized("token is empty").Send(c)
		}

		// Call gRPC auth service to verify token
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &pb.CheckAuthRequest{
			Token: token,
		}

		resp, err := remote.AuthGrpcClient.CheckAuthToken(ctx, req)
		if err != nil {
			return api.Response().InternalError(err).Send(c)
		}

		// Check if authentication failed
		if resp.Message != "" {
			return api.Response().Unauthorized(resp.Message).Send(c)
		}

		// Check if user info is present
		if resp.UserInfo == nil {
			return api.Response().Unauthorized("user info is missing in auth response").Send(c)
		}

		// Store user info in context for use by handlers
		c.Locals("userInfo", resp.UserInfo)
		c.Locals("userId", resp.UserInfo.UserId)

		return c.Next()
	}
}
