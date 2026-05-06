================================================================================
ROS-TO-RABBITMQ GO TELEMETRY AGENT - ARCHITECTURE OVERVIEW
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│                          SYSTEM COMPONENTS                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌──────────────────┐         ┌──────────────────┐         ┌──────────────────┐
│   ROS MASTER     │         │   GO TELEMETRY   │         │    RABBITMQ      │
│  (127.0.0.1:     │◄───────►│      AGENT       │────────►│   AMQP BROKER    │
│    11311)        │         │                  │         │                  │
└──────────────────┘         └──────────────────┘         └──────────────────┘
        │                             │                             │
        │ Publishes topics            │ Subscribes & converts       │ Receives
        │ /odom_with_amcl, etc.       │ ROS → Protobuf             │ protobuf
        │                             │                             │ payloads


================================================================================
CODE GENERATION PHASE (Offline - code-generator/)
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│  INPUT: ROS Package Definitions (discovered via ROS_PACKAGE_PATH)           │
└─────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
        ┌───────────────────────────────────────────────┐
        │   code-generator/rostoproto/                  │
        │   - FindRosPackages()                         │
        │   - Template-based code generation            │
        └───────────────────────────────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
                    ▼                               ▼
        ┌─────────────────────┐         ┌─────────────────────┐
        │  .proto files       │         │  Go ROS Messages    │
        │  telemetry/         │         │  telemetry/gengo/   │
        │  protobuf/ros/      │         │  ros/<pkg>/msg_*.go │
        └─────────────────────┘         └─────────────────────┘
                    │                               │
                    │ protoc                        │
                    ▼                               │
        ┌─────────────────────┐                    │
        │  Protobuf Go        │                    │
        │  telemetry/         │                    │
        │  genproto/ros/      │                    │
        │  <pkg>/<Msg>.pb.go  │                    │
        └─────────────────────┘                    │
                    │                               │
                    └───────────────┬───────────────┘
                                    ▼
        ┌───────────────────────────────────────────────┐
        │  Generated Artifacts:                         │
        │  - Converters (ROS → Protobuf)                │
        │  - Serializers (Protobuf → bytes)            │
        │  - Pipeline Builders                          │
        │  - Subscriber Wrappers                        │
        │                                               │
        │  Output: telemetry/gengo/ros/converter/       │
        │          converter.go                         │
        │          subscribers/ros/<pkg>/<Msg>.go       │
        └───────────────────────────────────────────────┘


================================================================================
RUNTIME ARCHITECTURE (telemetry/main.go)
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│                        INITIALIZATION SEQUENCE                               │
└─────────────────────────────────────────────────────────────────────────────┘

1. Load Configuration
   ├── config/telemetry_node.yml  (ROS master, node name, agent_id)
   └── config/rmq_config.yml      (RabbitMQ credentials)

2. Scan Generated Artifacts
   └── Build utils.GeneratorState map from telemetry/genproto/ros/

3. Register Builders
   └── converter.AssignBuilder() registers all generated pipeline builders

4. Start Conductor
   └── Hard-coded topic allowlist: ["/odom_with_amcl"]


┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONDUCTOR (channel.Conductor)                             │
│                   telemetry/cmd/channel/conductor.go                         │
└─────────────────────────────────────────────────────────────────────────────┘

                    ┌─────────────────────┐
                    │  Wait for ROS       │
                    │  Master Available   │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │  Discover Published │
                    │  ROS Topics         │
                    │  (Poll every 1 sec) │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │  Filter Topics      │
                    │  - In allowlist?    │
                    │  - Has generated    │
                    │    type support?    │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │  Build Pipelines    │
                    │  (one per topic)    │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │  Start All          │
                    │  Pipelines          │
                    └─────────────────────┘

    RESILIENCE: On ROS Master Failure
    └──► Shutdown pipelines → Wait for master → Rebuild → Restart


================================================================================
PIPELINE ARCHITECTURE (channel.Pipeline)
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│                     PER-TOPIC PIPELINE INSTANCE                              │
│                   telemetry/cmd/channel/pipeline.go                          │
└─────────────────────────────────────────────────────────────────────────────┘

    ROS Topic: /odom_with_amcl
    ROS Type:  nav_msgs/Odometry
    Proto Type: nav_msgs.Odometry

┌──────────────┐      ┌──────────────┐      ┌──────────────┐      ┌──────────┐
│ ROS          │      │              │      │              │      │ RabbitMQ │
│ Subscriber   │─────►│    Bridge    │─────►│ RMQ          │─────►│ Broker   │
│              │ chan │              │ chan │ Publisher    │ AMQP │          │
│ (goroslib)   │  in  │  Converter   │ out  │              │      │          │
└──────────────┘      │  Serializer  │      └──────────────┘      └──────────┘
                      └──────────────┘

