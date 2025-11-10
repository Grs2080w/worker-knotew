# Worker-knotew

Worker to process synchronization jobs (Google Drive / GitHub) using Supabase as a queue/token storage and Redis for caching. Exposes Prometheus metrics for monitoring.

## Overview
- Polls pending jobs in Supabase and processes updates in Google Drive or GitHub.
- Uses Redis for access token caching.
- Exports Prometheus metrics to `/metrics` on port 8080.
- Designed to run in concurrent mode (multiple workers) or isolated mode (one worker per process).

Important files:

- [main.go](/home/grs_s/dev/worker-knotew/main.go) — main entry with multiple workers.

- [packages/worker/main.go](/home/grs_s/dev/worker-knotew/packages/worker/main.go) — concurrent worker loop.

- [packages/worker/IsolatedWorker/root/main.go](/home/grs_s/dev/worker-knotew/packages/worker/IsolatedWorker/root/main.go) — isolated mode (single worker).
- [packages/supa/main.go](/home/grs_s/dev/worker-knotew/packages/supa/main.go) — Supabase abstractions.
- [packages/redis/main.go](/home/grs_s/dev/worker-knotew/packages/redis/main.go) — Redis client and token cache.
- [packages/prometheus/main.go](/home/grs_s/dev/worker-knotew/packages/prometheus/main.go) — recorded metrics.
- [prometheus.yml](/home/grs_s/dev/worker-knotew/prometheus.yml) — Example Prometheus configuration.

- [docker-compose.yml](/home/grs_s/dev/worker-knotew/docker-compose.yml) — Stack for Prometheus/Grafana and worker (expected image).

Important metrics:
- [`prometheus.JobExecutionDuration`](packages/prometheus/main.go)
- [`prometheus.JobsProcessed`](packages/prometheus/main.go)
- [`prometheus.JobsFailed`](packages/prometheus/main.go)
- [`prometheus.JobsError`](packages/prometheus/main.go)
- [`prometheus.IdleIterations`](packages/prometheus/main.go)
- [`prometheus.MemoryUsageBytes`](packages/prometheus/main.go)

## Roadmap

<div align="center"> 
    <img width="70%" src="./roadmap.png" >
</div>

---

## Considerations

This worker is made for Knotew, click the link below when it becomes available!