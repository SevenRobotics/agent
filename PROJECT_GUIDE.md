# Project Guide: ROS-to-RabbitMQ Go Agent

This repository contains a Go telemetry agent that discovers ROS 1 topics, subscribes to selected topics, converts ROS messages into Protocol Buffer messages, and publishes serialized protobuf payloads to RabbitMQ.

The codebase has two major parts:

- A ROS message generator in `code-generator/` that creates Go ROS message types, `.proto` files, protobuf Go bindings, converters, serializers, and typed pipeline builders.
- A runtime in `telemetry/`, `subscribers/`, and `publishers/` that wires ROS subscriptions to RabbitMQ publishers.

## Requirements

- Go `1.22.5` or newer.
- ROS 1 with a reachable ROS master. The current default is `127.0.0.1:11311`.
- RabbitMQ with a reachable AMQP endpoint.
- `protoc` and `protoc-gen-go` for regenerating protobuf bindings.
- A ROS environment where package paths are discoverable before running the generator. Source the relevant ROS workspace first, for example `source /opt/ros/<distro>/setup.bash` and, when applicable, the workspace `devel/setup.bash`.

## Configuration

Runtime configuration lives in `config/`.

- `config/telemetry_node.yml`
  - `name`: ROS node name used by the conductor when querying the ROS master.
  - `address`: ROS master address, for example `127.0.0.1:11311`.
  - `agent_id`: Robot/agent identifier used to build RabbitMQ exchange, routing key, and queue names.
- `config/rmq_config.yml`
  - `username`, `password`, `host`, and `vhost` for the RabbitMQ AMQP connection.

Do not treat the checked-in YAML as production-secret storage. For real deployments, inject credentials through your deployment layer or replace the local file during provisioning.

## Runtime Flow

The main telemetry runtime is `telemetry/main.go`.

1. It scans `telemetry/genproto/ros/` to build a `utils.GeneratorState` map of generated ROS message support.
2. It calls `converter.AssignBuilder()` from `telemetry/gengo/ros/converter/converter.go`. This registers generated builders under the `ros-rmq` builder key used by the conductor.
3. It loads RabbitMQ and ROS node config from `config/`.
4. It starts a `channel.Conductor`, currently with a hard-coded topic allowlist in `telemetry/main.go`.
5. The conductor waits for ROS master availability, discovers published topics, filters them to user-requested topics that have generated type support, builds pipelines, and starts them.
6. Each pipeline performs:
   - ROS subscription through `subscribers.NewRosSubscriber`.
   - ROS-to-protobuf conversion through a generated converter.
   - Protobuf serialization through a generated serializer.
   - RabbitMQ publishing through `publishers/rmq`.

Run the telemetry runtime with:

```bash
go run ./telemetry
```

For local generator-only experiments, `test.go` and `code-generator/main.go` both run the generator:

```bash
go run ./code-generator
go run ./test.go
```

## Topic Selection

At the moment, runtime topic selection is not config-driven. `telemetry/main.go` contains:

```go
topicList := []string{"/odom_with_amcl"}
```

Update this list when you want the agent to publish additional topics. The topic must be published by ROS and its ROS message type must exist in the generated artifacts. If the topic type is missing from generated support, the conductor will ignore it.

## RabbitMQ Naming

RabbitMQ names are derived in `telemetry/cmd/channel/conductor.go`.

- `agent_id` is normalized by lowercasing, removing `_` and `-`, and converting values like `amr001` or `amr1` to `amr.001`.
- Exchange: `<agent_id>.exchange`
- Routing key: `<agent_id>.<topic>`
- Queue: `<agent_id>.<topic>.q`

ROS topic names are normalized for internal maps by replacing nested `/` separators with `.`, so `/foo/bar` becomes `foo.bar`. The outbound routing key uses that dotted topic form.

Example for `agent_id: amr001` and topic `/odom_with_amcl`:

- Exchange: `amr.001.exchange`
- Routing key: `amr.001.odom_with_amcl`
- Queue: `amr.001.odom_with_amcl.q`

## Generated Artifacts

Generated files are a core part of this repository. Avoid manually editing them unless you are debugging generator output.

- `telemetry/protobuf/ros/<pkg>/<Msg>.proto`: generated proto definitions.
- `telemetry/genproto/ros/<pkg>/<Msg>.pb.go`: generated Go protobuf bindings.
- `telemetry/gengo/ros/<pkg>/msg_*.go`: generated Go ROS message types used by `goroslib`.
- `telemetry/gengo/ros/converter/converter.go`: generated ROS-to-protobuf converters, protobuf serializers, and `GetBuilderFromName`.
- `subscribers/ros/<pkg>/<Msg>.go`: generated subscriber-facing ROS message wrappers.

The generator code lives in `code-generator/rostoproto/`. Its defaults write generated output relative to that package:

- protobuf Go output: `telemetry/genproto/ros/`
- proto definitions: `telemetry/protobuf/ros/`
- goroslib Go message output: `telemetry/gengo/ros/`
- ROS subscriber type output: `subscribers/ros/`
- converter/builder output: `telemetry/gengo/ros/converter/converter.go`