CHANNELS (unbuffered):
├── in:   ROS messages (type S)
├── out:  Serialized protobuf ([]byte)
├── done: Shutdown signal
└── err:  Error reporting

COMPONENTS:

┌─────────────────────────────────────────────────────────────────────────────┐
│ 1. ROS Subscriber (subscribers/ros_subscriber.go)                           │
├─────────────────────────────────────────────────────────────────────────────┤
│ - Uses goroslib to subscribe to ROS topic                                   │
│ - Receives ROS messages via callbacks                                       │
│ - Forwards to 'in' channel                                                  │
│ - Handles callback→closed channel gracefully                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 2. Bridge (channel.Bridge)                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ - Reads from 'in' channel                                                   │
│ - Converts: ROS type S → Protobuf type P (generated converter)             │
│ - Serializes: Protobuf P → []byte (generated serializer)                   │
│ - Writes to 'out' channel                                                   │
│ - Recovers from closed channel during shutdown                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ 3. RMQ Publisher (publishers/rmq/rmq_publisher.go)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│ - Reads from 'out' channel                                                  │
│ - Publishes []byte to RabbitMQ                                              │
│ - Declares exchange, queue, binding on init                                 │
│ - Uses shared RMQ connection singleton                                      │
│ - Throttles error logging (1 per 5 seconds)                                 │
└─────────────────────────────────────────────────────────────────────────────┘


================================================================================
RABBITMQ INFRASTRUCTURE (publishers/rmq/)
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│                    CONNECTION SINGLETON                                      │
│                   (rmq.NewRabbitMQ)                                         │
└─────────────────────────────────────────────────────────────────────────────┘

    One Shared Connection Per Process
    ├── TCP Dial Timeout: 5 seconds
    ├── AMQP Heartbeat: 10 seconds
    └── Per-Pipeline Channel Clients

    RECONNECTION STRATEGY:
    ├── Monitor connection & channel closure
    ├── Exponential backoff: 1s → 30s cap
    ├── Recreate all client channels after reconnect
    └── Persistent messages (no offline buffer)


┌─────────────────────────────────────────────────────────────────────────────┐
│                         NAMING CONVENTION                                    │
└─────────────────────────────────────────────────────────────────────────────┘

    agent_id: "amr001" → normalized: "amr.001"
    topic:    "/odom_with_amcl" → normalized: "odom_with_amcl"

    Exchange:    <agent_id>.exchange          → "amr.001.exchange"
    Routing Key: <agent_id>.<topic>           → "amr.001.odom_with_amcl"
    Queue:       <agent_id>.<topic>.q         → "amr.001.odom_with_amcl.q"


================================================================================
DATA FLOW (End-to-End Message Journey)
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 1: ROS Message Published                                              │
└─────────────────────────────────────────────────────────────────────────────┘
    ROS Node publishes nav_msgs/Odometry to /odom_with_amcl
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 2: goroslib Subscriber Callback                                       │
└─────────────────────────────────────────────────────────────────────────────┘
    RosSubscriber[nav_msgs.Odometry] receives message
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 3: Forward to Pipeline                                                │
└─────────────────────────────────────────────────────────────────────────────┘
    Send nav_msgs.Odometry → 'in' channel (unbuffered)
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 4: Bridge Conversion                                                  │
└─────────────────────────────────────────────────────────────────────────────┘
    Generated Converter: nav_msgs.Odometry (ROS) 
                         → nav_msgs.Odometry (Protobuf)
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 5: Protobuf Serialization                                             │
└─────────────────────────────────────────────────────────────────────────────┘
    Generated Serializer: nav_msgs.Odometry (Protobuf) → []byte
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 6: Forward Serialized Message                                         │
└─────────────────────────────────────────────────────────────────────────────┘
    Send []byte → 'out' channel (unbuffered)
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  STEP 7: RabbitMQ Publish                                                   │
└─────────────────────────────────────────────────────────────────────────────┘
    RMQPublisher publishes to:
    - Exchange: "amr.001.exchange"
    - Routing Key: "amr.001.odom_with_amcl"
    - Queue: "amr.001.odom_with_amcl.q"
    - Message: Persistent, Content-Type: application/x-protobuf


