# CLAUDE.md

One tool per use case: decode args, call the use case, map the result. Register tools in `configurations/` under the `mcp_tools` group. Bearer auth lives in `middlewares/`; the host applies it.
