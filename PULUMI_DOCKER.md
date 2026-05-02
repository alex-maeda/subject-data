# Running Pulumi in Docker

This guide explains how to run `pulumi preview` and other Pulumi commands using Docker, both with a local backend (no AWS credentials needed) and with AWS-backed state.

## Overview

The project includes Docker-based tooling to run Pulumi commands in isolated containers. This is useful for:
- Running Pulumi without installing it locally
- Ensuring consistent Pulumi and Go versions across environments
- Testing infrastructure changes in a clean environment
- Running Pulumi in CI/CD pipelines

## Quick Start (Local Backend - No AWS Required)

For testing and development without real AWS resources:

```bash
# Build the Pulumi Docker image
docker compose -f docker-compose.yml build pulumi-local

# Run pulumi preview with local backend
./scripts/pulumi-docker.sh run-local preview

# Or use docker compose directly:
docker compose -f docker-compose.yml run --rm pulumi-local \
  pulumi preview --stack local
```

## Using AWS-Backed State

For production use with real AWS resources:

### 1. Configure AWS Credentials

Ensure AWS credentials are available in the container:

```bash
# Option A: Use AWS profile (recommended)
export AWS_PROFILE=your-profile

# Option B: Set environment variables
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_SESSION_TOKEN=...
```

### 2. Set Pulumi Passphrase

```bash
export PULUMI_CONFIG_PASSPHRASE="your-passphrase"
```

### 3. Login to Pulumi Backend

```bash
./scripts/pulumi-docker.sh login
```

Or manually:

```bash
docker compose -f docker-compose.yml run --rm pulumi \
  pulumi login s3://pulumi-state-subject-data-beta
```

### 4. Initialize Stack

```bash
./scripts/pulumi-docker.sh stack init beta
```

### 5. Configure Stack

```bash
# Set required configuration
./scripts/pulumi-docker.sh config set appEnv beta
./scripts/pulumi-docker.sh config set vpcCidr 10.0.0.0/16

# Set caller accounts (JSON)
./scripts/pulumi-docker.sh config set callerAccounts '{"sovraai/frontend-server":"arn:aws:iam::174581551884:role/apigw-invoke-role-frontend-server-ef18c68"}'
```

### 6. Run Preview

```bash
./scripts/pulumi-docker.sh preview
```

### 7. Apply Changes

```bash
./scripts/pulumi-docker.sh up
```

## Docker Compose Services

### pulumi-local

- **Purpose**: Local development and testing
- **Backend**: Local file system (`file:///root/.pulumi`)
- **AWS**: Not required
- **Use Case**: Testing Pulumi code changes, dry runs

### pulumi

- **Purpose**: Production deployments
- **Backend**: S3 bucket (configurable via `PULUMI_BACKEND_URL`)
- **AWS**: Required
- **Use Case**: Real infrastructure deployments

## Helper Script Usage

The `scripts/pulumi-docker.sh` script provides convenient commands:

