package dispatcher

import (
	"github.com/docker/swarmkit/api"
	//"github.com/docker/swarmkit/ca"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"github.com/docker/swarmkit/manager/state/store"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// CreateExecutorAttachment allows the node to request the resources
// allocation needed for a container attachment on the specific node.
// - Returns `InvalidArgument` if the Spec is malformed.
// - Returns `NotFound` if the Network is not found.
// - Returns an error if the creation fails.
func (d *Dispatcher) CreateExecutorAttachment(ctx context.Context, request *api.CreateExecutorAttachmentRequest) (*api.CreateExecutorAttachmentResponse, error) {
	att := &api.ExecutorAttachment{
		ID:         identity.NewID(),
		NodeID:     request.Spec.NodeID,
		Spec:       *request.Spec,
		Attachment: &api.NetworkAttachment{},
	}

	d.store.View(func(tx store.ReadTx) {
		att.Attachment.Network = store.GetNetwork(tx, att.Spec.Target)
		if att.Attachment.Network == nil {
			if networks, err := store.FindNetworks(tx, store.ByName(att.Spec.Target)); err == nil && len(networks) > 0 {
				att.Attachment.Network = networks[0]
			}
		}
	})
	if att.Attachment.Network == nil {
		return nil, grpc.Errorf(codes.NotFound, "network %s not found", att.Spec.Target)
	}

	for _, addr := range att.Spec.Addresses {
		att.Attachment.Addresses = append(att.Attachment.Addresses, addr)
	}

	if err := d.store.Update(func(tx store.Tx) error {
		return store.CreateAttachment(tx, att)
	}); err != nil {
		return nil, err
	}
	log.G(ctx).Infof("\nAttachment %s created and saved to store\n", att.ID)
	return &api.CreateExecutorAttachmentResponse{ID: att.ID}, nil
}

// RemoveExecutorAttachment allows the node to request the release of
// the resources associated to the container attachment.
// - Returns `InvalidArgument` if attachment ID is not provided.
// - Returns `NotFound` if the attachment is not found.
// - Returns an error if the deletion fails.
func (d *Dispatcher) RemoveExecutorAttachment(ctx context.Context, request *api.RemoveExecutorAttachmentRequest) (*api.RemoveExecutorAttachmentResponse, error) {
	if request.ID == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, ErrInvalidArgument.Error())
	}

	if err := d.store.Update(func(tx store.Tx) error {
		return store.DeleteAttachment(tx, request.ID)
	}); err != nil {
		if err == store.ErrNotExist {
			return nil, grpc.Errorf(codes.NotFound, "attachment %s not found", request.ID)
		}
		return nil, err
	}
	log.G(ctx).Infof("\nAttachment %s removed from store\n", request.ID)
	return &api.RemoveExecutorAttachmentResponse{}, nil
}
