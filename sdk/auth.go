package sdk

import (
	"context"

	pb "github.com/orka-org/orkacore/api/auth/v1"
)

// Register creates a new user account
func (c *Client) Register(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", err
	}
	return resp.AccessToken, resp.RefreshToken, nil
}

// Login authenticates a user
func (c *Client) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", "", err
	}
	return resp.AccessToken, resp.RefreshToken, nil
}

// ValidateToken validates a token and returns the associated user
func (c *Client) ValidateToken(ctx context.Context, token string) (bool, *User, error) {
	resp, err := c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, nil, err
	}
	if !resp.Valid {
		return false, nil, nil
	}
	return true, protoToUser(resp.User), nil
}

// GetUser retrieves user information by ID
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	return protoToUser(resp.User), nil
}
