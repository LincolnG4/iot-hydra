package runtimer

import (
	"context"
	"time"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/rs/zerolog/log"
)

type PodmanRuntime interface {
	CreateContainer(container Container) error
	CheckContainer(name string) (Container, error)
	StartContainer(name string) error
	StopContainer(name string) error
	DeleteContainer(name string) error
	ListContainers() ([]Container, error)
}

type PodmanConnector struct {
	ctx        context.Context
	socketPath string
	timeout    time.Duration
}

func NewConnector(socketPath string) (PodmanConnector, error) {
	var err error

	conn := PodmanConnector{
		socketPath: socketPath,
	}
	conn.ctx, err = conn.Start()
	if err != nil {
		return PodmanConnector{}, err
	}

	return conn, nil
}

// Start the connector
func (p *PodmanConnector) Start() (context.Context, error) {
	return bindings.NewConnection(context.Background(), p.socketPath)
}

// Conn returns the podman connector context. If closed, it starts again
func (p *PodmanConnector) Conn() context.Context {
	if p.ctx == nil {
		conn, _ := p.Start()
		return conn
	}
	return p.ctx
}

type PodmanManager struct {
	Conn    *PodmanConnector
	Options *ManagerOptions
}

type ManagerOptions struct {
	SocketPath string
	Timeout    time.Duration
}

func (p *PodmanManager) Connection() context.Context {
	return p.Conn.Conn()
}

// NewPodmanManager creates a new Podman runtimer manager
func NewPodmanManager(opt *ManagerOptions) (PodmanManager, error) {
	conn, err := NewConnector(opt.SocketPath)
	if err != nil {
		return PodmanManager{}, err
	}
	// Get Podman socket location
	return PodmanManager{
		Conn:    &conn,
		Options: opt,
	}, nil
}

type Container struct {
	Name   string
	Image  string
	Config map[string]any
	State  string
}

type ContainerConfig struct {
	// Ports to be expose hostPort:cointanerPort
	ExposedPort map[int]int
}

// CreateContainer trigger all the steps to start a container:
// pull image -> Create the container -> Start
func (p PodmanManager) CreateContainer(container Container) error {
	_, err := images.Pull(p.Connection(), container.Image, &images.PullOptions{})
	if err != nil {
		return err
	}
	s := specgen.NewSpecGenerator(container.Image, false)
	s.Name = container.Name
	createResponse, err := containers.CreateWithSpec(p.Connection(), s, nil)
	if err != nil {
		return err
	}

	log.Info().Msg("Container created.")
	if err := containers.Start(p.Connection(), createResponse.ID, nil); err != nil {
		return err
	}
	return nil
}

// CheckContainer inspect the status of a container by their name
func (p PodmanManager) CheckContainer(name string) (Container, error) {
	inspectData, err := containers.Inspect(p.Connection(), name, new(containers.InspectOptions).WithSize(true))
	if err != nil {
		return Container{}, err
	}

	container := Container{
		Name:  name,
		Image: inspectData.ImageName,
		State: inspectData.State.Status,
	}

	return container, nil
}

// StartContainer by container name or id
func (p PodmanManager) StartContainer(name string) error {
	return containers.Start(p.Connection(), name, &containers.StartOptions{})
}

// StopContainer by container name or id
func (p PodmanManager) StopContainer(name string) error {
	return containers.Stop(p.Connection(), name, &containers.StopOptions{})
}

// DeleteContainer by container name or id
func (p PodmanManager) DeleteContainer(name string) error {
	_, err := containers.Remove(p.Connection(), name, &containers.RemoveOptions{})
	return err
}

// ListContainers return a list of all container and their status
func (p PodmanManager) ListContainers() ([]Container, error) {
	listContainers, err := containers.List(p.Connection(), &containers.ListOptions{})
	if err != nil {
		return nil, err
	}

	containers := make([]Container, 0)
	for _, con := range listContainers {
		container := Container{
			Name:  con.Names[0],
			Image: con.Image,
			State: con.State,
		}
		containers = append(containers, container)
	}

	return containers, err
}
