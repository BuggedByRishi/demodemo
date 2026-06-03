-- +goose Up
--
-- search paths
--
ALTER ROLE tally_on_cloud_auth_role    SET search_path TO extensions, auth;
ALTER ROLE tally_on_cloud_app_role     SET search_path TO extensions, auth, proxmox, guacamole, core, tally;

--
-- revoke public (all users) access
--
REVOKE ALL ON SCHEMA extensions     FROM PUBLIC;
REVOKE ALL ON SCHEMA auth           FROM PUBLIC;
REVOKE ALL ON SCHEMA core           FROM PUBLIC;
REVOKE ALL ON SCHEMA proxmox        FROM PUBLIC;
REVOKE ALL ON SCHEMA guacamole      FROM PUBLIC;
REVOKE ALL ON SCHEMA tally          FROM PUBLIC;

--
-- grant schema usage
--
-- auth role can use auth and extensions schemas
GRANT USAGE ON SCHEMA extensions    TO tally_on_cloud_auth_role;
GRANT USAGE ON SCHEMA auth          TO tally_on_cloud_auth_role;
-- app role can use all schemas
GRANT USAGE ON SCHEMA extensions    TO tally_on_cloud_app_role;
GRANT USAGE ON SCHEMA auth          TO tally_on_cloud_app_role;
GRANT USAGE ON SCHEMA core          TO tally_on_cloud_app_role;
GRANT USAGE ON SCHEMA proxmox       TO tally_on_cloud_app_role;
GRANT USAGE ON SCHEMA guacamole     TO tally_on_cloud_app_role;
GRANT USAGE ON SCHEMA tally         TO tally_on_cloud_app_role;

--
-- default privileges for future objects
-- set before tables are created so all objects inherit these permissions
--
-- extensions schema: auth_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_auth_role;
-- auth schema: auth_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_auth_role;
-- extensions schema: app_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_app_role;
-- auth schema: app_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_app_role;
-- core schema: app_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA core          GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA core          GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA core          GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_app_role;
-- proxmox schema: app_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA proxmox       GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA proxmox       GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA proxmox       GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_app_role;
-- guacamole schema: app_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA guacamole     GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA guacamole     GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA guacamole     GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_app_role;
-- tally schema: app_role gets full DML
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA tally         GRANT SELECT, INSERT, UPDATE, DELETE    ON TABLES       TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA tally         GRANT USAGE, SELECT                     ON SEQUENCES    TO tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA tally         GRANT EXECUTE                           ON FUNCTIONS    TO tally_on_cloud_app_role;


-- +goose Down
-- revoke default privileges for tally schema for app role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA tally         REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA tally         REVOKE ALL ON SEQUENCES FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA tally         REVOKE ALL ON TABLES    FROM tally_on_cloud_app_role;
-- revoke default privileges for guacamole schema for app role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA guacamole     REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA guacamole     REVOKE ALL ON SEQUENCES FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA guacamole     REVOKE ALL ON TABLES    FROM tally_on_cloud_app_role;
-- revoke default privileges for proxmox schema for app role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA proxmox       REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA proxmox       REVOKE ALL ON SEQUENCES FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA proxmox       REVOKE ALL ON TABLES    FROM tally_on_cloud_app_role;
-- revoke default privileges for core schema for app role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA core          REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA core          REVOKE ALL ON SEQUENCES FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA core          REVOKE ALL ON TABLES    FROM tally_on_cloud_app_role;
-- revoke default privileges for auth schema for app role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          REVOKE ALL ON SEQUENCES FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          REVOKE ALL ON TABLES    FROM tally_on_cloud_app_role;
-- revoke default privileges for extensions schema for app role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    REVOKE ALL ON SEQUENCES FROM tally_on_cloud_app_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    REVOKE ALL ON TABLES    FROM tally_on_cloud_app_role;
-- revoke default privileges for auth schema for auth role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          REVOKE ALL ON SEQUENCES FROM tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA auth          REVOKE ALL ON TABLES    FROM tally_on_cloud_auth_role;
-- revoke default privileges for extensions schema for auth role
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    REVOKE ALL ON FUNCTIONS FROM tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    REVOKE ALL ON SEQUENCES FROM tally_on_cloud_auth_role;
ALTER DEFAULT PRIVILEGES FOR ROLE tally_on_cloud_admin IN SCHEMA extensions    REVOKE ALL ON TABLES    FROM tally_on_cloud_auth_role;
-- revoke schema usage
REVOKE USAGE ON SCHEMA tally        FROM tally_on_cloud_app_role;
REVOKE USAGE ON SCHEMA guacamole    FROM tally_on_cloud_app_role;
REVOKE USAGE ON SCHEMA proxmox      FROM tally_on_cloud_app_role;
REVOKE USAGE ON SCHEMA core         FROM tally_on_cloud_app_role;
REVOKE USAGE ON SCHEMA auth         FROM tally_on_cloud_app_role;
REVOKE USAGE ON SCHEMA extensions   FROM tally_on_cloud_app_role;
REVOKE USAGE ON SCHEMA auth         FROM tally_on_cloud_auth_role;
REVOKE USAGE ON SCHEMA extensions   FROM tally_on_cloud_auth_role;
-- reset search paths
ALTER ROLE tally_on_cloud_auth_role    RESET search_path;
ALTER ROLE tally_on_cloud_app_role     RESET search_path;