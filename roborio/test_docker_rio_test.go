package roborio

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type dockerRio struct {
	ctx context.Context
	cli *client.Client
	t   *testing.T
	id  string
}

func newDockerRio(t *testing.T) *dockerRio {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Fatal(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/ocurr/ssh-server", types.ImagePullOptions{})
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "ocurr/ssh-server",
		ExposedPorts: nat.PortSet{
			"22/tcp": struct{}{},
		},
	}, hostConfig, nil, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		t.Fatal(err)
	}

	return &dockerRio{
		ctx,
		cli,
		t,
		resp.ID,
	}
}

func (d *dockerRio) shutdown() {
	err := d.cli.ContainerStop(d.ctx, d.id, nil)
	if err != nil {
		d.t.Fatal(err)
	}
}
