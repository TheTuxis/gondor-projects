# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview
Projects microservice for Gondor platform. Go/Gin service handling project management: projects, tasks, phases, team members, and deliverables.

## Commands
- `make build` -- compile to bin/server
- `make run` -- run locally (needs PostgreSQL + Redis)
- `make test` -- run all tests with race detector
- `make lint` -- golangci-lint
- `make docker` -- build Docker image
- `make migrate-up` -- run database migrations
- `make migrate-down` -- rollback migrations

## Architecture
- `cmd/server/main.go` -- entry point, dependency injection, route registration
- `internal/config/` -- env-based configuration
- `internal/model/` -- GORM domain models (Project, Task, Phase, ProjectMember, Deliverable)
- `internal/repository/` -- database access layer
- `internal/service/` -- business logic
- `internal/handler/` -- HTTP handlers (Gin)
- `internal/middleware/` -- JWT auth (validate-only), logging
- `internal/pkg/jwt/` -- JWT validation (tokens issued by gondor-users-security)

## Key Decisions
- JWT tokens are validated only (issued by gondor-users-security service)
- Port 8002 (gondor-users-security is 8001)
- Database: gondor_projects (PostgreSQL, database-per-service)
- Soft delete for projects and tasks (GORM DeletedAt)
- Multi-tenancy via company_id on projects
- Project statuses: active, completed, cancelled, on_hold
- Task statuses: pending, in_progress, completed, cancelled
- Task priorities: low, medium, high, critical
- Member roles: pm, team_member, stakeholder, client
- Tasks support self-referential parent_id for subtasks
- All routes under /v1/projects/ prefix
- /health and /metrics skip JWT auth

## Database
PostgreSQL with GORM. Tables: projects, tasks, phases, project_members, deliverables.

## Environment Variables
- `PORT` (default: 8002)
- `DATABASE_URL` (default: postgres://gondor:gondor_dev@localhost:5432/gondor_projects?sslmode=disable)
- `JWT_SECRET` (default: change-me-in-production)
- `REDIS_URL` (default: localhost:6379)
- `NATS_URL` (default: nats://localhost:4222)
- `ENVIRONMENT` (default: development)
