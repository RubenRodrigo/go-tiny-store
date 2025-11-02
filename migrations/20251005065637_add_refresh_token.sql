-- Create "refresh_tokens" table
CREATE TABLE "refresh_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "token" text NOT NULL,
  "revoked" boolean NULL DEFAULT false,
  "expires_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_refresh_tokens_token" UNIQUE ("token"),
  CONSTRAINT "fk_users_refresh_token" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_refresh_tokens_deleted_at" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_deleted_at" ON "refresh_tokens" ("deleted_at");
-- Create index "idx_refresh_tokens_user_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_id" ON "refresh_tokens" ("user_id");
