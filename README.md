# Go Server for Communicating with Cassandra

This repository contains a Go server implementation for communicating with a Cassandra database. The Cassandra database is assumed to be running on Docker.

## Cassandra Setup

Ensure that Cassandra is running on Docker before running the Go server. The Cassandra database contains half a million records, which will be divided into two nodes for communication.

## Go Server

The Go server acts as an intermediary between the Cassandra database and clients. It uses the Cassandra Query Language (CQL) to query and manipulate data in the Cassandra database.

## Usage

1. **Start Cassandra on Docker:**
   - Run Cassandra on Docker with appropriate configurations to divide the data into two nodes.

2. **Run the Go Server:**
   - Run the Go server to start communicating with the Cassandra database.
   - Ensure that the server is configured to connect to the Cassandra cluster.

3. **Accessing Data:**
   - Once the Go server is running, clients can access data from the Cassandra database by sending requests to the server.
   - The server handles client requests and interacts with the Cassandra database using CQL queries.

## Dependencies

Ensure that the following dependencies are installed and configured:

- Go programming language
- Docker with Cassandra image

## Configuration

- Update the server configuration file (`config.go`) with appropriate settings for connecting to the Cassandra database.
- Modify any other configuration files as needed to match your environment and requirements.

## Contributions

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

