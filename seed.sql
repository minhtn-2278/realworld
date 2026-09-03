-- Development seed data for realworldapp.
--
-- All seeded users use the password: password
-- The statements are idempotent and require all application migrations to run first.

BEGIN;

INSERT INTO users (username, email, password_hash, bio, image)
VALUES
	(
		'seed_alice',
		'alice.seed@example.test',
		'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
		'Backend developer who enjoys Go, PostgreSQL, and practical API design.',
		'https://i.pravatar.cc/300?img=47'
	),
	(
		'seed_bob',
		'bob.seed@example.test',
		'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
		'Frontend developer exploring the Realworld API.',
		'https://i.pravatar.cc/300?img=12'
	),
	(
		'seed_carol',
		'carol.seed@example.test',
		'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
		'QA engineer testing feeds, favorites, and comments.',
		'https://i.pravatar.cc/300?img=32'
	)
ON CONFLICT (username) DO UPDATE
SET
	email = EXCLUDED.email,
	password_hash = EXCLUDED.password_hash,
	bio = EXCLUDED.bio,
	image = EXCLUDED.image,
	updated_at = CURRENT_TIMESTAMP;

INSERT INTO tags (name)
VALUES
	('go'),
	('echo'),
	('postgresql'),
	('redis'),
	('swagger')
ON CONFLICT (name) DO NOTHING;

INSERT INTO articles (slug, title, description, body, author_id, created_at, updated_at)
SELECT
	'getting-started-with-go-and-echo-seed0001',
	'Getting started with Go and Echo',
	'A short introduction to building a REST API with Go and Echo.',
	'This seeded article is useful for testing article listing, article details, tags, comments, and favorites.',
	u.id,
	CURRENT_TIMESTAMP - INTERVAL '3 days',
	CURRENT_TIMESTAMP - INTERVAL '3 days'
FROM users u
WHERE u.username = 'seed_alice'
ON CONFLICT (slug) DO UPDATE
SET
	title = EXCLUDED.title,
	description = EXCLUDED.description,
	body = EXCLUDED.body,
	author_id = EXCLUDED.author_id,
	updated_at = CURRENT_TIMESTAMP;

INSERT INTO articles (slug, title, description, body, author_id, created_at, updated_at)
SELECT
	'postgresql-indexes-for-api-queries-seed0002',
	'PostgreSQL indexes for API queries',
	'Common indexes that keep list endpoints responsive as data grows.',
	'This seeded article gives seed_alice a second article and lets seed_bob test a non-empty feed.',
	u.id,
	CURRENT_TIMESTAMP - INTERVAL '2 days',
	CURRENT_TIMESTAMP - INTERVAL '2 days'
FROM users u
WHERE u.username = 'seed_alice'
ON CONFLICT (slug) DO UPDATE
SET
	title = EXCLUDED.title,
	description = EXCLUDED.description,
	body = EXCLUDED.body,
	author_id = EXCLUDED.author_id,
	updated_at = CURRENT_TIMESTAMP;

INSERT INTO articles (slug, title, description, body, author_id, created_at, updated_at)
SELECT
	'testing-rest-apis-with-realistic-data-seed0003',
	'Testing REST APIs with realistic data',
	'Repeatable data makes manual API testing faster and more reliable.',
	'This seeded article belongs to seed_bob and is visible in seed_alice''s feed after logging in.',
	u.id,
	CURRENT_TIMESTAMP - INTERVAL '1 day',
	CURRENT_TIMESTAMP - INTERVAL '1 day'
FROM users u
WHERE u.username = 'seed_bob'
ON CONFLICT (slug) DO UPDATE
SET
	title = EXCLUDED.title,
	description = EXCLUDED.description,
	body = EXCLUDED.body,
	author_id = EXCLUDED.author_id,
	updated_at = CURRENT_TIMESTAMP;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id
FROM articles a
JOIN tags t ON t.name IN ('go', 'echo', 'swagger')
WHERE a.slug = 'getting-started-with-go-and-echo-seed0001'
ON CONFLICT DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id
FROM articles a
JOIN tags t ON t.name IN ('go', 'postgresql')
WHERE a.slug = 'postgresql-indexes-for-api-queries-seed0002'
ON CONFLICT DO NOTHING;

INSERT INTO article_tags (article_id, tag_id)
SELECT a.id, t.id
FROM articles a
JOIN tags t ON t.name IN ('go', 'redis')
WHERE a.slug = 'testing-rest-apis-with-realistic-data-seed0003'
ON CONFLICT DO NOTHING;

INSERT INTO user_follows (follower_id, following_id)
SELECT follower.id, following.id
FROM users follower
JOIN users following ON following.username = 'seed_alice'
WHERE follower.username IN ('seed_bob', 'seed_carol')
ON CONFLICT DO NOTHING;

INSERT INTO user_follows (follower_id, following_id)
SELECT follower.id, following.id
FROM users follower
JOIN users following ON following.username = 'seed_bob'
WHERE follower.username = 'seed_alice'
ON CONFLICT DO NOTHING;

INSERT INTO article_favorites (article_id, user_id)
SELECT a.id, u.id
FROM articles a
JOIN users u ON u.username IN ('seed_bob', 'seed_carol')
WHERE a.slug = 'getting-started-with-go-and-echo-seed0001'
ON CONFLICT DO NOTHING;

INSERT INTO article_favorites (article_id, user_id)
SELECT a.id, u.id
FROM articles a
JOIN users u ON u.username = 'seed_alice'
WHERE a.slug = 'testing-rest-apis-with-realistic-data-seed0003'
ON CONFLICT DO NOTHING;

INSERT INTO comments (body, article_id, author_id)
SELECT
	'Great starting point. The Swagger page makes the endpoints easy to explore.',
	a.id,
	u.id
FROM articles a
JOIN users u ON u.username = 'seed_bob'
WHERE a.slug = 'getting-started-with-go-and-echo-seed0001'
	AND NOT EXISTS (
		SELECT 1
		FROM comments c
		WHERE c.body = 'Great starting point. The Swagger page makes the endpoints easy to explore.'
			AND c.article_id = a.id
			AND c.author_id = u.id
	);

INSERT INTO comments (body, article_id, author_id)
SELECT
	'Useful checklist. I will use these records while testing pagination.',
	a.id,
	u.id
FROM articles a
JOIN users u ON u.username = 'seed_carol'
WHERE a.slug = 'testing-rest-apis-with-realistic-data-seed0003'
	AND NOT EXISTS (
		SELECT 1
		FROM comments c
		WHERE c.body = 'Useful checklist. I will use these records while testing pagination.'
			AND c.article_id = a.id
			AND c.author_id = u.id
	);

COMMIT;
