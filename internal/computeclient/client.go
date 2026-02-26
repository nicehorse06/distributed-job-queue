package computeclient

import (
	"context"

	computev1 "distributed-job-queue/internal/gen/compute/v1"
)

// Client wraps the gRPC compute client with domain-specific methods.
type Client struct {
	raw computev1.ComputeServiceClient
}

// New creates a compute client.
func New(raw computev1.ComputeServiceClient) *Client {
	return &Client{raw: raw}
}

// Square calls the Rust compute service for squaring a number.
func (c *Client) Square(ctx context.Context, value int64) (int64, error) {
	resp, err := c.raw.Square(ctx, &computev1.SquareRequest{Value: value})
	if err != nil {
		return 0, err
	}

	return resp.GetSquare(), nil
}
