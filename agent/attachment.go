package agent

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/log"
	"golang.org/x/net/context"
)

// Listener allow to receive notifications about non Task related objects.
type Listener interface {
	// Notify notifies the listener about the updates over the executor attachment objects.
	Notify([]*api.ExecutorAttachment)
}

// ExecutorAttachmentManager provides control over executor attachments on this node.
type ExecutorAttachmentManager interface {
	// CreateAttachment allows the node to request the allocation of resources
	// needed for a container attachment on this node.
	CreateAttachment(ctx context.Context, request *api.CreateExecutorAttachmentRequest) (*api.CreateExecutorAttachmentResponse, error)

	// RemoveAttachment allows the node to request the release of
	// the resources associated to the container attachment.
	RemoveAttachment(context.Context, *api.RemoveExecutorAttachmentRequest) (*api.RemoveExecutorAttachmentResponse, error)

	// Register allows clients to register to executor attachment notifications.
	Register(Listener) (string, error)

	// Leave leaves the notification clients pool.
	Leave(string)
}

// ExecutorAttachmentManager returns the executor attachment management point.
func (n *Node) ExecutorAttachmentManager() ExecutorAttachmentManager {
	return n.agent
}

// Notifier provides notifications for non Task objects.
type Notifier interface {
	// Notify notifies the listeners about the updates over the executor attachment objects.
	Notify(context.Context, []*api.ExecutorAttachment)
}

type notifier struct {
	listeners map[string]Listener
	sync.RWMutex
}

func newNotifier() *notifier {
	return &notifier{
		listeners: make(map[string]Listener),
	}
}

func (n *notifier) Notify(ctx context.Context, eal []*api.ExecutorAttachment) {
	for _, l := range n.listeners {
		log.G(ctx).Debugf("Notifier: Notifying listener %v:", l)
		for _, a := range eal {
			log.G(ctx).Debugf("%s (%s) (%v): %v", a.ID, a.Spec.Annotations.Name, a.Status.State, a.Attachment.Addresses)
		}
		l.Notify(eal)
	}
}

// Register allows clients to register for notifications.
func (a *Agent) Register(l Listener) (string, error) {
	if l == nil {
		return "", fmt.Errorf("invalid listener")
	}

	id := identity.NewID()
	logrus.Debugf("Notifier: Registering listener %v (%s)", l, id)
	a.notifier.Lock()
	a.notifier.listeners[id] = l
	a.notifier.Unlock()
	return id, nil
}

// Leave let the client leave the notification pool.
func (a *Agent) Leave(listener string) {
	logrus.Debugf("Notifier: Leaving listener (%s)", listener)
	a.notifier.Lock()
	delete(a.notifier.listeners, listener)
	a.notifier.Unlock()
}

// CreateAttachment allows the node to request the allocation of
// resources needed for a container attachment on this node.
func (a *Agent) CreateAttachment(ctx context.Context, request *api.CreateExecutorAttachmentRequest) (*api.CreateExecutorAttachmentResponse, error) {
	client := api.NewDispatcherClient(a.config.Conn)
	return client.CreateExecutorAttachment(ctx, request)
}

// RemoveAttachment allows the node to request the release of
// the resources associated to the container attachment.
func (a *Agent) RemoveAttachment(ctx context.Context, request *api.RemoveExecutorAttachmentRequest) (*api.RemoveExecutorAttachmentResponse, error) {
	client := api.NewDispatcherClient(a.config.Conn)
	return client.RemoveExecutorAttachment(ctx, request)
}
