CREATE TABLE IF NOT EXISTS "orders" (
    "id" VARCHAR(36) PRIMARY KEY,
    "user_id" VARCHAR(36) NOT NULL,
    "gross_amount" INTEGER NOT NULL,
    "status" VARCHAR(255) NOT NULL,
    "payment_method" VARCHAR(255) NOT NULL,
    "payment_url" VARCHAR(500) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);