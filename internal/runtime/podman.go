package runtime

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
)

type Podman struct {
	socketPath string
	conn       *context.Context
}

// NewPodmanRuntime creates a new Podman runtime manager
func NewPodmanRuntime() *Podman {
	// Get Podman socket location
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	return &Podman{
		socketPath: socket,
	}
}

func (p *Podman) Start() error {
	conn, err := bindings.NewConnection(context.Background(), p.socketPath)
	if err != nil {
		return err
	}

	p.conn = &conn

	containerLatestList, err := containers.List(conn, &containers.ListOptions{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Latest container is %s\n", containerLatestList)
	return nil
}
