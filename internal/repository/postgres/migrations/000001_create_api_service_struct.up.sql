CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE statuses AS ENUM ('OPEN', 'MERGED');

-- 1. Таблица для команд (TeamDBModel)
CREATE TABLE IF NOT EXISTS teams (
  team_name VARCHAR(255) PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE
);

-- 2. Таблица для пользователей (UserDBModel)
CREATE TABLE IF NOT EXISTS users (
  user_id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  username VARCHAR(255) NOT NULL,
  team_name VARCHAR(255) NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE,

  CONSTRAINT fk_team
  FOREIGN KEY(team_name)
  REFERENCES teams(team_name)
  ON DELETE RESTRICT
);

-- Индекс для быстрого поиска активных пользователей в команде (для логики назначения ревьюверов)
CREATE INDEX IF NOT EXISTS idx_users_team_active ON users (team_name, is_active) WHERE is_active = TRUE;

-- 3. Таблица для Pull Request'ов (PullRequestDBModel)
CREATE TABLE IF NOT EXISTS pull_requests (
  pull_request_id VARCHAR(255) PRIMARY KEY,
  pull_request_name VARCHAR(255) NOT NULL,
  author_id UUID NOT NULL,
  status statuses NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  merged_at TIMESTAMP WITH TIME ZONE,

  CONSTRAINT fk_author
  FOREIGN KEY(author_id)
  REFERENCES users(user_id)
  ON DELETE RESTRICT
);

-- Индекс для быстрого поиска PR по автору
CREATE INDEX IF NOT EXISTS idx_pr_author ON pull_requests (author_id);

-- 4. Таблица для связи PR и ревьюверов (PullRequestReviewerDBModel)
CREATE TABLE IF NOT EXISTS pr_reviewers (
  pull_request_id VARCHAR(255) NOT NULL,
  reviewer_id UUID NOT NULL,
  assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  PRIMARY KEY (pull_request_id, reviewer_id),

  CONSTRAINT fk_pr
  FOREIGN KEY(pull_request_id)
  REFERENCES pull_requests(pull_request_id)
  ON DELETE CASCADE, 

  CONSTRAINT fk_reviewer
  FOREIGN KEY(reviewer_id)
  REFERENCES users(user_id)
  ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_reviewer_pr ON pr_reviewers (reviewer_id);

-- 5. Таблица для токенов аутентификации (AuthTokenDBModel)
CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY DEFAULT (gen_random_uuid()),
  user_id UUID NOT NULL,
  token text UNIQUE NOT NULL,
  expires_at timestamp NOT NULL,
  created_at timestamp NOT NULL DEFAULT (now()),

  CONSTRAINT fk_token_user
  FOREIGN KEY(user_id)
  REFERENCES users(user_id)
  ON DELETE CASCADE
);

-- Дополнительно: Добавление начальных данных для админа
-- INSERT INTO teams (team_name) VALUES ('system') ON CONFLICT DO NOTHING;
-- INSERT INTO users (user_id, username, team_name, is_active, is_admin) VALUES ('00000000-0000-0000-0000-000000000000', 'System Admin', 'system', TRUE, TRUE) ON CONFLICT DO NOTHING;