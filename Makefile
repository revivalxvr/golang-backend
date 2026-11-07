# ============================
# MTLS Certificate Generator
# ============================

CERT_DIR=certs

CA_KEY=$(CERT_DIR)/ca.key
CA_CERT=$(CERT_DIR)/ca.crt
CA_SUBJ="/C=ID/ST=Jakarta/L=Jakarta/O=Dummmy/OU=IT Department/CN=Root CA"

SERVER_KEY=$(CERT_DIR)/server.key
SERVER_CSR=$(CERT_DIR)/server.csr
SERVER_CERT=$(CERT_DIR)/server.crt
SERVER_EXT=$(CERT_DIR)/server-ext.cnf
SERVER_SUBJ="/C=ID/ST=Jakarta/L=Jakarta/O=Dummmy/OU=Server/CN=server"

CLIENT_KEY=$(CERT_DIR)/client.key
CLIENT_CSR=$(CERT_DIR)/client.csr
CLIENT_CERT=$(CERT_DIR)/client.crt
CLIENT_SUBJ="/C=ID/ST=Jakarta/L=Jakarta/O=Dummmy/OU=Client/CN=client"

# ============================
# Ensure folder exists
# ============================
$(CERT_DIR):
	mkdir -p $(CERT_DIR)

# ============================
# Main Commands
# ============================
all: $(CERT_DIR) ca server client

# ============================
# Generate CA
# ============================
ca: $(CERT_DIR)
	openssl genrsa -out $(CA_KEY) 4096
	openssl req -x509 -new -nodes -key $(CA_KEY) -sha256 -days 3650 -out $(CA_CERT) \
		-subj $(CA_SUBJ)
	@echo "==> CA certificate generated at: $(CA_CERT)"

# ============================
# Generate Server Certificate
# ============================
server: $(CERT_DIR)
	openssl genrsa -out $(SERVER_KEY) 2048
	openssl req -new -key $(SERVER_KEY) -out $(SERVER_CSR) -subj $(SERVER_SUBJ)

	@echo "subjectAltName = @alt_names" >  $(SERVER_EXT)
	@echo "[alt_names]"               >> $(SERVER_EXT)
	@echo "DNS.1 = server"         	  >> $(SERVER_EXT)
	@echo "DNS.2 = localhost"         >> $(SERVER_EXT)
	@echo "DNS.3 = 127.0.0.1"         >> $(SERVER_EXT)

	openssl x509 -req -in $(SERVER_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial \
		-out $(SERVER_CERT) -days 365 -sha256 -extfile $(SERVER_EXT)

	@echo "==> Server certificate generated at: $(SERVER_CERT)"

# ============================
# Generate Client Certificate
# ============================
client: $(CERT_DIR)
	openssl genrsa -out $(CLIENT_KEY) 2048
	openssl req -new -key $(CLIENT_KEY) -out $(CLIENT_CSR) -subj $(CLIENT_SUBJ)

	openssl x509 -req -in $(CLIENT_CSR) -CA $(CA_CERT) -CAkey $(CA_KEY) -CAcreateserial \
		-out $(CLIENT_CERT) -days 365 -sha256

	@echo "==> Client certificate generated at: $(CLIENT_CERT)"

# ============================
# Clean
# ============================
clean:
	rm -rf $(CERT_DIR)
	@echo "==> All certificates removed (dir deleted)."
