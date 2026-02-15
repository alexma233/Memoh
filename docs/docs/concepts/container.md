# Container

Each bot runs in its own **isolated container**.

## What Isolation Gives You

- Separate filesystem per bot
- Separate runtime process space
- Controlled lifecycle (create/start/stop/delete)

This prevents one bot from interfering with another bot's execution environment.

## Why It Matters

Container isolation is the foundation that allows bots to run tools, commands, and file operations safely in parallel.

## Web UI Path

- `Bots > Select a bot > Container`

