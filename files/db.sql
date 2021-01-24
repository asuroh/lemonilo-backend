
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE roles AS ENUM ('admin', 'user');

DROP TABLE IF EXISTS "user";
CREATE TABLE "user" (
  "id" char(36) PRIMARY KEY DEFAULT uuid_generate_v4 (),
  "name" VARCHAR(30),
  "email" VARCHAR(50),
  "address" VARCHAR(30), 
  "username" VARCHAR(50) not null,
  "password" VARCHAR(100) not null,
  "image_path" text,
  "role" roles,
  "created_at" timestamp(6) DEFAULT now(),
  "updated_at" timestamp(6) DEFAULT now(),
  "deleted_at" timestamp(6)
);

INSERT INTO "public"."user" ("id", "name", "email",  "address",  "username", "password", "role", "image_path", "created_at", "updated_at", "deleted_at") VALUES
('0c334e29-609a-4c7d-856b-f004a1399b3b', 'admin', 'admin@lemonilo.com', 'jl. test', 'admin007', '$2a$14$Zr20eYNB/D/q0InBRVcfgunahuUodXyLqCSJZE.FYONkBe6QXtUVm', 'admin', '/c00p9aevvhfkgt1rtjg0_WhatsApp-Image-2021-01-14-at-19.28.21.jpeg', '2021-01-15 12:07:40.845362', '2021-01-15 13:12:09.801545', NULL);

