-- +goose Up
-- organization table
CREATE TABLE IF NOT EXISTS core.organizations(
    organization_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id INTEGER NOT NULL REFERENCES core.countries(country_id),
    state_id INTEGER NOT NULL REFERENCES core.states(state_id),
    theme_id UUID NOT NULL REFERENCES core.frontend_themes(theme_id),
    phone VARCHAR(20) NOT NULL UNIQUE,
    email extensions.citext NOT NULL UNIQUE,
    organization_name TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    logo_url TEXT UNIQUE,
    is_whitelabel BOOLEAN NOT NULL DEFAULT false, -- refer to partner_assignments below
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

 -- organization partner assignments (for whitelabels)
CREATE TABLE IF NOT EXISTS core.organization_partner_assignments(
    organization_id UUID PRIMARY KEY REFERENCES core.organizations(organization_id),
    partner_id UUID NOT NULL REFERENCES core.partners(partner_id),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- teams inside organizations
-- this table maps with the `user_groups` table from the Apache Guacamole schema
CREATE TABLE IF NOT EXISTS core.organization_teams(
    team_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES core.organizations(organization_id),
    team TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    UNIQUE (organization_id, team)
);


-- +goose Down
DROP TABLE IF EXISTS core.organization_teams;
DROP TABLE IF EXISTS core.organization_partner_assignments;
DROP TABLE IF EXISTS core.organizations;