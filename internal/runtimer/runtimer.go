package runtimer

type Runtimer struct {
	PodmanManager interface {
		CreateContainer(Container) error
		CheckContainer(string) (Container, error)
		StartContainer(string) error
		StopContainer(string) error
		DeleteContainer(string) error
		ListContainers() ([]Container, error)
	}
}

func NewRuntimer() Runtimer {
	return Runtimer{
		PodmanManager: PodmanManager{},
	}
}
