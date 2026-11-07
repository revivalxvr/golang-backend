#=========================== Setup Redis Cluster ===========================

redis-cli -a rahasia --cluster create redis-node-1:6379 redis-node-2:6379 redis-node-3:6379 redis-node-4:6379 redis-node-5:6379 redis-node-6:6379 --cluster-replicas 1
redis-cli -a rahasia cluster info
redis-cli -a rahasia cluster nodes

# single node
redis-benchmark -n 2000000 -c 500 -t set,get -r 1000000

# cluster
redis-benchmark -a rahasia -n 2000000 -c 500 -t set,get -r 1000000 --cluster

#=========================== Setup Elastic ===========================

# 1. Create Policy
curl -X PUT "http://localhost:9200/_ilm/policy/dummy_log_policy" \
  -H "Content-Type: application/json" \
  -d '{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_primary_shard_size": "10mb",
            "max_age": "1m"
          },
          "set_priority": {
            "priority": 100
          }
        }
      },
      "warm": {
        "min_age": "3m",
        "actions": {
          "allocate": {
            "include": {
              "data": "warm"
            }
          },
          "set_priority": {
            "priority": 50
          }
        }
      },
      "cold": {
        "min_age": "5m",
        "actions": {
          "allocate": {
            "include": {
              "data": "cold"
            }
          },
          "set_priority": {
            "priority": 0
          }
        }
      },
      "delete": {
        "min_age": "7m",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}'

curl -X PUT http://localhost:9200/_ilm/policy/restful_log_policy -H "Content-Type: application/json" -d'
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_size": "100mb",
            "max_age": "1m"
          },
          "set_priority": { "priority": 100 }
        }
      },
      "warm": {
        "min_age": "3m",
        "actions": {
          "allocate": { "include": { "data": "warm" } },
          "set_priority": { "priority": 50 }
        }
      },
      "cold": {
        "min_age": "5m",
        "actions": {
          "allocate": { "include": { "data": "cold" } },
          "set_priority": { "priority": 0 }
        }
      },
      "delete": {
        "min_age": "7m",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}'

curl -X PUT http://localhost:9200/_ilm/policy/grpc_log_policy -H "Content-Type: application/json" -d'
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_size": "100mb",
            "max_age": "1m"
          },
          "allocate": { "include": { "data": "hot" } },
          "set_priority": { "priority": 100 }
        }
      },
      "warm": {
        "min_age": "3m",
        "actions": {
          "allocate": { "include": { "data": "warm" } },
          "set_priority": { "priority": 50 }
        }
      },
      "cold": {
        "min_age": "5m",
        "actions": {
          "allocate": { "include": { "data": "cold" } },
          "set_priority": { "priority": 0 }
        }
      },
      "delete": {
        "min_age": "7m",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}'

curl -X PUT http://localhost:9200/_ilm/policy/rabbitmq_log_policy -H "Content-Type: application/json" -d'
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_size": "100mb",
            "max_age": "1m"
          },
          "set_priority": { "priority": 100 }
        }
      },
      "warm": {
        "min_age": "3m",
        "actions": {
          "allocate": { "include": { "data": "warm" } },
          "set_priority": { "priority": 50 }
        }
      },
      "cold": {
        "min_age": "5m",
        "actions": {
          "allocate": { "include": { "data": "cold" } },
          "set_priority": { "priority": 0 }
        }
      },
      "delete": {
        "min_age": "7m",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}'

# 2. Create Index Template
curl -X PUT http://localhost:9200/_index_template/dummy_log_template -H "Content-Type: application/json" -d'
{
  "index_patterns": ["dummy-log*"],
  "data_stream": {},
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0,
      "index.lifecycle.name": "dummy_log_policy",
      "index.refresh_interval": "30s"
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "request_id": { "type": "keyword" },
        "method":     { "type": "keyword" },
        "url":        { "type": "keyword" },
        "body":       { "type": "object", "enabled": false },
        "status_code": { "type": "integer" },
        "client_ip":   { "type": "ip" },
        "user_agent":  { "type": "text" },
        "user_id":     { "type": "keyword" },
        "latency":     { "type": "long" },
        "error_message": { "type": "text" }
      }
    }
  }
}'

curl -X PUT http://localhost:9200/_index_template/restful_log_template -H "Content-Type: application/json" -d'
{
  "index_patterns": ["restful-log*"],
  "data_stream": {},
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0,
      "index.lifecycle.name": "restful_log_policy",
      "index.refresh_interval": "30s"
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "request_id": { "type": "keyword" },
        "method":     { "type": "keyword" },
        "url":        { "type": "keyword" },
        "body":       { "type": "object", "enabled": false },
        "status_code": { "type": "integer" },
        "client_ip":   { "type": "ip" },
        "user_agent":  { "type": "text" },
        "user_id":     { "type": "keyword" },
        "latency":     { "type": "long" },
        "error": { "type": "text" }
      }
    }
  }
}'

