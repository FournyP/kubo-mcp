package container_builders

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const kuboImage = "ipfs/kubo:v0.43.1"

// kuboContainer is a value, so parallel tests share it without a racy lazy init.
var kuboContainer = &KuboContainerBuilder{}

// KuboContainerBuilder shares one offline Kubo node; content addressing keeps parallel tests apart.
type KuboContainerBuilder struct {
	container   testcontainers.Container
	apiUrl      string
	containerMu sync.Mutex
}

func GetKuboContainerBuilder() *KuboContainerBuilder {
	return kuboContainer
}

// ApiUrl starts the node on first use and returns its RPC API base URL.
func (b *KuboContainerBuilder) ApiUrl() string {
	b.containerMu.Lock()
	defer b.containerMu.Unlock()

	if b.container == nil {
		ctx := context.Background()
		req := testcontainers.ContainerRequest{
			Image:        kuboImage,
			ExposedPorts: []string{"5001/tcp"},
			Cmd:          []string{"daemon", "--offline"},
			WaitingFor: wait.ForHTTP("/api/v0/version").
				WithMethod(http.MethodPost).
				WithPort("5001/tcp").
				WithStartupTimeout(2 * time.Minute),
		}

		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			log.Fatalf("failed to start kubo container: %v", err)
		}
		b.container = container

		host, _ := container.Host(ctx)
		port, _ := container.MappedPort(ctx, "5001")
		b.apiUrl = fmt.Sprintf("http://%s:%s", host, port.Port())

		log.Println("Kubo container ready")
	}
	return b.apiUrl
}

func (b *KuboContainerBuilder) stop() {
	b.containerMu.Lock()
	defer b.containerMu.Unlock()

	if b.container != nil {
		_ = b.container.Terminate(context.Background())
		b.container = nil
		log.Println("Kubo container stopped")
	}
}

// StopKuboContainer is called from each integration package's TestMain.
func StopKuboContainer() {
	GetKuboContainerBuilder().stop()
}
