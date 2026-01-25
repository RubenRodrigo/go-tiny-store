-- Modify "password_reset_tokens" table
ALTER TABLE "password_reset_tokens" DROP COLUMN "deleted_at";
-- Modify "refresh_tokens" table
ALTER TABLE "refresh_tokens" DROP COLUMN "deleted_at";