curl -X PUT http://localhost:9200/_index_template/grpc_log_template -H "Content-Type: application/json" -d'
{
  "index_patterns": ["grpc-log*"],
  "data_stream": {},
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0,
      "index.lifecycle.name": "grpc_log_policy",
      "index.refresh_interval": "30s"
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "request_id": { "type": "keyword" },
        "method":     { "type": "keyword" },
        "body":       { "type": "object", "enabled": false },
        "status_code": { "type": "integer" },
        "app_id":      { "type": "keyword" },
        "user_id":     { "type": "keyword" },
        "latency":     { "type": "long" },
        "error": { "type": "text" }
      }
    }
  }
}'

curl -X PUT http://localhost:9200/_index_template/rabbitmq_log_template -H "Content-Type: application/json" -d'
{
  "index_patterns": ["rabbitmq-log*"],
  "data_stream": {},
  "template": {
    "settings": {
      "number_of_shards": 1,
      "number_of_replicas": 0,
      "index.lifecycle.name": "rabbitmq_log_policy",
      "index.refresh_interval": "30s"
    },
    "mappings": {
      "properties": {
        "@timestamp": { "type": "date" },
        "request_id": { "type": "keyword" },
        "exchange":     { "type": "keyword" },
        "routing_key":   { "type": "keyword" },
        "queue":        { "type": "keyword" },
        "consumer_tag":  { "type": "keyword" },
        "payload":       { "type": "object", "enabled": false },
        "content_type":  { "type": "keyword" },
        "delivery_mode": { "type": "integer" },
        "app_id":        { "type": "keyword" },
        "acked":         { "type": "boolean" },
        "latency":     { "type": "long" },
        "error": { "type": "text" }
      }
    }
  }
}'

# 3. Create Data Stream (replaces initial index creation)
curl -X PUT http://localhost:9200/_data_stream/dummy-log
curl -X PUT http://localhost:9200/_data_stream/restful-log
curl -X PUT http://localhost:9200/_data_stream/grpc-log
curl -X PUT http://localhost:9200/_data_stream/rabbitmq-log

# 4. Set polling interval *for testing*
curl -X PUT http://localhost:9200/_cluster/settings -H "Content-Type: application/json" -d'
{
  "transient": {
    "indices.lifecycle.poll_interval": "30s"
  }
}'

# Check data stream exists
curl -X GET "http://localhost:9200/_data_stream/dummy-log?pretty"
curl -X GET "http://localhost:9200/_data_stream/restful-log?pretty"
curl -X GET "http://localhost:9200/_data_stream/grpc-log?pretty"
curl -X GET "http://localhost:9200/_data_stream/rabbitmq-log?pretty"

# Sample document search
curl -X GET "http://localhost:9200/dummy-log/_search?pretty"
curl -X GET "http://localhost:9200/restful-log/_search?pretty"
curl -X GET "http://localhost:9200/grpc-log/_search?pretty"
curl -X GET "http://localhost:9200/rabbitmq-log/_search?pretty"

# Check policy
curl -X GET http://localhost:9200/_ilm/policy/dummy_log_policy | jq
curl -X GET http://localhost:9200/_ilm/policy/restful_log_policy | jq
curl -X GET http://localhost:9200/_ilm/policy/grpc_log_policy | jq
curl -X GET http://localhost:9200/_ilm/policy/rabbitmq_log_policy | jq

# Check index template
curl -X GET "http://localhost:9200/_index_template/dummy_log_template" | jq
curl -X GET "http://localhost:9200/_index_template/restful_log_template" | jq
curl -X GET "http://localhost:9200/_index_template/grpc_log_template" | jq
curl -X GET "http://localhost:9200/_index_template/rabbitmq_log_template" | jq

# Check detail index
GET http://localhost:9200/dummy-log/_ilm/explain | jq
GET http://localhost:9200/restful-log/_ilm/explain | jq
GET http://localhost:9200/grpc-log/_ilm/explain | jq
GET http://localhost:9200/rabbitmq-log/_ilm/explain | jq

# Check all index
curl -X GET "http://localhost:9200/_cat/indices?v"

# Check node attributes
GET http://localhost:9200/_nodes?filter_path=**.attributes | jq

# Delete
curl -X DELETE http://localhost:9200/dummy-log
curl -X DELETE http://localhost:9200/restful-log
curl -X DELETE http://localhost:9200/grpc-log
curl -X DELETE http://localhost:9200/rabbitmq-log

# Delete data stream
curl -X DELETE http://localhost:9200/_data_stream/dummy-log
curl -X DELETE http://localhost:9200/_data_stream/restful-log
curl -X DELETE http://localhost:9200/_data_stream/grpc-log
curl -X DELETE http://localhost:9200/_data_stream/rabbitmq-log

# Delete index template
curl -X DELETE http://localhost:9200/_index_template/dummy_log_template
curl -X DELETE http://localhost:9200/_index_template/restful_log_template
curl -X DELETE http://localhost:9200/_index_template/grpc_log_template
curl -X DELETE http://localhost:9200/_index_template/rabbitmq_log_template

# Delete policy
curl -X DELETE http://localhost:9200/_ilm/policy/dummy_log_policy
curl -X DELETE http://localhost:9200/_ilm/policy/restful_log_policy
curl -X DELETE http://localhost:9200/_ilm/policy/grpc_log_policy
curl -X DELETE http://localhost:9200/_ilm/policy/rabbitmq_log_policy