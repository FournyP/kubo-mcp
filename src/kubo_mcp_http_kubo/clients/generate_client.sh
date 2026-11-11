#!/usr/bin/env sh
set -eu

# Regenerates the Kubo RPC client from the hand-maintained kubo_rpc.v3.yaml subset.

CLIENTS_DIR="$(CDPATH= cd -- "$(dirname "$0")" && pwd)"
SPEC_PATH="${CLIENTS_DIR}/kubo_rpc.v3.yaml"
OUTPUT_FILE="${CLIENTS_DIR}/oapi_client.gen.go"

OAPI_CODEGEN_VERSION="v2.5.1"

echo "Generating Kubo RPC client with oapi-codegen ${OAPI_CODEGEN_VERSION}"
(
	cd "${CLIENTS_DIR}"
	go run "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@${OAPI_CODEGEN_VERSION}" \
		-generate types,client \
		-package clients \
		-o "${OUTPUT_FILE}" \
		"${SPEC_PATH}"
)

echo "Done"
