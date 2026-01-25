-- Modify "refresh_tokens" table
ALTER TABLE "refresh_tokens" DROP CONSTRAINT "fk_users_refresh_token", DROP COLUMN "revoked", ADD CONSTRAINT "fk_users_refresh_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
