CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE users (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email                VARCHAR(255) UNIQUE,
    encrypted_password   TEXT,
    email_confirmed_at   TIMESTAMPTZ,
    invited_at           TIMESTAMPTZ,
    confirmation_token   VARCHAR(255) DEFAULT '',
    confirmation_sent_at TIMESTAMPTZ,
    recovery_token       VARCHAR(255) DEFAULT '',
    recovery_sent_at     TIMESTAMPTZ,
    last_sign_in_at      TIMESTAMPTZ,
    first_name           VARCHAR(255) DEFAULT '',
    last_name            VARCHAR(255) DEFAULT '',
    avatar_url           TEXT DEFAULT '',
    display_name         VARCHAR(255) DEFAULT '',
    is_active            BOOLEAN DEFAULT TRUE,
    app_metadata         JSONB DEFAULT '{}'::jsonb,
    user_metadata        JSONB DEFAULT '{}'::jsonb,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;

-- OAuth identities
CREATE TABLE identities (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider        VARCHAR(255) NOT NULL,
    provider_id     VARCHAR(255) NOT NULL,
    identity_data   JSONB DEFAULT '{}'::jsonb,
    last_sign_in_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_identities_provider ON identities(provider, provider_id);
CREATE INDEX idx_identities_user_id ON identities(user_id);

-- Sessions
CREATE TABLE sessions (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip         INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    not_after  TIMESTAMPTZ
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- Refresh tokens with rotation support
CREATE TABLE refresh_tokens (
    id         BIGSERIAL PRIMARY KEY,
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token      VARCHAR(255) NOT NULL,
    parent     VARCHAR(255),
    revoked    BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_session_id ON refresh_tokens(session_id);

-- One-time tokens (confirmation, recovery, magic link)
CREATE TABLE one_time_tokens (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  VARCHAR(255) NOT NULL,
    token_type  VARCHAR(50) NOT NULL,
    relates_to  VARCHAR(255) DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_one_time_tokens_user_id ON one_time_tokens(user_id);
CREATE INDEX idx_one_time_tokens_relates_to ON one_time_tokens(relates_to);

-- Updated_at trigger
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_users BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_sessions BEFORE UPDATE ON sessions FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_timestamp_refresh_tokens BEFORE UPDATE ON refresh_tokens FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
