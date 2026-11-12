# CLAUDE.md

Entry point and DI root. Bearer auth wraps every route but `/health`; open when `auth.api_keys` is empty. Copy `appsettings.Example.json` or `.env.example` to run locally; never commit the filled copies.
