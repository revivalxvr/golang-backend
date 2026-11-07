CONF_DIR=./conf
PREFIX=redis-node
PORT=6379
BIND_IP=0.0.0.0
CLUSTER_ENABLED=yes
REQUIREPASS=rahasia
MASTERAUTH=rahasia
CLUSTER_NODE_TIMEOUT=5000
CLUSTER_ANNOUNCE_PORT=6379
CLUSTER_ANNOUNCE_BUS_PORT=16379
APPENDONLY=yes
NODES=1 2 3 4 5 6

.PHONY: all-redis-conf redis-conf clean-redis-conf

all-redis-conf: clean-redis-conf redis-conf

redis-conf: ${CONF_DIR}
	@for n in $(NODES); do \
		echo "Creating ${CONF_DIR}/${PREFIX}-$$n.conf"; \
		echo "port ${PORT}" > ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "bind ${BIND_IP}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "protected-mode no" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "cluster-enabled ${CLUSTER_ENABLED}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "requirepass ${REQUIREPASS}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "masterauth ${MASTERAUTH}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "cluster-config-file nodes-$$n.conf" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "cluster-node-timeout ${CLUSTER_NODE_TIMEOUT}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "cluster-announce-ip 172.48.0.1$$n" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "cluster-announce-port ${CLUSTER_ANNOUNCE_PORT}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "cluster-announce-bus-port ${CLUSTER_ANNOUNCE_BUS_PORT}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
		echo "appendonly ${APPENDONLY}" >> ${CONF_DIR}/${PREFIX}-$$n.conf; \
	done

${CONF_DIR}:
	mkdir -p ${CONF_DIR}

clean-redis-conf:
	rm -rf ${CONF_DIR}