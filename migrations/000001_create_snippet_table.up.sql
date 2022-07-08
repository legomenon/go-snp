CREATE TABLE IF NOT EXISTS "snippets" (
  "id" SERIAL PRIMARY KEY,
  "title" varchar(100),
  "content" text,
  "date_created" timestamp DEFAULT (now()),
  "expires" timestamp 
);