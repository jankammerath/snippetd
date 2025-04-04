package snippetd

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

const containerdSocket = "/run/containerd/containerd.sock"
const snippetNamespace = "snippetd"

type CodeRuntime struct {
	client *containerd.Client
}

func NewCodeRuntime() (*CodeRuntime, error) {
	client, err := containerd.New(containerdSocket)
	if err != nil {
		return nil, err
	}

	runtime := &CodeRuntime{client: client}

	// create the namespace if it does not exist
	err = runtime.createNamespace()
	if err != nil {
		return nil, err
	}

	// remove all container leftovers
	err = runtime.removeContainers()
	if err != nil {
		return nil, err
	}

	return runtime, nil
}

// creates the namespace if it does not exist
func (runtime *CodeRuntime) createNamespace() error {
	ctx := context.Background()
	namespaces, err := runtime.client.NamespaceService().List(ctx)
	if err != nil {
		return err
	}

	for _, ns := range namespaces {
		if ns == snippetNamespace {
			return nil // Namespace already exists
		}
	}

	// Create the namespace if it doesn't exist
	err = runtime.client.NamespaceService().Create(ctx, snippetNamespace, map[string]string{})
	if err != nil {
		return err
	}

	return nil
}

// removes any containers all containers from the namespace
func (runtime *CodeRuntime) removeContainers() error {
	ctx := namespaces.WithNamespace(context.Background(), snippetNamespace)
	containers, err := runtime.client.Containers(ctx)
	if err != nil {
		return err
	}

	for _, container := range containers {
		err = container.Delete(ctx, containerd.WithSnapshotCleanup)
		if err != nil {
			return err
		}
	}

	return nil
}
