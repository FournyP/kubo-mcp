package services_test

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/errors"
	"github.com/FournyP/kubo-mcp/src/kubo_mcp_http_downloader/services"
)

type WhenCheckingADestinationTestingSuite struct {
	sut *services.DestinationGuard
}

func WhenCheckingADestinationBeforeEach() *WhenCheckingADestinationTestingSuite {
	return &WhenCheckingADestinationTestingSuite{sut: services.NewDestinationGuard()}
}

func TestWhenCheckingADestination(t *testing.T) {
	t.Parallel()

	t.Run("Given an address that cannot be parsed", func(t *testing.T) {
		t.Parallel()

		t.Run("Should return a blocked destination error", func(t *testing.T) {
			t.Parallel()

			suite := WhenCheckingADestinationBeforeEach()

			err := suite.sut.Control("tcp", "kubo:5001", nil)

			var blockedDestinationError *errors.BlockedDestinationError
			assert.ErrorAs(t, err, &blockedDestinationError)
		})
	})

	t.Run("Given a loopback address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be public", func(t *testing.T) {
			t.Parallel()

			assert.False(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("127.0.0.1")))
		})
	})

	t.Run("Given a private IPv4 address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be public", func(t *testing.T) {
			t.Parallel()

			assert.False(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("10.1.2.3")))
		})
	})

	t.Run("Given a cloud metadata address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be public", func(t *testing.T) {
			t.Parallel()

			assert.False(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("169.254.169.254")))
		})
	})

	t.Run("Given a carrier-grade NAT address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be public", func(t *testing.T) {
			t.Parallel()

			assert.False(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("100.64.0.1")))
		})
	})

	t.Run("Given a unique local IPv6 address", func(t *testing.T) {
		t.Parallel()

		// Railway's private network hands out fd12::/16 addresses.
		t.Run("Should not be public", func(t *testing.T) {
			t.Parallel()

			assert.False(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("fd12::1")))
		})
	})

	t.Run("Given an IPv4-mapped loopback address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be public", func(t *testing.T) {
			t.Parallel()

			assert.False(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("::ffff:127.0.0.1")))
		})
	})

	t.Run("Given a NAT64 address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should not be public", func(t *testing.T) {
			t.Parallel()

			assert.False(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("64:ff9b::a00:1")))
		})
	})

	t.Run("Given a public IPv4 address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be public", func(t *testing.T) {
			t.Parallel()

			assert.True(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("1.1.1.1")))
		})
	})

	t.Run("Given a public IPv6 address", func(t *testing.T) {
		t.Parallel()

		t.Run("Should be public", func(t *testing.T) {
			t.Parallel()

			assert.True(t, WhenCheckingADestinationBeforeEach().sut.IsPublic(netip.MustParseAddr("2606:4700::1111")))
		})
	})
}
