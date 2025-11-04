package internal

import (
	"context"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/bgrewell/dtac-agent/cmd/plugins/docker/plugin/internal/utilities"
)

// NewDockerClientWrapper returns a new wrapper around the docker client
func NewDockerClientWrapper() (wrapper *DockerClientWrapper, err error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return &DockerClientWrapper{
		client: client,
	}, nil
}

// DockerClientWrapper is a type for making working with the docker client a little more friendly
type DockerClientWrapper struct {
	client *docker.Client
}

// ListImages returns a list of the docker images on the system
func (w *DockerClientWrapper) ListImages(options ...ListImageOptions) (images []*ImageInfo, err error) {
	o := &listImageOptions{}
	for _, option := range options {
		option(o)
	}

	lio := docker.ListImagesOptions{
		Filters: o.filters,
		All:     o.all,
		Digests: o.digests,
		Filter:  o.filter,
	}

	images = make([]*ImageInfo, 0)
	imgs, err := w.client.ListImages(lio)
	if err != nil {
		return nil, err
	}
	for _, img := range imgs {
		ii := &ImageInfo{
			ID:          img.ID,
			RepoTags:    img.RepoTags,
			Created:     utilities.ConvertEpochTimeToTimestamp(img.Created),
			Size:        utilities.ConvertBytesToHumanReadable(img.Size),
			VirtualSize: utilities.ConvertBytesToHumanReadable(img.VirtualSize),
			ParentID:    img.ParentID,
			RepoDigests: img.RepoDigests,
			Labels:      img.Labels,
		}
		images = append(images, ii)
	}

	return images, nil
}

// ListConfigs returns a list of configurations if swarm is installed
func (w *DockerClientWrapper) ListConfigs() ([]swarm.Config, error) {
	return w.client.ListConfigs(docker.ListConfigsOptions{})
}

// ListContainers returns a list of docker containers on the system
func (w *DockerClientWrapper) ListContainers(options ...ListContainerOptions) ([]docker.APIContainers, error) {
	o := &listContainerOptions{}
	for _, option := range options {
		option(o)
	}

	lco := docker.ListContainersOptions{
		All:    o.all,
		Size:   o.size,
		Limit:  o.limit,
		Since:  o.since,
		Before: o.before,
	}

	return w.client.ListContainers(lco)
}

// ListNodes returns a list of nodes if docker swarm is enabled
func (w *DockerClientWrapper) ListNodes() ([]swarm.Node, error) {
	return w.client.ListNodes(docker.ListNodesOptions{
		Filters: nil,
		Context: nil,
	})
}

// ListNetworks returns a list of docker networks on the system
func (w *DockerClientWrapper) ListNetworks() ([]docker.Network, error) {
	return w.client.ListNetworks()
}

// ListPlugins returns a list of plugins on the system if swarm is installed
func (w *DockerClientWrapper) ListPlugins() ([]docker.PluginDetail, error) {
	return w.client.ListPlugins(context.Background())
}

// ListSecrets returns a list of secrets if swarm is installed
func (w *DockerClientWrapper) ListSecrets() ([]swarm.Secret, error) {
	return w.client.ListSecrets(docker.ListSecretsOptions{
		Filters: nil,
		Context: nil,
	})
}

// ListServices returns a list of docker services if swarm is installed
func (w *DockerClientWrapper) ListServices() ([]swarm.Service, error) {
	return w.client.ListServices(docker.ListServicesOptions{
		Filters: nil,
		Status:  false,
		Context: nil,
	})
}

// ListTasks returns a list of tasks if swarm is installed
func (w *DockerClientWrapper) ListTasks() ([]swarm.Task, error) {
	return w.client.ListTasks(docker.ListTasksOptions{
		Filters: nil,
		Context: nil,
	})
}

// ListVolumes returns a list of volumes on the system
func (w *DockerClientWrapper) ListVolumes() ([]docker.Volume, error) {
	return w.client.ListVolumes(docker.ListVolumesOptions{
		Filters: nil,
		Context: nil,
	})
}
