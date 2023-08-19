package syncer

import (
	"context"
	"reflect"

	v1 "agones.dev/agones/pkg/apis/agones/v1"
	"github.com/Octops/agones-event-broadcaster/pkg/events"
	"github.com/Octops/octops-image-syncer/pkg/runtime/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type ImageServiceClient interface {
	ImageStatus(ctx context.Context, request *pb.ImageStatusRequest) (*pb.ImageStatusResponse, error)
	PullImage(ctx context.Context, request *pb.PullImageRequest) (*pb.PullImageResponse, error)
}

// FleetImageSyncer implements the Broker interface used by the Agones Event Broadcaster to notify events
type FleetImageSyncer struct {
	imageClient ImageServiceClient
}

func NewFleetImageSyncer(client ImageServiceClient) *FleetImageSyncer {
	return &FleetImageSyncer{imageClient: client}
}

func (f *FleetImageSyncer) BuildEnvelope(event events.Event) (*events.Envelope, error) {
	envelope := &events.Envelope{}
	envelope.AddHeader("event_type", event.EventType().String())
	envelope.Message = event.(events.Message)

	return envelope, nil
}

func (f *FleetImageSyncer) SendMessage(envelope *events.Envelope) error {
	message := envelope.Message.(events.Message).Content()
	eventType := envelope.Header.Headers["event_type"]

	fleet, err := f.Unwrap(message)
	if err != nil {
		return errors.Wrap(err, "failed to process event")
	}

	switch eventType {
	case "fleet.events.added":
		fallthrough
	case "fleet.events.updated":
		return f.HandleAddedUpdated(fleet)
	case "fleet.events.deleted":
		//TODO: Consider a flag to decide if the image must be removed when a fleet is deleted
		//It may cause a race condition with running gameservers that are still in Terminating state
		log.Logger().Infof("fleet %s deleted", fleet.Name)
	}

	return nil
}

func (f *FleetImageSyncer) HandleAddedUpdated(fleet *v1.Fleet) error {
	image := fleet.Spec.Template.Spec.Template.Spec.Containers[0].Image
	fields := logrus.Fields{
		"fleet": fleet.GetName(),
		"image": image,
	}

	if ok, err := f.CheckImageStatus(image); err != nil {
		return errors.Wrap(err, "failed to check image status")
	} else if ok {
		log.Logger().WithFields(fields).Info("image already present")

		return nil
	}

	ref, err := f.PullImage(image)
	if err != nil {
		return errors.Wrap(err, "failed to pull image")
	}

	log.Logger().WithFields(fields).WithField("ref", ref).Info("fleet synced")

	return nil
}

func (f *FleetImageSyncer) Unwrap(message interface{}) (*v1.Fleet, error) {
	if fleet, ok := message.(*v1.Fleet); ok {
		return fleet, nil
	} else if fleet, ok := reflect.ValueOf(message).Field(1).Interface().(*v1.Fleet); ok {
		return fleet, nil
	}

	return nil, errors.New("message content is not a v1.Fleet")
}

func (f *FleetImageSyncer) CheckImageStatus(image string) (bool, error) {
	statusRequest := createImageStatusRequest(image)

	status, err := f.imageClient.ImageStatus(context.Background(), statusRequest)
	if err != nil {
		return false, errors.Wrap(err, "failed to get image status")
	}

	if status.Image != nil && len(status.Image.Id) > 0 {
		return true, nil
	}

	//Image is not present
	return false, nil
}

func (f *FleetImageSyncer) PullImage(image string) (string, error) {
	request := createPullImageRequest(image)

	resp, err := f.imageClient.PullImage(context.Background(), request)
	if err != nil {
		return "", errors.Wrap(err, "failed to pull image")
	}

	return resp.GetImageRef(), nil
}

func createPullImageRequest(image string) *pb.PullImageRequest {
	return &pb.PullImageRequest{
		Image: &pb.ImageSpec{
			Image: image,
		},
	}
}

func createImageStatusRequest(image string) *pb.ImageStatusRequest {
	return &pb.ImageStatusRequest{
		Image: &pb.ImageSpec{
			Image: image,
		},
	}
}
