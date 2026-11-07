CREATE TABLE IF NOT EXISTS "product_orders" (
    "id" UUID PRIMARY KEY,
    "order_id" VARCHAR(36) NOT NULL, -- Disamakan ukurannya menjadi VARCHAR(36)
    "product_id" UUID NOT NULL,
    "quantity" INT NOT NULL,
    "price" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
  
    CONSTRAINT fk_order 
        FOREIGN KEY (order_id) 
        REFERENCES "orders"(id) 
        ON DELETE CASCADE -- Atau NO ACTION jika Anda ingin mencegah data order terhapus saat produk terkait dihapus
);