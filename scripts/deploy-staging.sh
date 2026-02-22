#!/bin/bash
# Deployment Script for E-Commerce API Staging
# This script is executed on the staging server via SSH
# It handles:
# - Docker image pull & verification
# - Graceful container shutdown
# - Health check validation
# - Automated rollback on failure

set -euo pipefail

# =====================
# Configuration
# =====================
REGISTRY="${1:-ghcr.io}"
IMAGE_NAME="${2:-akbarandriansyah22/BackendProject_and_Portofolio/e-commerce-api}"
TAG="${3:-staging}"
ENV_FILE="${4:-/home/deploy/.env.staging}"
CONTAINER_NAME="ecommerce-api"
NETWORK_NAME="ecommerce-network"
PORT="8080"
HEALTH_CHECK_TIMEOUT=30

# =====================
# Logging
# =====================
log_info() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] [INFO] $*"
}

log_error() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] [ERROR] $*" >&2
}

log_success() {
  echo "[$(date +'%Y-%m-%d %H:%M:%S')] [SUCCESS] $*"
}

# =====================
# Validation
# =====================
validate_inputs() {
  if [ ! -f "$ENV_FILE" ]; then
    log_error "Environment file not found: $ENV_FILE"
    exit 1
  fi
  
  if ! command -v docker &> /dev/null; then
    log_error "Docker is not installed"
    exit 1
  fi
  
  log_info "Inputs validated"
}

# =====================
# Docker Network Setup
# =====================
ensure_network() {
  if ! docker network ls | grep -q "$NETWORK_NAME"; then
    log_info "Creating Docker network: $NETWORK_NAME"
    docker network create "$NETWORK_NAME" || true
  fi
}

# =====================
# Image Pull & Verification
# =====================
pull_image() {
  local image="${REGISTRY}/${IMAGE_NAME}:${TAG}"
  log_info "Pulling image: $image"
  
  if ! docker pull "$image"; then
    log_error "Failed to pull image: $image"
    exit 1
  fi
  
  # Verify image exists locally
  if ! docker image inspect "$image" > /dev/null; then
    log_error "Image verification failed"
    exit 1
  fi
  
  log_success "Image pulled successfully"
}

# =====================
# Backup Current Container
# =====================
backup_current() {
  if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    local backup_name="${CONTAINER_NAME}.backup.$(date +%s)"
    log_info "Backing up current container as $backup_name"
    
    docker rename "$CONTAINER_NAME" "$backup_name" || true
    # Keep backup for 1 hour
    docker update --restart=no "$backup_name" || true
  fi
}

# =====================
# Graceful Shutdown
# =====================
graceful_shutdown() {
  if docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    log_info "Gracefully shutting down existing container..."
    
    # Send SIGTERM and wait for graceful shutdown
    docker stop -t 10 "$CONTAINER_NAME" 2>/dev/null || true
    
    # Force remove if still running
    docker rm -f "$CONTAINER_NAME" 2>/dev/null || true
    
    log_info "Container shutdown complete"
  fi
}

# =====================
# Start New Container
# =====================
start_container() {
  local image="${REGISTRY}/${IMAGE_NAME}:${TAG}"
  
  log_info "Starting new container from image: $image"
  
  docker run -d \
    --name "$CONTAINER_NAME" \
    --restart unless-stopped \
    --network "$NETWORK_NAME" \
    --env-file "$ENV_FILE" \
    -p "127.0.0.1:${PORT}:8080" \
    --health-cmd='curl -f http://localhost:8080/health || exit 1' \
    --health-interval=10s \
    --health-timeout=5s \
    --health-retries=3 \
    --log-driver json-file \
    --log-opt max-size=10m \
    --log-opt max-file=3 \
    "$image" || {
    log_error "Failed to start container"
    return 1
  }
  
  log_info "Container started, waiting for health check..."
}

# =====================
# Health Check
# =====================
wait_for_health() {
  local elapsed=0
  local interval=2
  
  while [ $elapsed -lt $HEALTH_CHECK_TIMEOUT ]; do
    local health_status=$(docker inspect --format='{{.State.Health.Status}}' "$CONTAINER_NAME" 2>/dev/null || echo "none")
    
    case "$health_status" in
      "healthy")
        log_success "Container is healthy"
        return 0
        ;;
      "unhealthy")
        log_error "Container marked as unhealthy"
        docker logs "$CONTAINER_NAME" | tail -20
        return 1
        ;;
      *)
        log_info "Health status: $health_status (waiting...)"
        ;;
    esac
    
    sleep $interval
    elapsed=$((elapsed + interval))
  done
  
  log_error "Health check timeout after ${HEALTH_CHECK_TIMEOUT}s"
  docker logs "$CONTAINER_NAME" | tail -20
  return 1
}

# =====================
# Smoke Tests
# =====================
run_smoke_tests() {
  log_info "Running smoke tests..."
  
  local health_response
  health_response=$(curl -s -o /dev/null -w "%{http_code}" "http://127.0.0.1:${PORT}/health" || echo "000")
  
  if [ "$health_response" != "200" ]; then
    log_error "Health endpoint failed: HTTP $health_response"
    return 1
  fi
  
  log_success "Smoke tests passed"
}

# =====================
# Rollback on Failure
# =====================
rollback() {
  log_error "Deployment failed, attempting rollback..."
  
  docker stop -t 5 "$CONTAINER_NAME" 2>/dev/null || true
  docker rm "$CONTAINER_NAME" 2>/dev/null || true
  
  # Restore backup if exists
  local backup=$(docker ps -a --format '{{.Names}}' | grep "^${CONTAINER_NAME}.backup" | sort -r | head -1)
  if [ -n "$backup" ]; then
    log_info "Restoring backup container: $backup"
    docker rename "$backup" "$CONTAINER_NAME"
    docker start "$CONTAINER_NAME"
    
    # Wait for rollback container
    sleep 5
    if docker inspect --format='{{.State.Running}}' "$CONTAINER_NAME" | grep -q true; then
      log_success "Rollback successful"
      return 0
    fi
  fi
  
  log_error "Rollback failed"
  return 1
}

# =====================
# Cleanup
# =====================
cleanup_old_backups() {
  log_info "Cleaning up old backup containers..."
  
  # Remove backups older than 1 hour
  for backup in $(docker ps -a --format '{{.Names}}' | grep "^${CONTAINER_NAME}.backup"); do
    local created=$(docker inspect --format='{{.Created}}' "$backup")
    local created_epoch=$(date -d "$created" +%s)
    local now=$(date +%s)
    local age=$((now - created_epoch))
    
    if [ $age -gt 3600 ]; then
      log_info "Removing old backup: $backup"
      docker rm -f "$backup" || true
    fi
  done
}

# =====================
# Main Execution
# =====================
main() {
  log_info "Starting deployment..."
  log_info "Image: ${REGISTRY}/${IMAGE_NAME}:${TAG}"
  
  validate_inputs
  ensure_network
  pull_image
  backup_current
  graceful_shutdown
  
  if ! start_container; then
    log_error "Failed to start container"
    return 1
  fi
  
  if ! wait_for_health; then
    log_error "Container health check failed"
    if ! rollback; then
      log_error "Rollback also failed - manual intervention required"
      return 1
    fi
    return 1
  fi
  
  if ! run_smoke_tests; then
    log_error "Smoke tests failed"
    if ! rollback; then
      log_error "Rollback also failed"
      return 1
    fi
    return 1
  fi
  
  cleanup_old_backups
  log_success "Deployment completed successfully"
  return 0
}

main "$@"
