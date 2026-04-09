# Project Guide: ROS-to-RabbitMQ Go Agent

This document provides a detailed overview of the Go Agent, its architecture, folder structure, and setup instructions. This agent is responsible for generating code for ROS topic conversions, subscribing to those topics, and forwarding the data to RabbitMQ.

---

## 🚀 Project Overview

The **Go Agent** acts as a bridge between a ROS (Robot Operating System) environment and a RabbitMQ message broker. It automates the process of:
1.  **Code Generation**: Converting ROS message definitions into Go structs and Protocol Buffers (Proto).
2.  **Subscription**: Dynamically subscribing to ROS topics.
3.  **Publishing**: Serializing ROS data and publishing it to specific RabbitMQ exchanges/queues.

---

## 🛠 Setup & Installation

### Prerequisites
-   **Go**: Version 1.22.5 or higher.
-   **ROS**: Compatible with ROS 1 (via `goroslib`). Ensure a ROS Master is running (default `127.0.0.1:11311`).
-   **RabbitMQ**: An accessible RabbitMQ instance.

### Installation
1.  Clone the repository.
2.  Install dependencies:
    ```bash
    go mod download
    ```

### Configuration
Update the YAML files in the `config/` directory:
-   **`config/rmq_config.yml`**: Configure RabbitMQ credentials, host, and vhost.
-   **`config/telemetry_node.yml`**: Configure the ROS node name and Master address.

### Running the Agent
Currently, the application can be tested using `test.go`:
```bash
go run test.go
```

---

## 🏗 Core Architecture

The agent follows a pipeline-based architecture:

1.  **Generation Phase**: The `code-generator` parses ROS `.msg` files and generates equivalent Go and Proto definitions. This ensures type safety and compatibility between ROS and external systems.
2.  **Subscription Phase**: Based on configuration, the `subscribers` module initializes ROS nodes and subscribes to specified topics using `goroslib`.
3.  **Transmission Phase**: Data received from ROS topics is passed through internal Go channels, serialized (usually to Proto), and sent to RabbitMQ via the `publishers` module.

---

## 📂 Detailed Folder Structure

| Folder | Description | Status |
| :--- | :--- | :--- |
| `api/` | Intended for external API definitions (GRPC/REST). | 🚧 Placeholder |
| `cmd/` | Main entry points for the application. | 🚧 Placeholder |
| `code-generator/` | **Heart of the agent.** Logic for parsing ROS messages and generating Go/Proto code. | ✅ Active |
| `config/` | YAML and Go files for managing RMQ and ROS configurations. | ✅ Active |
| `iface/` | Defines interfaces for publishers, subscribers, and other core components to ensure decoupling. | ✅ Active |
| `logging/` | Custom logging utilities for the agent. | ✅ Active |
| `publishers/` | Implementations for sending data to RabbitMQ (`publishers/rmq`). | ✅ Active |
| `services/` | Contains core business logic and gRPC service implementations. | 🚧 Partial Placeholder |
| `subscribers/` | Implementations for ROS topic subscriptions. | ✅ Active |
| `telemetry/` | Logic for monitoring agent health and performance metrics. | ✅ Active |
| `utils/` | Shared utility functions and state management (e.g., `GeneratorState`). | ✅ Active |

### Key Files
-   **`go.mod`**: Project dependencies and module definition.
-   **`test.go`**: Current main testing entry point demonstrating the generator usage.
-   **`subscribers/ros_subscriber.go`**: Core logic for interfacing with ROS topics.
-   **`publishers/rmq/rabbitmq.go`**: Wrapper for RabbitMQ client operations.

---

## ✨ Key Features

-   **Automated ROS-to-Proto Conversion**: Reduces manual effort in maintaining message definitions across different systems.
-   **Generic Subscriber/Publisher Models**: Uses Go Generics to handle various message types seamlessly.
-   **Resilient RabbitMQ Integration**: Supports exchange/queue declaration and binding management.
-   **Telemetry Support**: Built-in structures for monitoring node health.

---

## 🗺 Dev Roadmap & Placeholders

The following areas are identified for further feature addition:
1.  **`api/`**: Implementation of a control API to manage the agent remotely.
2.  **`services/grpc/task`**: A planned service for executing complex tasks or workflows via gRPC.
3.  **`cmd/`**: Consolidating multiple tools (generator, runner, CLI) into a unified binary.
4.  **`services/grpc/common/status.go`**: Extending agent status reporting.

---

> [!TIP]
> When adding new ROS message types, always run the `code-generator` first to ensure all internal types are synchronized.
