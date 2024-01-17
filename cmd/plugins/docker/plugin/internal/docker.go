package internal

import (
	"context"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/intel-innersource/frameworks.automation.dtac.agent/cmd/plugins/docker/plugin/internal/utilities"
)

func NewDockerClientWrapper() (wrapper *DockerClientWrapper, err error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return &DockerClientWrapper{
		client: client,
	}, nil
}

type DockerClientWrapper struct {
	client *docker.Client
}

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

func (w *DockerClientWrapper) ListConfigs() ([]swarm.Config, error) {
	return w.client.ListConfigs(docker.ListConfigsOptions{})
}

func (w *DockerClientWrapper) ListContainers() ([]docker.APIContainers, error) {
	return w.client.ListContainers(docker.ListContainersOptions{
		All:     false,
		Size:    false,
		Limit:   0,
		Since:   "",
		Before:  "",
		Filters: nil,
	})
}

func (w *DockerClientWrapper) ListNodes() ([]swarm.Node, error) {
	return w.client.ListNodes(docker.ListNodesOptions{
		Filters: nil,
		Context: nil,
	})
}

func (w *DockerClientWrapper) ListNetworks() ([]docker.Network, error) {
	return w.client.ListNetworks()
}

func (w *DockerClientWrapper) ListPlugins() ([]docker.PluginDetail, error) {
	return w.client.ListPlugins(context.Background())
}

func (w *DockerClientWrapper) ListSecrets() ([]swarm.Secret, error) {
	return w.client.ListSecrets(docker.ListSecretsOptions{
		Filters: nil,
		Context: nil,
	})
}

func (w *DockerClientWrapper) ListServices() ([]swarm.Service, error) {
	return w.client.ListServices(docker.ListServicesOptions{
		Filters: nil,
		Status:  false,
		Context: nil,
	})
}

func (w *DockerClientWrapper) ListTasks() ([]swarm.Task, error) {
	return w.client.ListTasks(docker.ListTasksOptions{
		Filters: nil,
		Context: nil,
	})
}

func (w *DockerClientWrapper) ListVolumes() ([]docker.Volume, error) {
	return w.client.ListVolumes(docker.ListVolumesOptions{
		Filters: nil,
		Context: nil,
	})
}
