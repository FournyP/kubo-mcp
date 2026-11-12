// Railway Infrastructure as Code: railway config plan | apply
//
// An apply deletes every resource this file does not declare, so link it to a
// project dedicated to kubo-mcp.
//
// Kubo is not declared here. Railway's private network does not cross projects,
// so a node in another project needs its public URL and auth.
// Export secrets for the first apply; later runs omit them and preserve() keeps Railway's.
//
//   export AUTH__API_KEYS=$(openssl rand -hex 32)
//   export KUBO__API_URL=https://<kubo-domain> KUBO__AUTHORIZATION="Bearer <kubo secret>"

import { defineRailway, github, preserve, project, service } from "railway/iac";

const REPO = "FournyP/kubo-mcp";

// Matched by name, so keep it identical to Railway: a mismatch is a delete and recreate.
const MCP_SERVICE = "kubo-mcp";

const PORT = "3000";

/** Push the value from the local environment if present, else keep Railway's. */
const fromEnvOrPreserve = (name: string) => process.env[name] ?? preserve();

export default defineRailway(() => {
  const mcp = service(MCP_SERVICE, {
    source: github(REPO, { branch: "main" }),
    build: {
      builder: "DOCKERFILE",
      dockerfilePath: "src/kubo_mcp_host/Dockerfile",
      watchPatterns: ["/go.mod", "/go.sum", "/src/**"],
    },
    deploy: {
      // /health skips auth, so it stays green without a key.
      healthcheckPath: "/health",
    },
    env: {
      // Railway routes the domain to PORT; HOST__PORT is what the binary binds.
      PORT,
      HOST__PORT: PORT,

      // Comma-separated bearer tokens. Never leave it empty on a public domain.
      AUTH__API_KEYS: fromEnvOrPreserve("AUTH__API_KEYS"),

      KUBO__API_URL: fromEnvOrPreserve("KUBO__API_URL"),
      KUBO__AUTHORIZATION: fromEnvOrPreserve("KUBO__AUTHORIZATION"),

      IPFS__PUBLIC_GATEWAY_URL: fromEnvOrPreserve("IPFS__PUBLIC_GATEWAY_URL"),
      MCP__ENABLED_TOOLS: fromEnvOrPreserve("MCP__ENABLED_TOOLS"),
      LOGGING__LEVEL: process.env.LOGGING__LEVEL ?? "info",
    },
  });

  return project("Kubo MCP", { resources: [mcp] });
});
