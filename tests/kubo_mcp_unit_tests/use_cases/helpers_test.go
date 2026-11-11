package use_cases_test

import "github.com/FournyP/kubo-mcp/src/kubo_mcp/settings"

const testGatewayUrl = "https://gateway.example"

func newTestIpfsSettings(maxUploadSizeBytes, maxReadSizeBytes int64) *settings.IpfsSettings {
	return &settings.IpfsSettings{
		MaxUploadSizeBytes: maxUploadSizeBytes,
		MaxReadSizeBytes:   maxReadSizeBytes,
		PublicGatewayUrl:   testGatewayUrl,
	}
}
