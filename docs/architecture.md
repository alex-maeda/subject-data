# Subject Search — Infrastructure Architecture

## Overview

Subject Data Service a simple server-to-server API deployed to its own AWS account (one per environment). Calls are authenticated via cross-account IAM roles.

## Request Flow

```mermaid
flowchart LR
    subgraph caller["Caller Account"]
        CS[caller-server]
    end

    subgraph ss["Template Service Account"]
        subgraph apigw["API Gateway (HTTP API)"]
            IAM_AUTH[IAM Auth]
        end

        VPC_LINK[VPC Link]

        subgraph vpc["VPC"]
            subgraph ecs["ECS Fargate Task"]
                APP[subject-data\n:8380]
                TS[tailscale sidecar\nSOCKS5 :1055]
            end
            CM[Cloud Map\nService Discovery]
        end

        subgraph rds["RDS PostgreSQL"]
            DB[(db.t4g.micro\nPostgreSQL 16)]
        end

        subgraph secrets["Secrets Manager"]
            S1[ts-authkey]
            S2[database-url]
            S3[api-auth-tokens]
        end
    end

    subgraph hetzner["Hetzner"]
        HS[SerivceOnHetzner]
    end

    CS -- "1. AssumeRole\n(cross-account)" --> INVOKE_ROLE
    INVOKE_ROLE[Invoke Role] --> IAM_AUTH
    IAM_AUTH --> VPC_LINK
    VPC_LINK --> CM
    CM --> APP
    APP -- "SOCKS5 proxy" --> TS
    TS -- "Tailscale\ntailnet" --> HS
    APP -- "PostgreSQL\n:5432" --> DB
    secrets -. "injected at\ntask start" .-> ecs
```

## Cross-Account Auth

```mermaid
sequenceDiagram
    participant CS as caller-server<br/>(caller account)
    participant STS as AWS STS
    participant APIGW as API Gateway<br/>(subject-data account)
    participant ECS as ECS Task

    CS->>STS: AssumeRole(invoke-role ARN,<br/>ExternalId: subject-data-ENV)
    STS-->>CS: Temporary credentials
    CS->>APIGW: Signed request (SigV4)
    APIGW->>APIGW: Verify IAM auth
    APIGW->>ECS: VPC Link → Cloud Map → task :8380
    ECS-->>APIGW: Response
    APIGW-->>CS: Response
```

## Key Details

| Component             | Detail                                                   |
|-----------------------|----------------------------------------------------------|
| **Database**          | RDS PostgreSQL 16, db.t4g.micro, 20GB gp3, single-AZ    |
| **Launch type**       | Fargate (512 CPU / 1024 MB)                              |
| **Tailscale mode**    | Userspace (`TS_USERSPACE=true`, SOCKS5 on `:1055`)       |
| **HS connectivity**   | App → SOCKS5 proxy → Tailscale tailnet → `dev:9200`      |
| **API auth**          | IAM (SigV4) at API Gateway, cross-account invoke role    |
| **Service discovery** | Cloud Map private DNS (`subject-data-ENV.local`)         |
| **Secrets**           | Secrets Manager, injected by ECS at task start           |
| **State**             | Pulumi S3 backend (`s3://pulumi-state-subject-data-ENV`) |
| **Region**            | `us-east-2`                                              |
