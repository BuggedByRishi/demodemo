-- +goose Up
-- agents (employees) table
-- this table maps with the `users` table from the Apache Guacamole schema
CREATE TABLE IF NOT EXISTS core.agents(
    agent_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES core.organizations(organization_id),
    agent_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email extensions.citext NOT NULL,
    guacamole_password BYTEA NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    UNIQUE (organization_id, phone),
    UNIQUE (organization_id, email)
);

-- agent team assignments for the organization
CREATE TABLE IF NOT EXISTS core.agent_teams(
    agent_id UUID NOT NULL REFERENCES core.agents(agent_id),
    team_id UUID NOT NULL REFERENCES core.organization_teams(team_id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (agent_id, team_id)
);

-- agent oauth associations
-- this association is created for the oauth providers
-- so the main database is connected to the auth database via this association
CREATE TABLE IF NOT EXISTS core.agent_oauth_associations(
    agent_id UUID PRIMARY KEY REFERENCES core.agents(agent_id),
    user_id UUID NOT NULL UNIQUE REFERENCES auth.users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS core.agent_oauth_associations;
DROP TABLE IF EXISTS core.agent_teams;
DROP TABLE IF EXISTS core.agents;