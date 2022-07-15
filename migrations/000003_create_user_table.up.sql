CREATE TABLE IF NOT EXISTS "users" (
	"id" SERIAL PRIMARY KEY NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"email" VARCHAR(255) NOT NULL,
	"hashed_password" CHAR(60) NOT NULL,
	"created" TIMESTAMP NOT NULL DEFAULT (now())
);

ALTER TABLE "users" ADD CONSTRAINT "users_uc_email" UNIQUE (email);

