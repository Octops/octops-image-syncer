package clients

import (
	"context"
	"google.golang.org/grpc"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type ImageServiceClient struct {
	conn   *grpc.ClientConn
	client pb.ImageServiceClient
}

func (c *ImageServiceClient) ImageStatus(ctx context.Context, request *pb.ImageStatusRequest) (*pb.ImageStatusResponse, error) {
	return c.client.ImageStatus(ctx, request)
}

func (c *ImageServiceClient) PullImage(ctx context.Context, request *pb.PullImageRequest) (*pb.PullImageResponse, error) {
	return c.client.PullImage(ctx, request)
}

func NewImageServiceClient(conn *grpc.ClientConn) *ImageServiceClient {
	return &ImageServiceClient{
		conn, pb.NewImageServiceClient(conn),
	}
}
