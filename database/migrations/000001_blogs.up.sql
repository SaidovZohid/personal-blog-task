CREATE TABLE IF NOT EXISTS "users" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR,
    "email" VARCHAR NOT NULL,
    "password" VARCHAR NOT NULL,
    "role" VARCHAR NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "posts" (
    "id" SERIAL PRIMARY KEY,
    "header" VARCHAR NOT NULL,
    "content" TEXT NOT NULL,
    "user_id" INTEGER NOT NULL REFERENCES users(id),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "comments" (
    "id" SERIAL PRIMARY KEY,
    "content" TEXT NOT NULL,
    "post_id" INTEGER NOT NULL REFERENCES posts(id),
    "user_id" INTEGER NOT NULL REFERENCES users(id),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "replies" (
    "id" SERIAL PRIMARY KEY,
    "content" TEXT NOT NULL,
    "comment_id" INTEGER NOT NULL REFERENCES comments(id),
    "user_id" INTEGER NOT NULL REFERENCES users(id),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);