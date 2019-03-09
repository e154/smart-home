-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE roles (
  name        VARCHAR(255) NOT NULL PRIMARY KEY,
  description VARCHAR(255) NOT NULL,
  parent      VARCHAR(255) NULL CONSTRAINT parent_at_roles_fk REFERENCES roles (name) ON UPDATE CASCADE ON DELETE RESTRICT,
  created_at  TIMESTAMPTZ  NOT NULL,
  updated_at  TIMESTAMPTZ  NOT NULL
);

CREATE TABLE permissions (
  id           BIGSERIAL PRIMARY KEY,
  role_name    varchar(255) NOT NULL CONSTRAINT role_at_permissions_fk REFERENCES roles (name) ON UPDATE CASCADE ON DELETE CASCADE,
  package_name varchar(255) NOT NULL,
  level_name   varchar(255) NOT NULL
);

CREATE UNIQUE INDEX permissions_unq
  ON permissions (role_name, package_name, level_name);

CREATE TABLE users (
  id                     BIGSERIAL PRIMARY KEY,
  nickname               VARCHAR(255) NOT NULL,
  first_name             VARCHAR(255) NOT NULL,
  last_name              VARCHAR(255) NOT NULL,
  encrypted_password     VARCHAR(255) NOT NULL,
  email                  VARCHAR(255) NOT NULL,
  lang                   VARCHAR(2)   NOT NULL DEFAULT 'en',
  history                JSONB                 DEFAULT '{}',
  status                 VARCHAR(255) NOT NULL DEFAULT 'active',
  reset_password_token   VARCHAR(255) NOT NULL,
  authentication_token   VARCHAR(255) NOT NULL,
  image_id               BIGINT CONSTRAINT image_at_users_fk REFERENCES images (id) ON UPDATE CASCADE ON DELETE SET NULL,
  sign_in_count          BIGINT       NOT NULL,
  current_sign_in_ip     VARCHAR(255) NOT NULL,
  last_sign_in_ip        VARCHAR(255) NULL,
  user_id                BIGINT       NULL CONSTRAINT user_at_users_fk REFERENCES users (id) ON UPDATE CASCADE ON DELETE SET NULL,
  role_name              VARCHAR(255) CONSTRAINT role_name_at_users_fk REFERENCES roles (name) ON UPDATE CASCADE ON DELETE RESTRICT,
  reset_password_sent_at TIMESTAMPTZ  NULL,
  current_sign_in_at     TIMESTAMPTZ  NULL,
  last_sign_in_at        TIMESTAMPTZ  NULL,
  created_at             TIMESTAMPTZ  NOT NULL,
  updated_at             TIMESTAMPTZ  NOT NULL,
  deleted_at             TIMESTAMPTZ  NULL
);

CREATE UNIQUE INDEX email_2_users_unq
  ON users (email);
CREATE UNIQUE INDEX nickname_2_users_unq
  ON users (nickname);
CREATE UNIQUE INDEX authentication_token_2_users_unq
  ON users (authentication_token);

CREATE INDEX email_2_user_idx
  ON users (email);
CREATE INDEX authentication_token_2_user_idx
  ON users (authentication_token);

CREATE TABLE user_metas (
  id      BIGSERIAL PRIMARY KEY,
  user_id BIGINT CONSTRAINT user_at_user_metas_fk REFERENCES users (id) ON UPDATE CASCADE ON DELETE SET NULL,
  key     VARCHAR(255) NOT NULL,
  value   text         NOT NULL
);

CREATE UNIQUE INDEX kay_user_2_user_metas_unq
  ON user_metas (key, user_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS workers CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS user_metas CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;

