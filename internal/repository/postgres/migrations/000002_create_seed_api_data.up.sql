TRUNCATE TABLE pr_reviewers, pull_requests, sessions, users, teams RESTART IDENTITY CASCADE;

INSERT INTO teams (team_name)
VALUES
    ('backend'),
    ('frontend'),
    ('mobile');

INSERT INTO users (username, team_name, is_active, is_admin)
VALUES
    ('admin_backend', 'backend', TRUE, TRUE);

INSERT INTO users (username, team_name, is_active)
VALUES
    ('alice', 'backend', TRUE),
    ('bob', 'backend', TRUE),
    ('charlie', 'backend', TRUE),
    ('denis', 'backend', TRUE),
    ('igor', 'backend', TRUE),
    ('kate', 'backend', TRUE),
    ('leo', 'backend', FALSE);

INSERT INTO users (username, team_name, is_active)
VALUES
    ('mike', 'frontend', TRUE),
    ('nina', 'frontend', TRUE),
    ('olga', 'frontend', TRUE),
    ('pavel', 'frontend', FALSE),
    ('roma', 'frontend', TRUE),
    ('sofia', 'frontend', TRUE);

INSERT INTO users (username, team_name, is_active)
VALUES
    ('tanya', 'mobile', TRUE),
    ('vlad', 'mobile', TRUE),
    ('yana', 'mobile', TRUE),
    ('zoya', 'mobile', TRUE),
    ('kirill', 'mobile', FALSE);

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
SELECT
    'pr-1001',
    'Add search endpoint',
    u.user_id,
    'OPEN'
FROM users u WHERE u.username = 'alice';

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
SELECT
    'pr-1002',
    'Fix login handler',
    u.user_id,
    'OPEN'
FROM users u WHERE u.username = 'bob';

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
SELECT
    'pr-1003',
    'Refactor caching',
    u.user_id,
    'MERGED'
FROM users u WHERE u.username = 'charlie';

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
SELECT
    'pr-1004',
    'Optimize DB queries',
    u.user_id,
    'OPEN'
FROM users u WHERE u.username = 'denis';

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
SELECT
    'pr-1005',
    'Implement GraphQL layer',
    u.user_id,
    'OPEN'
FROM users u WHERE u.username = 'mike';

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
SELECT
    'pr-1006',
    'Update mobile UI',
    u.user_id,
    'MERGED'
FROM users u WHERE u.username = 'tanya';

INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
SELECT 'pr-1001', u.user_id FROM users u WHERE u.username = 'bob';

INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
SELECT 'pr-1002', u.user_id FROM users u WHERE u.username = 'denis';

INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
SELECT 'pr-1003', u.user_id FROM users u WHERE u.username = 'igor';

INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
SELECT 'pr-1004', u.user_id FROM users u WHERE u.username = 'kate';

INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
SELECT 'pr-1005', u.user_id FROM users u WHERE u.username = 'nina';
