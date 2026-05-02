#!/bin/bash
# Helper script to run Pulumi commands in Docker
# Usage: ./scripts/pulumi-docker.sh [preview|up|down|destroy|login|logout|stack|config] [args...]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Default values
ENV=${PULUMI_ENV:-local}
STACK=${PULUMI_STACK:-default}
SERVICE="subject-data"

# Determine backend URL
if [ "$ENV" = "local" ]; then
    BACKEND_URL="file:///root/.pulumi"
    COMPOSE_FILE="docker-compose.yml"
    SERVICE_NAME="pulumi-local"
else
    BACKEND_URL="s3://pulumi-state-${SERVICE}-${ENV}"
    COMPOSE_FILE="docker-compose.yml"
    SERVICE_NAME="pulumi"
fi

# Export for docker-compose
export PULUMI_BACKEND_URL="$BACKEND_URL"
export PULUMI_ENV="$ENV"
export PULUMI_STACK="$STACK"

COMMAND="${1:-help}"
shift || true

case "$COMMAND" in
    login)
        log_info "Logging into Pulumi backend: $BACKEND_URL"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi login "$BACKEND_URL"
        ;;

    logout)
        log_info "Logging out from Pulumi backend"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi logout
        ;;

    stack)
        SUBCMD="${1:-list}"
        shift || true
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi stack "$SUBCMD" "$@"
        ;;

    init)
        STACK_NAME="${1:-$STACK}"
        log_info "Initializing stack: $STACK_NAME"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi stack init "$STACK_NAME"
        ;;

    config)
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi config "$@"
        ;;

    preview)
        log_info "Running pulumi preview for stack: $STACK"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi preview --stack "$STACK" "$@"
        ;;

    up)
        log_info "Running pulumi up for stack: $STACK"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi up --stack "$STACK" --yes "$@"
        ;;

    down)
        log_info "Running pulumi destroy for stack: $STACK"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi destroy --stack "$STACK" --yes "$@"
        ;;

    destroy)
        log_info "Running pulumi destroy for stack: $STACK"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi destroy --stack "$STACK" --yes "$@"
        ;;

    refresh)
        log_info "Running pulumi refresh for stack: $STACK"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi refresh --stack "$STACK" --yes "$@"
        ;;

    output)
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi stack output --stack "$STACK" "$@"
        ;;

    logs)
        docker compose -f "$COMPOSE_FILE" logs --follow "$SERVICE_NAME"
        ;;

    build)
        log_info "Building pulumi provider plugins..."
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" \
            pulumi plugin install resource aws v6.83.2
        ;;

    shell)
        log_info "Starting interactive shell in pulumi container"
        docker compose -f "$COMPOSE_FILE" run --rm "$SERVICE_NAME" bash
        ;;

    run-local)
        log_info "Running pulumi preview with local backend (no AWS credentials needed)"
        docker compose -f docker-compose.yml run --rm pulumi-local \
            pulumi preview --stack "$STACK" "$@"
        ;;

    help|--help|-h)
        cat << EOF
Pulumi Docker Helper Script

Usage: $0 [command] [args...]

Commands:
  login              Login to Pulumi backend
  logout             Logout from Pulumi backend
  stack [cmd]        Manage stacks (init, list, select, rm)
  config             Manage stack configuration
  preview            Run pulumi preview (dry-run)
  up                 Run pulumi up (apply changes)
  down/destroy       Run pulumi destroy (remove all resources)
  refresh            Refresh pulumi state
  output             Show stack outputs
  logs               View container logs
  build              Build/install pulumi plugins
  shell              Start interactive shell in container
  run-local          Run pulumi preview with local backend (no AWS)
  help               Show this help message

Environment Variables:
  PULUMI_ENV         Environment (local|beta|prod) - default: local
  PULUMI_STACK       Stack name - default: default
  PULUMI_CONFIG_PASSPHRASE  Pulumi config passphrase
  AWS_PROFILE        AWS profile to use
  AWS_REGION         AWS region - default: us-east-2

Examples:
  # Run preview with local backend (no AWS credentials needed)
  ./scripts/pulumi-docker.sh run-local

  # Run preview for beta environment
  PULUMI_ENV=beta ./scripts/pulumi-docker.sh preview

  # Initialize a new stack
  ./scripts/pulumi-docker.sh stack init beta

  # Set configuration values
  ./scripts/pulumi-docker.sh config set appEnv beta
  ./scripts/pulumi-docker.sh config set vpcCidr 10.0.0.0/16

  # Run full preview
  ./scripts/pulumi-docker.sh preview

  # Apply changes
  ./scripts/pulumi-docker.sh up
EOF
        ;;

    *)
        log_error "Unknown command: $COMMAND"
        echo "Run '$0 help' for usage information"
        exit 1
        ;;
esac
