# eventctl

A command-line tool for inspecting and interacting with Kafka topics. Useful for debugging event-driven systems — peek at recent messages, publish test events, and manage topics without leaving the terminal.

## Installation

```bash
go install github.com/faridtmammadov/eventctl@latest
```

Or build from source:

```bash
git clone https://github.com/faridtmammadov/eventctl
cd eventctl
go build -o eventctl .
```

## Configuration

Create `~/.eventctl/config.yaml` to define named broker connections:

```yaml
connections:
  local-kafka:
    brokers:
      - localhost:9092

  staging-kafka:
    brokers:
      - kafka1.staging.example.com:9092
      - kafka2.staging.example.com:9092

default: local-kafka
```

Copy the provided example as a starting point:

```bash
cp config.example.yaml ~/.eventctl/config.yaml
```

The `default` key sets which connection is used when no `--connection` flag is passed. You can also skip the config file entirely and point directly at a broker with `--broker`.

## Local Kafka with Redpanda

The repo includes a `docker-compose.yaml` that starts a single-node [Redpanda](https://redpanda.com) instance — a Kafka-compatible broker ideal for local experimentation:

```bash
docker compose up -d
```

Redpanda will be available at `localhost:19092`. Update your config to match:

```yaml
connections:
  local-kafka:
    brokers:
      - localhost:19092

default: local-kafka
```

## Commands

### peek — read recent messages from a topic

```bash
# Print the last message from a topic (uses default connection)
eventctl peek orders

# Print the last 5 messages
eventctl peek orders -n 5

# Use a named connection from config
eventctl peek orders -c staging-kafka -n 10

# Point directly at a broker without a config file
eventctl peek orders --broker localhost:19092
```

### publish — send a message to a topic

```bash
# Publish an inline message
eventctl publish orders --message '{"id": 1, "status": "new"}'

# Publish with an explicit key
eventctl publish orders --key order-123 --message '{"id": 1, "status": "new"}'

# Pipe from stdin
echo '{"id": 1, "status": "new"}' | eventctl publish orders

# Publish to a staging broker
eventctl publish orders -c staging-kafka --message '{"id": 2}'
```

### topic — manage topics

```bash
# Create a topic with default settings (1 partition, replication factor 1)
eventctl topic create orders

# Create a topic with multiple partitions
eventctl topic create orders --partitions 6

# Create a topic on a specific broker
eventctl --broker localhost:19092 topic create shipping --partitions 3 --replication-factor 1
```

## Global flags

| Flag | Short | Description |
|---|---|---|
| `--connection` | `-c` | Named connection from `~/.eventctl/config.yaml` |
| `--broker` | `-b` | Broker address(es), comma-separated. Overrides config. |
