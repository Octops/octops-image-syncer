package clients

import (
	"google.golang.org/grpc"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

type ImageServiceClient struct {
	conn *grpc.ClientConn
	pb.ImageServiceClient
}

func NewImageServiceClient(conn *grpc.ClientConn) *ImageServiceClient {
	return &ImageServiceClient{
		conn, pb.NewImageServiceClient(conn),
	}
}
