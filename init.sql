-- -------------------------------------------------------------
-- TablePlus 5.9.0(538)
--
-- https://tableplus.com/
--
-- Database: test
-- Generation Time: 2024-04-01 17:38:48.7520
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "public"."posts";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."posts" (
    "id" uuid NOT NULL DEFAULT uuid7(),
    "title" text,
    "body" text,
    "created" timestamp NOT NULL DEFAULT now(),
    "updated" timestamp NOT NULL DEFAULT now(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    PRIMARY KEY ("id")
);

DROP TABLE IF EXISTS "public"."users";
-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Table Definition
CREATE TABLE "public"."users" (
    "id" uuid NOT NULL DEFAULT uuid7(),
    "created_at" timestamptz,
    "updated_at" timestamptz,
    "deleted_at" timestamptz,
    "email" text,
    "password" text,
    "first_name" text,
    "last_name" text,
    "verified" bool DEFAULT false,
    "code" text,
    PRIMARY KEY ("id")
);

