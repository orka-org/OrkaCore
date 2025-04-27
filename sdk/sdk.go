package sdk

import (
	"context"

	pb "github.com/orka-org/orkacore/api/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents the Auth service client
type Client struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

// ClientOption is a function that configures the client
type ClientOption func(*Client)

func NewClient(ctx context.Context, target string, opts ...ClientOption) (*Client, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := &Client{
		conn:   conn,
		client: pb.NewAuthServiceClient(conn),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}
