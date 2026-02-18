-- 0010_client_type_to_model
-- Move client_type from llm_providers to models, rename to new values, drop unsupported types.

-- 1) Add client_type column to models (nullable)
ALTER TABLE models ADD COLUMN IF NOT EXISTS client_type TEXT;

-- 2) Migrate data from provider to model with name mapping
UPDATE models SET client_type = CASE p.client_type
    WHEN 'openai' THEN 'openai-responses'
    WHEN 'openai-compat' THEN 'openai-completions'
    WHEN 'anthropic' THEN 'anthropic-messages'
    WHEN 'google' THEN 'google-generative-ai'
END
FROM llm_providers p
WHERE models.llm_provider_id = p.id
  AND p.client_type IN ('openai', 'openai-compat', 'anthropic', 'google');

-- 3) Delete models whose provider had an unsupported client_type (still NULL after step 2, AND type is chat)
DELETE FROM models WHERE client_type IS NULL AND type = 'chat';

-- 4) For embedding models, leave client_type as NULL (they don't need one)

-- 5) Delete providers with unsupported client_type
DELETE FROM llm_providers WHERE client_type NOT IN ('openai', 'openai-compat', 'anthropic', 'google');

-- 6) Add CHECK constraints: client_type values + chat models must have client_type
ALTER TABLE models ADD CONSTRAINT models_client_type_check
  CHECK (client_type IS NULL OR client_type IN ('openai-responses', 'openai-completions', 'anthropic-messages', 'google-generative-ai'));
ALTER TABLE models ADD CONSTRAINT models_chat_client_type_check
  CHECK (type != 'chat' OR client_type IS NOT NULL);

-- 7) Drop client_type from llm_providers
ALTER TABLE llm_providers DROP CONSTRAINT IF EXISTS llm_providers_client_type_check;
ALTER TABLE llm_providers DROP COLUMN IF EXISTS client_type;
