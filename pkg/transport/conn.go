package transport

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func NewConn(target string) (*grpc.ClientConn, error) {
	if len(target) == 0 {
		return nil, errors.New("target is null, it should be a remote endpoint or a unix domain socket")
	}

	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial %s", target)
	}

	return conn, nil
}
