package dockerinitiator

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	dockerinitiator "github.com/Storytel/go-docker-initiator"
)

// PubSubInstance contains the instance config for a PubSub image
type PubSubInstance struct {
	*dockerinitiator.Instance
	project string
	PubSubConfig
}

var (
	DefaultImage = "storytel/google-cloud-pubsub-emulator"

	DefaultCmd = []string{"--host=0.0.0.0", "--port=8262"}

	DefaultExposedPort = "8262"
)

// PubSubConfig contains configs for pubsub
type PubSubConfig struct {
	// ProbeTimeout specifies the timeout for the probing.
	// A timeout results in a startup error, if left empty a default value is used
	ProbeTimeout time.Duration

	// Image specifies the image used for the Mysql docker instance.
	// If left empty it will be set to DefaultImage
	Image string

	// Cmd is the commands that will run in the container
	// Is left empty it will be set to DefaultCmd
	Cmd []string

	// ExposedPort sets the exposed port of the container
	// If left empty it will be set to DefaultExposedPort
	ExposedPort string
}

// PubSub will create a PubSub instance container
func PubSub(config PubSubConfig) (*PubSubInstance, error) {

	if config.ProbeTimeout == 0 {
		config.ProbeTimeout = 10 * time.Second
	}

	if config.Image == "" {
		config.Image = DefaultImage
	}

	if config.ExposedPort == "" {
		config.ExposedPort = DefaultExposedPort
	}

	if len(config.Cmd) == 0 {
		config.Cmd = DefaultCmd
	}

	i, err := dockerinitiator.CreateContainer(
		dockerinitiator.ContainerConfig{
			Image:         config.Image,
			Cmd:           config.Cmd,
			ContainerPort: config.ExposedPort,
		},
		dockerinitiator.HTTPProbe{})
	if err != nil {
		return nil, err
	}

	project := "__docker_initiator__project-" + strconv.Itoa(rand.Int())[:8]
	psi := &PubSubInstance{
		i,
		project,
		config,
	}

	if err = psi.Probe(psi.ProbeTimeout); err != nil {
		return nil, err
	}

	return psi, nil
}

// Setenv sets the required variables for running against the emulator
func (psi *PubSubInstance) Setenv() error {
	err := os.Setenv("PUBSUB_EMULATOR_HOST", psi.GetHost())
	if err != nil {
		return err
	}

	err = os.Setenv("GOOGLE_CLOUD_PROJECT", psi.GetProject())
	if err != nil {
		return err
	}

	return nil
}

// GetProject fetches the project for the pubsub instance
func (psi *PubSubInstance) GetProject() string {
	return psi.project
}
