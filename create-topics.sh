#!/bin/bash

echo "Waiting for Kafka to be ready..."
sleep 10

echo "Creating Kafka topics..."

# Create main topics
kafka-topics --create --if-not-exists --topic user.created --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic user.updated --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic user.deleted --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092

kafka-topics --create --if-not-exists --topic booking.created --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic booking.updated --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic booking.deleted --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092

kafka-topics --create --if-not-exists --topic unit.created --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic unit.updated --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic unit.deleted --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092

# Create DLT topics
kafka-topics --create --if-not-exists --topic user.created.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic user.updated.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic user.deleted.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092

kafka-topics --create --if-not-exists --topic booking.created.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic booking.updated.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic booking.deleted.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092

kafka-topics --create --if-not-exists --topic unit.created.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic unit.updated.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092
kafka-topics --create --if-not-exists --topic unit.deleted.dlt --partitions 1 --replication-factor 1 --bootstrap-server localhost:9092

echo "All topics created successfully!"
echo "Listing all topics:"
kafka-topics --list --bootstrap-server localhost:9092 