package services_test

import (
	"os"
	"testing"

	"github.com/FournyP/kubo-mcp/tests/common_tests/container_builders"
)

func TestMain(m *testing.M) {
	code := m.Run()
	container_builders.StopKuboContainer()
	os.Exit(code)
}
