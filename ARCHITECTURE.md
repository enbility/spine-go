# SPINE-Go Architecture Documentation

This document explains the overall architecture of the SPINE-Go library, which provides an implementation of the EEBUS SPINE 1.3 specification in Go.

## Overview

SPINE-Go is a library that implements the SPINE (Smart Premises Interoperable Neutral-message Exchange) protocol, which is part of the EEBUS specification for smart energy management systems. The library provides a complete stack for creating and managing SPINE devices that can communicate with each other over networks.

## Core Components

The architecture is organized into several key packages, each serving a specific purpose:

### 1. API Package (`api/`)

Contains the core interfaces that define contracts for all major components. This package serves as the foundation for the entire architecture and enables dependency injection and testing through mocking.

**Key Interfaces:**
- `DeviceInterface` - Common device functionality
- `DeviceLocalInterface` - Local device management
- `DeviceRemoteInterface` - Remote device representation  
- `EntityInterface` - Entity management
- `FeatureInterface` - Feature functionality
- `SenderInterface` - Message sending capabilities
- `EventHandlerInterface` - Event handling

### 2. Model Package (`model/`)

Contains the Go representation of the SPINE data model with proper JSON serialization support and EEBUS tags for generic feature-to-function mapping.

**Key Components:**
- Data structures for all SPINE message types
- Command and function definitions
- Network management types
- Address and identification structures

### 3. SPINE Package (`spine/`)

The main implementation package containing the core business logic for SPINE devices, entities, features, functions, and data management.

## Hierarchical Architecture

The SPINE architecture follows a hierarchical structure:

```
Device
├── Entity (0..n)
    ├── Feature (0..n)
        ├── Function (0..n)
            └── Data
```

### Device Level

**DeviceLocal** (`device_local.go`)
- Represents the local SPINE device
- Manages local entities and remote device connections
- Handles incoming SPINE messages and routing
- Provides node management, subscription management, and binding management
- Contains device information like brand, model, serial number

**DeviceRemote** (`device_remote.go`)
- Represents a remote SPINE device
- Manages remote entities discovered through detailed discovery
- Handles incoming messages from the associated device
- Maintains connection to the local device

### Entity Level

**EntityLocal** (`entity_local.go`)
- Represents a local entity within the device
- Manages local features
- Handles heartbeat management for non-device-information entities
- Examples: EV charging point, heat pump, battery storage

**EntityRemote** (`entity_remote.go`)
- Represents a remote entity from a connected device
- Manages remote features discovered through detailed discovery
- Maintains reference to parent remote device

### Feature Level

**FeatureLocal** (`feature_local.go`)
- Implements specific SPINE feature functionality
- Manages function data and operations
- Handles subscriptions and bindings to remote features
- Processes incoming messages for the feature
- Examples: Device Diagnosis, Load Control, Measurement

**FeatureRemote** (`feature_remote.go`)
- Represents a remote feature
- Stores data received from remote devices
- Maintains operation capabilities and response delays

## Communication Architecture

### Message Flow

1. **Outgoing Messages:**
   ```
   Local Feature → Sender → SHIP Connection → Network
   ```

2. **Incoming Messages:**
   ```
   Network → SHIP Connection → Device.ProcessCmd() → Feature.HandleMessage()
   ```

### Key Communication Components

**Sender** (`send.go`)
- Handles all outgoing SPINE messages
- Manages message counters and caching
- Supports different message types: Request, Reply, Notify, Write
- Implements message deduplication for certain message types

**Message Processing**
- `DeviceLocal.ProcessCmd()` - Main entry point for incoming messages
- Routes messages to appropriate local features
- Handles acknowledgments and error responses
- Manages message validation and addressing

## Management Systems

### Node Management (`nodemanagement.go`)

The Node Management feature is present on every device (Entity 0, Feature 0) and handles:

- **Detailed Discovery**: Exchange of device, entity, and feature information
- **Destination Lists**: Available communication endpoints
- **Use Case Data**: Supported use cases and their configurations
- **Subscription Management**: Managing data subscriptions between features
- **Binding Management**: Managing persistent connections between features

### Subscription Manager (`subscription_manager.go`)

Manages subscriptions between client and server features:
- Tracks active subscriptions between local and remote features
- Handles subscription requests and deletions
- Automatically notifies subscribed features when data changes
- Manages subscription data persistence

### Binding Manager (`binding_manager.go`)

Manages bindings (persistent connections) between features:
- Tracks active bindings between client and server features
- Handles binding requests and deletions
- Enforces binding constraints (e.g., one remote binding per local server feature)
- Manages binding data persistence

### Heartbeat Manager (`heartbeat_manager.go`)

Manages heartbeat functionality for device diagnosis:
- Sends periodic heartbeat messages to subscribed remote features
- Configurable heartbeat intervals
- Automatically starts/stops based on subscription state
- Used for connection monitoring and fault detection

## Event System (`events.go`)

The library includes a comprehensive event system for notifying applications about important state changes:

**Event Types:**
- `EventTypeDeviceChange` - Device connection/disconnection
- `EventTypeEntityChange` - Entity addition/removal
- `EventTypeSubscriptionChange` - Subscription state changes
- `EventTypeBindingChange` - Binding state changes  
- `EventTypeDataChange` - Feature data updates

**Event Handling Levels:**
- `EventHandlerLevelCore` - For internal library components (synchronous)
- `EventHandlerLevelApplication` - For application code (asynchronous)

## Data Flow Examples

### Device Discovery Process

1. Local device connects to remote device via SHIP
2. Local device requests detailed discovery from remote device
3. Remote device responds with device, entity, and feature information
4. Local device creates corresponding remote device, entity, and feature objects
5. Events are published for each discovered component
6. Applications can then establish subscriptions and bindings

### Data Subscription Process

1. Local client feature requests subscription to remote server feature
2. Subscription manager validates and stores the subscription
3. Remote device acknowledges the subscription
4. When remote feature data changes, notifications are sent to subscriber
5. Local feature updates its cached data and publishes events

### Feature Communication

1. Local feature wants to read data from remote feature
2. Sender creates and sends a read request message
3. Remote feature receives request and sends reply with data
4. Local feature processes reply, updates cache, and publishes events
5. Applications receive events and can access the updated data

## Integration Points

### SHIP Integration

SPINE-Go integrates with the SHIP (Smart Home IP) protocol layer:
- Uses SHIP for device discovery and connection establishment
- SHIP handles network-level communication and security
- SPINE messages are transported as SHIP payload

### Application Integration

Applications integrate with SPINE-Go through:
- Event subscriptions for state change notifications
- Direct API calls for data access and control
- Feature-specific helper libraries for use case implementations

## Thread Safety

The library is designed to be thread-safe:
- Uses mutexes to protect shared data structures
- Event publishing is handled safely across goroutines
- Message processing includes proper synchronization
- Managers use appropriate locking for concurrent access

## Testing Architecture

The architecture supports comprehensive testing through:
- **Mocks Package** (`mocks/`) - Auto-generated mocks for all interfaces
- **Integration Tests** (`integration_tests/`) - End-to-end testing scenarios
- **Unit Tests** - Extensive unit test coverage for individual components
- **Test Helpers** - Common testing utilities and data generators

This modular and interface-driven architecture provides flexibility, testability, and clear separation of concerns while maintaining compliance with the EEBUS SPINE specification.
