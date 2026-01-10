# Agent CLI

A command-line interface for the personal housekeeper assistant agent.

## Setup

1. Create a `.env` file in the project root (not in this directory) with the following variables:

```env
# Main chat model
MODEL=gpt-4o
BASE_URL=https://api.openai.com/v1
API_KEY=your-api-key-here
MODEL_CLIENT_TYPE=openai

# Embedding model for memory search (using mem0ai)
EMBEDDING_MODEL=text-embedding-3-small
EMBEDDING_BASE_URL=https://api.openai.com/v1
EMBEDDING_API_KEY=your-api-key-here
EMBEDDING_CLIENT_TYPE=openai
EMBEDDING_DIMENSIONS=1536

# Summary model for memory generation (optional, defaults to main model)
SUMMARY_MODEL=gpt-4o-mini
SUMMARY_BASE_URL=https://api.openai.com/v1
SUMMARY_API_KEY=your-api-key-here
SUMMARY_CLIENT_TYPE=openai

# Database
DATABASE_URL=postgresql://user:password@localhost:5432/byte
```

2. Make sure the database is set up and running (required for memory storage).

## Usage

Run the CLI from the agent package:

```bash
pnpm start
```

Or with Bun directly:

```bash
bun run index.ts
```

## Features

- **Interactive Chat**: Type your messages and get responses from the AI agent
- **Long-term Memory**: Conversations are automatically saved with LLM-generated summaries
- **Context Loading**: Automatically loads recent conversations (last 60 minutes)
- **Memory Search**: The agent can search through past conversations using natural language and embeddings
- **Tool Calling**: Supports automatic tool execution with multi-step reasoning
- **Multi-Provider Support**: Works with OpenAI, Anthropic, and Google AI (via Vercel AI SDK)

## Commands

- Type your message and press Enter to chat
- Type `exit` or `quit` to close the application

## Environment Variables

### Required

- `MODEL`: The main LLM model ID (e.g., `gpt-4o`, `claude-3-5-sonnet-20241022`)
- `BASE_URL`: The API base URL for the main model
- `API_KEY`: Your API key for the main model
- `EMBEDDING_MODEL`: The embedding model for memory search (e.g., `text-embedding-3-small`)
- `DATABASE_URL`: PostgreSQL connection string with pgvector extension enabled

### Optional

- `MODEL_CLIENT_TYPE`: The model provider type (default: `openai`, options: `openai`, `anthropic`, `google`)
- `EMBEDDING_BASE_URL`: Base URL for embedding API (default: same as `BASE_URL`)
- `EMBEDDING_API_KEY`: API key for embedding (default: same as `API_KEY`)
- `EMBEDDING_CLIENT_TYPE`: Provider type for embedding (default: `openai`)
- `EMBEDDING_DIMENSIONS`: The dimensions of the embedding model (default: `1536`)
- `SUMMARY_MODEL`: The model used to summarize conversations for memory (default: same as `MODEL`)
- `SUMMARY_BASE_URL`: Base URL for summary model (default: same as `BASE_URL`)
- `SUMMARY_API_KEY`: API key for summary model (default: same as `API_KEY`)
- `SUMMARY_CLIENT_TYPE`: Provider type for summary model (default: same as `MODEL_CLIENT_TYPE`)

## Memory System

The agent uses [mem0ai](https://github.com/mem0ai/mem0) for sophisticated memory management:

1. **Conversation Storage**: After each conversation, mem0ai uses an LLM to extract and store key information
2. **Embedding Generation**: The extracted information is converted to embedding vectors using your configured embedding model
3. **Vector Storage**: Embeddings are stored in PostgreSQL with pgvector extension
4. **Semantic Search**: The agent can search past memories using natural language queries via vector similarity
5. **Context Loading**: Recent conversations are automatically loaded from history into context

### Model Recommendations

- **Main Model**: Use a powerful model like `gpt-4o` or `claude-3-5-sonnet-20241022` for best conversation quality
- **Summary Model**: Use a cheaper/faster model like `gpt-4o-mini` for memory extraction to save costs
- **Embedding Model**: Use `text-embedding-3-small` (1536 dims) or `text-embedding-3-large` (3072 dims) for OpenAI

### Database Requirements

The memory system requires PostgreSQL with the `pgvector` extension installed. Make sure:
1. PostgreSQL is installed and running
2. pgvector extension is enabled: `CREATE EXTENSION IF NOT EXISTS vector;`
3. mem0ai will automatically create the required tables on first run

