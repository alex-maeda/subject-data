# subject-data

REST API for reading and writing Subjects, their Properties, and the Records used to derive those properties. Server-to-server only — never called directly from the web app.

# For Developers

## Prerequisites

- Go 1.26+
- golangci-lint (for `make lint`)

## Quick start

```bash
# Run the server (default port 8380)
make run
```

The server listens on `:8380` by default. Override with `PORT`:

```bash
PORT=9090 make run
```

## Commands

| Command | What it does |
|---|---|
| `make build` | Compile binary to `bin/run_service` |
| `make run` | Start dev server |
| `make test` | Run all tests |
| `make lint` | Run golangci-lint |
| `make fmt` | Format all Go files |

## Environment variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8380` | Server listen port |
| `AUTH_TOKENS` | _(unset)_ | Comma-separated Bearer tokens for S2S auth. If unset, auth is disabled (local dev). |
| `CORS_ALLOW_ORIGINS` | _(unset)_ | When set, enables CORS for the given origin (e.g. `http://localhost:5173`). Leave unset in production. |

## Project structure

```
cmd/run_service/     Entry point, signal handling
internal/
  api/               HTTP server, router, handlers, request/response types
  middleware/         Bearer token auth middleware
infra/
  pulumi/            AWS infrastructure (ECS Fargate, API Gateway, VPC)
```

# For Clients

AWS accounts:
* Beta: `174581551884`
* Prod: `493628259544`

## Calling as a human on a laptop

1. Your best bet right not is to use Tailscale, where you can hit the API or use the Swagger UI: http://subject-data-beta-1.tail5ab057.ts.net:8380/docs/index.html
2. You'll need an auth token.
   1. If you switch AWS roles into 174581551884, you can see auth tokens as an ASM secret (`subject-data/beta/api-auth-tokens`) .
   2. In the Swagger UI, the top right should have an "Authorize" button. Click that, and enter "Bearer {the-auth-token-you-got-from-secrets-manager}"

## Calling as another service

1. There is an API Gateway, with invokeURL https://ytegidxhv5.execute-api.us-east-2.amazonaws.com
2. The client will need the same auth key that you used for Swagger
3. The client will need an invoking role
   1. For `frontend-server`, it already exists: `arn:aws:iam::174581551884:role/apigw-invoke-role-frontend-server-ef18c68`
   2. For others, we'll need to update the `subject-data:callerAccounts` config in `infra/pulumu/Pulumi.beta.yaml`.
   3. If you switch roles into `174581551884`, you can view the details, but I'll also put them here:
      1. Trust relationship:
      ```json
      {
          "Version": "2012-10-17",
          "Statement": [
              {
                  "Effect": "Allow",
                  "Principal": {
                      "AWS": "arn:aws:iam::835107812595:root"
                  },
                  "Action": "sts:AssumeRole",
                  "Condition": {
                      "StringEquals": {
                          "sts:ExternalId": "subject-data-beta"
                      }
                  }
              }
         ]
      }
      ```
      2. Permissions:
      ```json
      {
          "Version": "2012-10-17",
          "Statement": [
              {
                  "Action": "execute-api:Invoke",
                  "Effect": "Allow",
                  "Resource": "arn:aws:execute-api:us-east-2:174581551884:ytegidxhv5/*"
              }
          ]
      }
      ```