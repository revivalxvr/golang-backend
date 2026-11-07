CREATE TABLE IF NOT EXISTS "refresh_tokens" (
   "id" UUID PRIMARY KEY,
    "user_id" UUID NOT NULL,
    "token" VARCHAR(512) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
   CONSTRAINT "fk_refresh_tokens_user"
        FOREIGN KEY ("user_id") 
        REFERENCES "users" ("id") 
        ON DELETE CASCADE
)