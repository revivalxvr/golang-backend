CREATE TABLE IF NOT EXISTS "credentials" (
    "user_id" UUID PRIMARY KEY,
    "password" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "fk_credentials_user"
        FOREIGN KEY ("user_id") 
        REFERENCES "users" ("id") 
        ON DELETE CASCADE
)