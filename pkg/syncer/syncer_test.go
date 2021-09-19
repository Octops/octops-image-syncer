package syncer

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"testing"
)

func TestFleetImageSyncer_CheckImageStatus(t *testing.T) {
	testCases := []struct {
		name     string
		image    string
		response *pb.ImageStatusResponse
		want     bool
		wantErr  bool
		err      error
	}{
		{
			name:     "image is not present",
			image:    "gameserver:latest",
			response: &pb.ImageStatusResponse{},
			want:     false,
			wantErr:  false,
		},
		{
			name:  "image is present",
			image: "gameserver:latest",
			response: &pb.ImageStatusResponse{
				Image: &pb.Image{
					Id: "sha256:f8cdc89145cb0b5d6ee2ea95968310c45e4f453dd24ac682ff13f50f0d4b921d",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name:     "error checking status",
			image:    "gameserver:latest",
			response: &pb.ImageStatusResponse{},
			want:     false,
			wantErr:  true,
			err:      errors.New("failed to get image status"),
		},
	}

	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := new(imageServiceClient)
			imageSyncer := NewFleetImageSyncer(client)
			request := createImageStatusRequest(tc.image)
			client.On("ImageStatus", ctx, request).Return(tc.response, tc.err)

			ok, err := imageSyncer.CheckImageStatus(tc.image)
			require.Equal(t, tc.want, ok)
			require.Equal(t, err != nil, tc.wantErr)
			require.ErrorIs(t, err, tc.err)
		})
	}
}

type imageServiceClient struct {
	mock.Mock
}

func (m *imageServiceClient) ImageStatus(ctx context.Context, request *pb.ImageStatusRequest) (*pb.ImageStatusResponse, error) {
	args := m.Called(ctx, request)

	return args.Get(0).(*pb.ImageStatusResponse), args.Error(1)
}

func (m *imageServiceClient) PullImage(ctx context.Context, request *pb.PullImageRequest) (*pb.PullImageResponse, error) {
	args := m.Called(ctx, request)

	return args.Get(0).(*pb.PullImageResponse), args.Error(1)
}
