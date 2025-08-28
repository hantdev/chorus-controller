package worker

import (
	"context"

	pb "github.com/clyso/chorus/proto/gen/go/chorus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	addr string
}

func New(addr string) *Client { return &Client{addr: addr} }

func (c *Client) Dial(ctx context.Context) (pb.ChorusClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(c.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return pb.NewChorusClient(conn), conn, nil
}
