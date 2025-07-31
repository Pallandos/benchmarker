package containers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// ContainerInfo holds ID and Name of a container
type ContainerInfo struct {
	ID   string
	Name string
}

// ListContainerInfosByStack returns a slice of ContainerInfo (ID, Name) for a given stack name
func ListContainerInfos(stackName string) ([]ContainerInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(context.Background())

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	infos := make([]ContainerInfo, 0, len(containers))
	for _, c := range containers {
		if c.Labels["com.docker.compose.project"] == stackName || c.Labels["com.docker.stack.namespace"] == stackName {
			name := ""
			if len(c.Names) > 0 {
				name = c.Names[0]
			}
			infos = append(infos, ContainerInfo{ID: c.ID, Name: name})
		}
	}
	return infos, nil
}
