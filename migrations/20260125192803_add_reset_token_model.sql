-- Create "password_reset_tokens" table
CREATE TABLE "password_reset_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "token_hash" text NOT NULL,
  "expires_at" timestamptz NULL,
  "used_at" timestamptz NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_password_reset_tokens_token_hash" UNIQUE ("token_hash"),
  CONSTRAINT "fk_users_password_reset_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_password_reset_tokens_deleted_at" to table: "password_reset_tokens"
CREATE INDEX "idx_password_reset_tokens_deleted_at" ON "password_reset_tokens" ("deleted_at");
-- Create index "idx_password_reset_tokens_user_id" to table: "password_reset_tokens"
CREATE INDEX "idx_password_reset_tokens_user_id" ON "password_reset_tokens" ("user_id");
