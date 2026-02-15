# Bot

A **bot** is the primary runtime entity in Memoh.

Each bot has its own:

- Configuration
- Container lifecycle
- Memory scope
- Channel bindings
- Model assignments

## Key Settings

- **max-load-time** (`max_context_load_time`): how many minutes of recent conversation context are loaded into prompts
- **language**: preferred language for interaction (default is `auto`)
- **chat model / memory model / embedding model**: model IDs used by this bot

## Why It Matters

The bot abstraction allows Memoh to isolate behavior and resources per agent, while keeping management centralized in one Web UI.

## Web UI Path

- `Bots > Select a bot > Settings`

