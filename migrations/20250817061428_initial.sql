-- Create "categories" table
CREATE TABLE "categories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(100) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_categories_deleted_at" to table: "categories"
CREATE INDEX "idx_categories_deleted_at" ON "categories" ("deleted_at");
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" text NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create "products" table
CREATE TABLE "products" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  "price" numeric NOT NULL,
  "disabled" boolean NULL DEFAULT false,
  "stock" bigint NULL DEFAULT 0,
  "category_id" uuid NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_categories_products" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "products" ("deleted_at");
-- Create "product_images" table
CREATE TABLE "product_images" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "product_id" uuid NOT NULL,
  "url" character varying(500) NOT NULL,
  "alt_text" character varying(255) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_images" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_product_images_deleted_at" to table: "product_images"
CREATE INDEX "idx_product_images_deleted_at" ON "product_images" ("deleted_at");
-- Create index "idx_product_images_product_id" to table: "product_images"
CREATE INDEX "idx_product_images_product_id" ON "product_images" ("product_id");