The generator discovers ROS packages via `code-generator/rostoproto/util.FindRosPackages()` and skips packages in the blacklist inside `code-generator/rostoproto/cmd.go`.

## Regenerating ROS and Protobuf Code

Before regenerating, make sure ROS package paths are available in the shell and `protoc` plus `protoc-gen-go` are installed.

```bash
go run ./code-generator
```

After regeneration:

```bash
go test ./...
```

Expected regeneration side effects include changes under `telemetry/protobuf/ros/`, `telemetry/genproto/ros/`, `telemetry/gengo/ros/`, `subscribers/ros/`, and `telemetry/gengo/ros/converter/converter.go`.

## Important Packages

- `telemetry/cmd/channel`
  - `conductor.go`: discovers ROS topics, builds pipeline configs, handles ROS master loss/recovery, and starts pipelines.
  - `pipeline.go`: owns one ROS subscriber, bridge, and RabbitMQ publisher for one topic.
  - `bridge.go`: converts messages from ROS type `S` to protobuf type `P`.
  - `builder.go`: generic builder wrapper used by generated builder registrations.
- `subscribers`
  - `subscriber.go`: generic subscriber interface.
  - `ros_subscriber.go`: `goroslib` subscriber implementation.
- `publishers`
  - `publisher.go`: generic publisher interface.
  - `rmq/rabbitmq.go`: RabbitMQ connection singleton, channel clients, reconnect loop, exchanges, queues, bindings, publish, and consume.
  - `rmq/rmq_publisher.go`: typed protobuf publisher with serializer hook.
  - `rmq/rmq_subscriber.go`: RabbitMQ consumer helper.
- `config`
  - `telemetry_config.go`: shared config structs for ROS, RabbitMQ, and ROS-RMQ pipeline setup.
- `utils`
  - `generator_state.go`: generated-message state and global builder registry.
- `services`
  - Currently skeletal gRPC/task scaffolding.
- `cmd`, `telemetry/cmd/init`, top-level `Makefile`, `telemetry/Makefile`, and `services/Makefile`
  - Currently placeholders or empty.

## Pipeline Details

`channel.NewPipeline` creates unbuffered `in`, `out`, `done`, and error channels. It initializes:

- `subscribers.NewRosSubscriber[S]` for the ROS topic.
- `rmq.NewRabbitMQ`, which is a singleton connection per process.
- One RabbitMQ client/channel per pipeline name.
- `rmq.NewRMQPublisher[P]`, which declares the topic exchange, queue, and binding.
- A `Bridge[S, P]` with generated converter and serializer functions.

`pipeline.Start` launches the bridge, initializes the ROS subscriber, then launches the publisher. `pipeline.Shutdown` signals `done` and closes the ROS subscriber/node.

Current channel behavior to keep in mind:

- Pipeline channels are unbuffered, so slow RabbitMQ publishing can apply backpressure to conversion and ROS callback forwarding.
- `Bridge` and ROS subscriber callbacks recover from sends to closed channels, but lifecycle changes should still be made carefully.
- Publisher send errors are throttled to avoid log flooding during RabbitMQ outages.
- RabbitMQ reconnect recreates client channels, but publisher declarations and bindings are performed during initial configure. Reconnection behavior should be tested when changing RMQ lifecycle code.

## Development Commands

Common commands:

```bash
go mod download
go test ./...
go run ./code-generator
go run ./telemetry
```

Build packages without running the long-lived telemetry process:

```bash
go test ./... -run '^$'
```

Format Go changes:

```bash
gofmt -w <files>
```

## Adding Support for a New ROS Topic

1. Ensure the ROS package that defines the topic message is discoverable in the shell.
2. Run `go run ./code-generator`.
3. Confirm generated support exists in `telemetry/genproto/ros/`, `telemetry/gengo/ros/`, and `telemetry/gengo/ros/converter/converter.go`.
4. Add the ROS topic path to `topicList` in `telemetry/main.go`.
5. Run `go test ./...`.
6. Run `go run ./telemetry` with ROS master and RabbitMQ available.
7. Confirm the expected RabbitMQ exchange, routing key, and queue are created.

## Coding Guidelines

- Prefer existing generic interfaces (`Subscriber`, `Publisher`, `Pipeline`, `Builder`) before adding new abstractions.
- Keep handwritten runtime changes out of generated trees unless you are changing the generator.
- When changing message conversion or serialization behavior, update the generator templates/functions in `code-generator/rostoproto/`, then regenerate.
- Keep topic and routing-name behavior consistent with `conductor.go`; downstream consumers depend on those names.
- Treat `telemetry/main.go` as the current runtime entrypoint, not the placeholder packages under `cmd/`.
- Be careful with shared RabbitMQ singleton state in tests or multi-pipeline changes.

## Known Gaps

- Topic allowlisting is hard-coded in `telemetry/main.go`.
- There is no committed automated integration test for ROS master discovery or RabbitMQ publish/consume behavior.
- Several app/service entrypoints and Makefiles are placeholders.
- The checked-in config format has no environment-variable override layer yet.
- Generated files are large, so reviews should distinguish generator changes from regeneration churn.
