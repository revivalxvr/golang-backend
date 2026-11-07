-- 1. Buat Tabel Categories
CREATE TABLE "categories" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) NOT NULL,
    "deskripsi" VARCHAR(255) NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 2. Buat Tabel Products (category_id dibuat nullable & ON DELETE SET NULL)
CREATE TABLE "products" (
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "price" DECIMAL NOT NULL,
    "sku" VARCHAR(255) NOT NULL,
    "image_id" VARCHAR(255) NOT NULL,
    "image" VARCHAR(500) NOT NULL,
    "stock" INT NOT NULL,
    "category_id" INT NULL, 
    "length" FLOAT NOT NULL,
    "width" FLOAT NOT NULL,
    "height" FLOAT NOT NULL,
    "deskripsi" VARCHAR(255) NOT NULL,
    "badge" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_products_categories FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE RESTRICT
);