================================================================================
RESILIENCE & ERROR HANDLING
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│  ROS MASTER FAILURE RECOVERY                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│  1. Conductor polls ROS master every 1 second                               │
│  2. On scan failure:                                                        │
│     └──► Shutdown all pipelines                                             │
│     └──► Wait for ROS master reachable                                      │
│     └──► Reset state (topics, builders, pipelines)                          │
│     └──► Rebuild pipelines                                                  │
│     └──► Restart pipelines                                                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│  RABBITMQ CONNECTION FAILURE RECOVERY                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  1. Monitor connection & channel closure events                             │
│  2. On connection loss:                                                     │
│     └──► Reconnect with exponential backoff (1s → 30s)                      │
│     └──► Recreate all client channels                                       │
│  3. On individual channel loss:                                             │
│     └──► Recreate channel if connection alive                               │
│  4. Messages during outage: LOST (no offline buffer)                        │
│  5. Error logging: Throttled to 1 per 5 seconds                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│  CHANNEL BACKPRESSURE                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│  - Unbuffered channels create natural backpressure                          │
│  - Slow RMQ publish blocks Bridge                                           │
│  - Bridge blocks ROS subscriber callback                                    │
│  - Consider buffering for high-frequency topics                             │
└─────────────────────────────────────────────────────────────────────────────┘


================================================================================
DIRECTORY STRUCTURE
================================================================================

ros-to-rabbitmq-agent/
├── code-generator/
│   ├── rostoproto/              # Core generator logic
│   │   ├── cmd.go              # Package discovery, blacklist
│   │   └── util.go             # FindRosPackages()
│   └── main.go                 # Generator entrypoint
│
├── telemetry/
│   ├── main.go                 # Runtime entrypoint (TOPIC ALLOWLIST HERE)
│   ├── cmd/
│   │   └── channel/
│   │       ├── conductor.go    # Topic discovery & lifecycle
│   │       ├── pipeline.go     # Per-topic pipeline
│   │       ├── bridge.go       # Conversion logic
│   │       └── builder.go      # Builder wrapper
│   ├── protobuf/ros/           # Generated .proto files
│   ├── genproto/ros/           # Generated .pb.go files
│   └── gengo/ros/
│       ├── <pkg>/msg_*.go      # Generated ROS messages
│       └── converter/
│           └── converter.go    # Generated converters & builders
│
├── subscribers/
│   ├── subscriber.go           # Generic interface
│   ├── ros_subscriber.go       # goroslib implementation
│   └── ros/<pkg>/<Msg>.go      # Generated wrappers
│
├── publishers/
│   ├── publisher.go            # Generic interface
│   └── rmq/
│       ├── rabbitmq.go         # Connection singleton
│       ├── rmq_publisher.go    # Typed publisher
│       └── rmq_subscriber.go   # Consumer helper
│
├── config/
│   ├── telemetry_node.yml      # ROS master, agent_id
│   └── rmq_config.yml          # RabbitMQ credentials
│
└── utils/
    └── generator_state.go      # Builder registry


================================================================================
KEY WORKFLOWS
================================================================================

┌─────────────────────────────────────────────────────────────────────────────┐
│  ADDING A NEW ROS TOPIC                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│  1. Source ROS workspace: source /opt/ros/<distro>/setup.bash              │
│  2. Run generator: go run ./code-generator                                  │
│  3. Verify generated artifacts in:                                          │
│     - telemetry/genproto/ros/<pkg>/                                         │
│     - telemetry/gengo/ros/<pkg>/                                            │
│     - telemetry/gengo/ros/converter/converter.go                            │
│  4. Add topic to allowlist in telemetry/main.go:                            │
│     topicList := []string{"/odom_with_amcl", "/new_topic"}                  │
│  5. Test: go test ./...                                                     │
│  6. Run: go run ./telemetry                                                 │
│  7. Verify RabbitMQ exchange/queue creation                                 │
└─────────────────────────────────────────────────────────────────────────────┘


================================================================================
CURRENT LIMITATIONS & GAPS
================================================================================

❌ Hard-coded topic allowlist (in telemetry/main.go)
❌ No config-driven topic selection
❌ No automated integration tests for recovery scenarios
❌ No offline message buffering during RMQ outages
❌ No environment variable override for config
❌ Placeholder cmd/ and Makefile structure
❌ No end-to-end delivery guarantees
❌ Review burden for large generated file diffs


================================================================================
DEVELOPMENT COMMANDS
================================================================================

# Download dependencies
go mod download

# Run all tests
go test ./...

# Regenerate code (requires ROS env sourced)
go run ./code-generator

# Run telemetry agent
go run ./telemetry

# Build without running
go test ./... -run '^$'

# Format code
gofmt -w <files>