```bash
# Show help
./scripts/pulumi-docker.sh help

# Login to backend
./scripts/pulumi-docker.sh login

# Initialize a stack
./scripts/pulumi-docker.sh stack init beta

# List stacks
./scripts/pulumi-docker.sh stack list

# Set configuration
./scripts/pulumi-docker.sh config set appEnv beta
./scripts/pulumi-docker.sh config set vpcCidr 10.0.0.0/16

# View configuration
./scripts/pulumi-docker.sh config

# Run preview (dry-run)
./scripts/pulumi-docker.sh preview

# Apply changes
./scripts/pulumi-docker.sh up

# Destroy all resources
./scripts/pulumi-docker.sh destroy

# Refresh state
./scripts/pulumi-docker.sh refresh

# View outputs
./scripts/pulumi-docker.sh output

# Interactive shell
./scripts/pulumi-docker.sh shell
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PULUMI_ENV` | Environment (local, beta, prod) | `local` |
| `PULUMI_STACK` | Stack name | `default` |
| `PULUMI_CONFIG_PASSPHRASE` | Pulumi config passphrase | - |
| `PULUMI_BACKEND_URL` | Backend URL (s3:// or file://) | `file:///root/.pulumi` |
| `AWS_PROFILE` | AWS profile name | `default` |
| `AWS_REGION` | AWS region | `us-east-2` |
| `AWS_ACCESS_KEY_ID` | AWS access key | - |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key | - |
| `AWS_SESSION_TOKEN` | AWS session token | - |

## Stack Configuration

Required configuration for each stack:

```bash
# App environment (beta, prod)
pulumi config set appEnv beta

# VPC CIDR block
pulumi config set vpcCidr 10.0.0.0/16

# Caller accounts (JSON map)
pulumi config set callerAccounts '{"sovraai/frontend-server":"arn:aws:iam::174581551884:role/apigw-invoke-role-frontend-server-ef18c68"}'
```

## Examples

### Example 1: Local Preview (No AWS)

```bash
# Build and run local preview
docker compose -f docker-compose.yml build pulumi-local
docker compose -f docker-compose.yml run --rm pulumi-local \
  pulumi preview --stack local
```

### Example 2: Beta Environment Preview

```bash
# Set environment
export PULUMI_ENV=beta
export PULUMI_CONFIG_PASSPHRASE="your-passphrase"
export AWS_PROFILE=beta-account

# Login and preview
./scripts/pulumi-docker.sh login
./scripts/pulumi-docker.sh preview
```

### Example 3: Initialize New Stack

```bash
# Initialize beta stack
./scripts/pulumi-docker.sh stack init beta

# Configure
./scripts/pulumi-docker.sh config set appEnv beta
./scripts/pulumi-docker.sh config set vpcCidr 10.0.0.0/16

# Preview
./scripts/pulumi-docker.sh preview
```

### Example 4: Deploy to Production

```bash
# Switch to prod
export PULUMI_ENV=prod
export PULUMI_STACK=prod
export PULUMI_CONFIG_PASSPHRASE="your-passphrase"

# Login
./scripts/pulumi-docker.sh login

# Preview
./scripts/pulumi-docker.sh preview

# Deploy (after review)
./scripts/pulumi-docker.sh up
```

## Troubleshooting

### Issue: "passphrase incorrect"

**Solution**: Ensure `PULUMI_CONFIG_PASSPHRASE` is set correctly:

```bash
export PULUMI_CONFIG_PASSPHRASE="your-passphrase"
```

### Issue: "failed to login to backend"

**Solution**: Check AWS credentials and backend URL:

```bash
aws sts get-caller-identity  # Verify AWS credentials
echo $PULUMI_BACKEND_URL      # Verify backend URL
```

### Issue: "error creating container"

**Solution**: Rebuild the Docker image:

```bash
docker compose -f docker-compose.yml build --no-cache
```

### Issue: "plugin not found"

**Solution**: Install required plugins:

```bash
./scripts/pulumi-docker.sh build
```

## CI/CD Integration

For GitHub Actions or other CI/CD systems:

```yaml
- name: Run Pulumi Preview
  run: |
    docker compose -f docker-compose.yml build pulumi
    docker compose -f docker-compose.yml run --rm pulumi \
      pulumi preview --stack beta
  env:
    PULUMI_CONFIG_PASSPHRASE: ${{ secrets.PULUMI_CONFIG_PASSPHRASE }}
    AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
    AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    AWS_REGION: us-east-2
```

## File Structure

```
.
 docker-compose.yml          # Docker Compose configuration
 Dockerfile.pulumi          # Pulumi Docker image
 scripts/
    pulumi-docker.sh       # Helper script
 .env.pulumi.example        # Environment variables template
 infra/pulumi/              # Pulumi project
     Pulumi.yaml
     main.go
     go.mod
```

## Notes

- The local backend (`file://`) stores state locally and is not shared
- For team collaboration, use S3 backend with proper IAM permissions
- Always review `pulumi preview` output before running `pulumi up`
- Keep your Pulumi passphrase secure - it encrypts your stack's secrets
- The Docker image includes Go for building plugins if needed

## See Also

- [Pulumi CLI Documentation](https://www.pulumi.com/docs/cli/)
- [Pulumi AWS Provider](https://www.pulumi.com/registry/packages/aws/)
- [Infrastructure Bootstrap Guide](docs/infra-bootstrap.md)
