package sdk

import (
	"context"
	"fmt"

	userv1 "github.com/user/go-templates/template-grpc-sdk/gen/go/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	User userv1.UserServiceClient
	conn *grpc.ClientConn
}

func NewClient(target string) (*Client, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return &Client{
		User: userv1.NewUserServiceClient(conn),
		conn: conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// Example method
func (c *Client) GetUser(ctx context.Context, id string) (*userv1.User, error) {
	resp, err := c.User.GetUser(ctx, &userv1.GetUserRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.GetUser(), nil
}
