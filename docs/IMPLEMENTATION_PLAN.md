# distributed-job-queue — Unified Roadmap (Codex-Oriented)

This document describes the intended direction and evolution of the `distributed-job-queue` project.
It is written for AI-assisted development (Codex) to keep implementation decisions consistent
with the architectural intent.

The roadmap prioritizes:
1) correctness and reliability first,
2) measurable performance improvements second,
3) distributed-system readiness third.

---

## Guiding Principles

- **Correctness before optimization**: build a reliable queue with durable state and safe concurrency.
- **Clear separation of responsibilities**: API handles requests; workers execute jobs; database is the source of truth.
- **Failure is normal**: timeouts, retries, and idempotency are fundamental.
- **Bounded concurrency**: avoid unbounded goroutines; apply backpressure.
- **Evolve by replacing infrastructure layers, not rewriting the core model**.

---

## Phase 1 — Foundation (Pure Go + PostgreSQL)

### Goal
Build a production-oriented asynchronous job queue service using **Golang** for both API and workers,
with **PostgreSQL** as the system of record. Focus on correctness, concurrency safety, and clear state transitions.

### Architecture
- **API Service (Go + Gin)**
  - Receives HTTP requests
  - Validates input
  - Creates jobs in PostgreSQL in `queued` state
  - Returns `job_id` immediately (non-blocking request path)
- **Worker Service (Go)**
  - Polls/claims runnable jobs from PostgreSQL
  - Executes jobs using a bounded worker pool
  - Updates job state and results in PostgreSQL
- **PostgreSQL**
  - Stores job payload, state, attempts, timestamps, and results
  - Enforces correctness via constraints and transactions

### Key Deliverables
#### 1) Database Schema & Migrations
- `jobs` table with:
  - `id` (primary key)
  - `type` (job type)
  - `payload` (JSON)
  - `status` (queued/running/succeeded/failed)
  - `attempts`, `max_attempts`
  - `run_at` / `next_run_at` (optional scheduling)
  - `locked_at`, `locked_by` (optional worker claiming fields)
  - `last_error`, `result`
  - `created_at`, `updated_at`
- Constraints:
  - idempotency via unique key (optional in v1 but recommended)
- Indexes:
  - support "find next runnable jobs" efficiently
- Migrations managed via a Go migration tool (goose / golang-migrate)

#### 2) Concurrency Model (Worker Pool)
- Implement bounded concurrency using:
  - a worker pool + job channel, or
  - semaphore pattern
- No unbounded goroutine spawning.
- Support backpressure (do not claim more jobs than can be processed).

#### 3) Safe Job Claiming & State Transitions
- Only one worker should process a given job at a time.
- Use transactional state transitions.
- Prefer row-level locking patterns (implementation detail may evolve):
  - claim job in a transaction
  - update state `queued -> running` atomically

#### 4) Timeouts, Retries, Backoff
- Per-job execution timeout via `context`.
- Retry failed jobs with exponential backoff.
- Persist attempts and error messages in PostgreSQL.
- Define a final failure state once retry limits are exceeded.

#### 5) Graceful Shutdown
- API: stop accepting new requests, shutdown HTTP server with timeout.
- Worker: stop claiming new jobs, finish in-progress tasks or respect timeouts, persist state safely.

#### 6) Local Deployment
- `docker-compose` orchestrates:
  - API service
  - Worker service
  - PostgreSQL
- Developer workflow:
  - `make up`, `make down`, `make test`, `make lint`

### Phase 1 “Done” Criteria
- Minimal set of endpoints to create jobs and query job status.
- Jobs are processed concurrently with bounded concurrency.
- Job state is durable and correct across restarts.
- Retries and timeouts behave predictably.
- Race conditions are addressed at the database boundary (not via shared memory).

---

## Phase 2 — Optimization (Optional Rust Compute Engine via gRPC)

### Goal
Introduce a **Rust Compute Engine** for CPU-bound tasks only, once a real bottleneck is identified.
The Go worker transitions from "executor" to "dispatcher/orchestrator" for specific job types.

**Important scope control**: Phase 2 should start with exactly one CPU-heavy job type and
include a small benchmark to justify the added complexity.

### Architecture Shift
- **API Service (Go + Gin)**: unchanged (HTTP, validation, job creation).
- **Worker Service (Go)**:
  - continues to claim jobs from PostgreSQL
  - for CPU-heavy job types, sends payload to Rust via gRPC
  - handles gRPC timeouts and failures
  - persists results and job state in PostgreSQL
- **Compute Engine (Rust)**:
  - stateless gRPC server
  - performs CPU-bound computation
  - returns result deterministically

### Key Deliverables
#### 1) Protobuf Contract
- Define `.proto` under `/proto`
- Include versioning (e.g., `compute.v1`)
- Strict request/response schema:
  - `job_type`, `payload`, optional `trace_id`
  - result data and error codes

#### 2) Rust gRPC Server (Tonic + Tokio)
- Implement gRPC server in Rust:
  - concurrency handled by Tokio runtime
  - multi-threaded execution as needed
- Compute engine should be:
  - **stateless**
  - **side-effect-free** (pure function) to make retries safe

#### 3) Go gRPC Client Integration
- Connection management and timeouts.
- Failure handling policy:
  - if gRPC fails, job remains retryable
  - avoid “double retry storms” by clearly defining where retries happen

#### 4) Minimal Benchmark
- Measure at least one of:
  - throughput improvement
  - p95/p99 latency stability
- Benchmark should be reproducible locally.

### Phase 2 “Done” Criteria
- One job type is routed to Rust compute engine.
- System remains correct under failure/retry scenarios.
- Clear justification for Rust exists (via benchmark or observed behavior).

---

## Phase 3 — Distributed Readiness & Deployment (Kubernetes-Friendly)

### Goal
Make the system deployable and scalable in a multi-instance environment without changing core logic.
This phase focuses on operational readiness, scaling, and observability.

### Deliverables
- Container images for API and worker.
- Deployment manifests (or Helm chart) for:
  - API service
  - Worker service
  - PostgreSQL (dev-only; production would use managed DB)
- Horizontal scaling:
  - multiple API instances
  - multiple worker instances
- Observability baseline:
  - structured logging with correlation IDs
  - basic metrics (job counts, durations, failures) if feasible

### Scaling Strategy (Practical First)
- Use HPA/KEDA (preferred) based on:
  - worker CPU usage (basic)
  - or custom metrics (queue depth) if available

### Optional Stretch Goal
- A Go-based Kubernetes operator that scales workers/compute-engine based on queue depth.
  - This is optional and should not block core phases.

### Phase 3 “Done” Criteria
- System can run with multiple API/worker replicas.
- No duplicate job processing under scale (correct claiming).
- Deployed services can be restarted/rolled without corrupting state.

---

## Monorepo Directory Structure (Unified)

```text
/distributed-job-queue
├── /api-service (Go)          # Gin API service
├── /worker-service (Go)       # Job worker / orchestrator
├── /compute-engine (Rust)     # Optional Phase 2 component
├── /proto                     # Shared gRPC definitions
├── /migrations                # SQL migrations
├── docker-compose.yml
└── docs/
    └── PROJECT_OVERVIEW.md
```

## Notes for Codex (Implementation Intent)

* Prefer simple, explicit code over clever abstractions.

* Keep Phase 1 implementation fully functional without Phase 2.

* Phase 2 must not pollute Phase 1 core logic:
  * route only specific job types to Rust
  * keep compute engine stateless/pure

* Database correctness is the primary mechanism for coordinating concurrency:

* do not rely on in-memory locks across processes

* Always preserve deterministic behavior under retries and restarts.