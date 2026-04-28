# Infrastructure Bootstrap — subject-data

Manual bootstrap steps for provisioning a new subject-data AWS account.
Run these **once per account** before Pulumi can manage ongoing infrastructure.

Each step starts with an account verification check. **Read the expected output before running anything.**

---

## Prerequisites

- AWS CLI v2 installed
- Access to the sovra-ai management account (770826159442) for account creation
- Access to the target account (via root or IAM) for provisioning steps

---

## Step 0 — Create the AWS account

Run from the **management account** (770826159442).

```bash
# If you're able to log in using a username & password in the browser, you can do the same for the AWS CLI:
aws login --profile me-in-sovra-root

# Verify you're in the management account
aws sts get-caller-identity --profile me-in-sovra-root
# Expected: Account 770826159442

# Create the account
aws organizations create-account \
  --profile me-in-sovra-root \
  --email "tech+subject-data-<ENV>@sovra.ai" \
  --account-name "subject-data-<ENV>"

# Poll until SUCCEEDED (takes 30-60s)
aws organizations list-create-account-status \
  --profile me-in-sovra-root \
  --states IN_PROGRESS \
  --query 'CreateAccountStatuses[?AccountName==`subject-data-<ENV>`]'
```

Record the new account ID — you'll need it for every subsequent step.

---

## Step 1 — Assume a role into the new account

New accounts in an Organization have a default `OrganizationAccountAccessRole` that the management account can assume.

```bash
# Replace ACCOUNT_ID with the new account ID from Step 0
ACCOUNT_ID=<ACCOUNT_ID>
ENV=<ENV>  # beta or prod

# Add a profile
echo "[profile bootstrap-subject-data-${ENV}]
role_arn = arn:aws:iam::${ACCOUNT_ID}:role/OrganizationAccountAccessRole
source_profile = me-in-sovra-root
role_session_name = bootstrap-subject-data-${ENV}
region = us-east-2" >> ~/.aws/config

# Verify you're in the right account
aws sts get-caller-identity --profile bootstrap-subject-data-${ENV}
# Expected: Account $ACCOUNT_ID, Arn contains "OrganizationAccountAccessRole"
```

---

## Step 2 — S3 state bucket for Pulumi

One bucket per account. Pulumi state for this service-environment lives here.

```bash
BUCKET="pulumi-state-subject-data-${ENV}"

aws s3api create-bucket \
  --profile bootstrap-subject-data-${ENV} \
  --bucket "${BUCKET}" \
  --region us-east-2 \
  --create-bucket-configuration LocationConstraint=us-east-2

aws s3api put-bucket-versioning \
  --profile bootstrap-subject-data-${ENV} \
  --bucket "${BUCKET}" \
  --versioning-configuration Status=Enabled

aws s3api put-bucket-encryption \
  --profile bootstrap-subject-data-${ENV} \
  --bucket "${BUCKET}" \
  --server-side-encryption-configuration '{
    "Rules": [{"ApplyServerSideEncryptionByDefault": {"SSEAlgorithm": "AES256"}}]
  }'

aws s3api put-public-access-block \
  --profile bootstrap-subject-data-${ENV} \
  --bucket "${BUCKET}" \
  --public-access-block-configuration '{
    "BlockPublicAcls": true,
    "IgnorePublicAcls": true,
    "BlockPublicPolicy": true,
    "RestrictPublicBuckets": true
  }'

# Verify
aws s3api head-bucket --profile bootstrap-subject-data-${ENV} --bucket "${BUCKET}" && echo "OK: ${BUCKET} exists"
```

---

## Step 3 — ECR repository

ECR is a bootstrap resource (like the S3 state bucket) because the CI/CD pipeline pushes images before `pulumi up` runs. Pulumi looks up this repo by name — it does not create it.

```bash
aws ecr create-repository \
  --profile bootstrap-subject-data-${ENV} \
  --repository-name subject-data \
  --region us-east-2 \
  --image-scanning-configuration scanOnPush=true \
  --encryption-configuration encryptionType=AES256

# Lifecycle policy: keep last 20 images, expire untagged after 7 days
aws ecr put-lifecycle-policy \
  --profile bootstrap-subject-data-${ENV} \
  --repository-name subject-data \
  --region us-east-2 \
  --lifecycle-policy-text '{
    "rules": [
      {
        "rulePriority": 1,
        "description": "Expire untagged images after 7 days",
        "selection": {
          "tagStatus": "untagged",
          "countType": "sinceImagePushed",
          "countUnit": "days",
          "countNumber": 7
        },
        "action": { "type": "expire" }
      },
      {
        "rulePriority": 2,
        "description": "Keep last 20 tagged images",
        "selection": {
          "tagStatus": "tagged",
          "tagPrefixList": ["sha-"],
          "countType": "imageCountMoreThan",
          "countNumber": 20
        },
        "action": { "type": "expire" }
      }
    ]
  }'

# Verify
aws ecr describe-repositories \
  --profile bootstrap-subject-data-${ENV} \
  --repository-names subject-data --region us-east-2 \
  --query 'repositories[0].repositoryUri' --output text
```

---

## Step 4 — GitHub Actions OIDC provider

Only needed once per account. If multiple repos deploy to the same account, they share this provider.

```bash
# Check if it already exists (skip creation if so)
aws iam list-open-id-connect-providers \
  --profile bootstrap-subject-data-${ENV} \
  --query "OpenIDConnectProviderList[?ends_with(Arn, 'token.actions.githubusercontent.com')]"

# Create the OIDC provider
THUMBPRINT="d89e3bd43d5d909b47a18977aa9d5ce36cee184c"

aws iam create-open-id-connect-provider \
  --profile bootstrap-subject-data-${ENV} \
  --url "https://token.actions.githubusercontent.com" \
  --client-id-list "sts.amazonaws.com" \
  --thumbprint-list "${THUMBPRINT}"
```

> **Note:** The thumbprint above is GitHub's current OIDC thumbprint. AWS doesn't actually
> validate it for GitHub (they use a library of trusted CAs), but the field is required.
> If in doubt, verify at https://github.blog/changelog/2023-06-27-github-actions-update-on-oidc-integration-with-aws/

---

## Step 5 — IAM deploy role for GitHub Actions

This role is assumed by GitHub Actions via OIDC. Scoped to the `sovraai/subject-data` repo.

```bash
# Get the OIDC provider ARN
OIDC_ARN=$(aws iam list-open-id-connect-providers \
  --profile bootstrap-subject-data-${ENV} \
  --query "OpenIDConnectProviderList[?ends_with(Arn, 'token.actions.githubusercontent.com')].Arn" \
  --output text)

# Create the trust policy
cat > /tmp/trust-policy.json << TRUSTEOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": { "Federated": "${OIDC_ARN}" },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:sovraai/subject-data:*"
        }
      }
    }
  ]
}
TRUSTEOF

aws iam create-role \
  --profile bootstrap-subject-data-${ENV} \
  --role-name github-actions-infra-deploy \
  --assume-role-policy-document file:///tmp/trust-policy.json \
  --description "GitHub Actions deploy role for subject-data (${ENV})"

# Attach permissions — start broad, tighten after first successful deploy
aws iam attach-role-policy \
  --profile bootstrap-subject-data-${ENV} \
  --role-name github-actions-infra-deploy \
  --policy-arn arn:aws:iam::aws:policy/AdministratorAccess

rm /tmp/trust-policy.json

# Verify
aws iam get-role \
  --profile bootstrap-subject-data-${ENV} \
  --role-name github-actions-infra-deploy \
  --query 'Role.Arn' --output text
```

> **TODO:** Replace `AdministratorAccess` with a scoped policy after the first deploy
> succeeds and we know exactly which permissions Pulumi needs. Track in a follow-up ticket.

---

## Step 6 — Generate Pulumi passphrase

One passphrase per repo, shared across all environments. Store a backup in the management account.

```bash
# Generate
PASSPHRASE=$(openssl rand -base64 32)

# Verify you're in the management account
aws sts get-caller-identity --profile me-in-sovra-root
# Expected: Account 770826159442

# Store in Secrets Manager as a break-glass backup
aws secretsmanager create-secret \
  --profile me-in-sovra-root \
  --name "pulumi/subject-data/config-passphrase" \
  --description "Pulumi config passphrase for sovraai/subject-data (all envs)" \
  --secret-string "${PASSPHRASE}" \
  --region us-east-2

# Set as a repo-level secret in GitHub (available to all environments)
gh secret set PULUMI_CONFIG_PASSPHRASE \
  --repo sovraai/subject-data \
  --body "${PASSPHRASE}"
```

> **Skip the Secrets Manager step** if the passphrase already exists (e.g. bootstrapping prod
> after beta). Retrieve it with:
> ```bash
> aws secretsmanager get-secret-value \
>   --profile me-in-sovra-root \
>   --secret-id "pulumi/subject-data/config-passphrase" \
>   --region us-east-2 \
>   --query SecretString --output text
> ```

---

## Step 7 — Initialize Pulumi stack

The stack must exist in the S3 backend before CI can run `pulumi up`. Run from the repo's `infra/pulumi/` directory with **beta account creds assumed** (Step 1).

```bash
cd infra/pulumi

#  The following command is wrapped in a subshell so the env vars don't leak:
(
  eval $(aws configure export-credentials --profile bootstrap-subject-data-${ENV} --format env)
  AWS_REGION=us-east-2 \
  PULUMI_BACKEND_URL=s3://pulumi-state-subject-data-${ENV} \
  PULUMI_CONFIG_PASSPHRASE="${PASSPHRASE}" \
  pulumi stack init ${ENV}
)
```

> The `PASSPHRASE` variable should still be in your shell from Step 6. If not, retrieve it:
> ```bash
> PASSPHRASE=$(aws secretsmanager get-secret-value \
>   --profile me-in-sovra-root \
>   --secret-id "pulumi/subject-data/config-passphrase" \
>   --region us-east-2 \
>   --query SecretString --output text)
> ```

---

## Step 8 — Configure GitHub Environment

```bash
REPO="sovraai/subject-data"
ROLE_ARN="<role ARN from Step 5>"
ENV=<ENV>  # beta or prod

# Create the environment
gh api --method PUT "repos/${REPO}/environments/${ENV}"

# Add the deploy role ARN as an environment variable
gh variable set AWS_DEPLOY_ROLE_ARN \
  --repo "${REPO}" \
  --env "${ENV}" \
  --body "${ROLE_ARN}"

# Add the Slack webhook URL for failure notifications (repo-level, shared across envs)
gh secret set SLACK_WEBHOOK_URL \
  --repo "${REPO}" \
  --body "<your-slack-webhook-url>"

# Verify
gh variable list --repo "${REPO}" --env "${ENV}"
gh secret list --repo "${REPO}"
```

---

## Step 9 — Verify end-to-end

```bash
# Verify the account exists
aws organizations list-accounts \
  --profile me-in-sovra-root \
  --query "Accounts[?Name=='subject-data-${ENV}'].{Id:Id,Name:Name,Status:Status}" \
  --output table
```

Then trigger a CI run from the repo to confirm OIDC auth works.

---

## Step 10 — Populate secrets

Prepare to use pulumi commands by "logging in" to use the bucket:

```bash
(
  eval $(aws configure export-credentials --profile bootstrap-subject-data-${ENV} --format env)
  AWS_REGION=us-east-2 \
  PULUMI_BACKEND_URL=s3://${BUCKET} \
  PULUMI_CONFIG_PASSPHRASE="${PASSPHRASE}" \
  pulumi login s3://${BUCKET}
)

(
  eval $(aws configure export-credentials --profile bootstrap-subject-data-${ENV} --format env)
  AWS_REGION=us-east-2 \
  PULUMI_BACKEND_URL=s3://${BUCKET} \
  PULUMI_CONFIG_PASSPHRASE="${PASSPHRASE}" \
  pulumi up
)
```

Run `pulumi up` to create the Secrets Manager resources, populate their values.
Run from the **beta account** (assumed role from Step 1).

### Tailscale auth key

Generate in the [Tailscale admin console](https://login.tailscale.com/admin/settings/keys) → Settings → Keys → Generate auth key.
Use: **Reusable** + **Ephemeral** (so ECS tasks can cycle without burning the key and nodes auto-deregister on task stop). Max expiry is 90 days — rotate manually until automation is in place.

```bash
aws secretsmanager put-secret-value \
  --profile bootstrap-subject-data-${ENV} \
  --secret-id "subject-data/${ENV}/ts-authkey" \
  --secret-string "<key from Tailscale admin>" \
  --region us-east-2
```

### Optional - add Tailscale URLs

The Tailscale sidecar connects the task to the tailnet. Use the Tailscale hostname for any node you connect to

```bash
aws secretsmanager put-secret-value \
  --profile bootstrap-subject-data-${ENV} \
  --secret-id "subject-data/${ENV}/es-url" \
  --secret-string "http://dev:9200" \
  --region us-east-2
```

### API auth tokens

Server-to-server bearer tokens that callers (e.g. frontend-server) send to authenticate with the subject-data API. Generate and share with the calling service — they store the same value in their own secrets namespace.

```bash
AUTH_TOKEN=$(openssl rand -base64 32)

# Store in subject-data's secrets
aws secretsmanager put-secret-value \
  --profile bootstrap-subject-data-${ENV} \
  --secret-id "subject-data/${ENV}/api-auth-tokens" \
  --secret-string "${AUTH_TOKEN}" \
  --region us-east-2

echo "Share this token with the calling service (e.g. frontend-server)."
echo "Token: ${AUTH_TOKEN}"
```

### Verify all secrets have values

```bash
for secret in ts-authkey es-url api-auth-tokens; do
  echo -n "subject-data/${ENV}/${secret}: "
  aws secretsmanager get-secret-value \
    --profile bootstrap-subject-data-${ENV} \
    --secret-id "subject-data/${ENV}/${secret}" \
    --region us-east-2 \
    --query 'VersionId' --output text
done
```

Then restart the ECS service to pick up the secrets:

```bash
aws ecs update-service \
  --profile bootstrap-subject-data-${ENV} \
  --cluster "subject-data-${ENV}" \
  --service "subject-data-${ENV}" \
  --force-new-deployment \
  --region us-east-2
```

---

## Checklist

- [ ] Account created in AWS Organizations
- [ ] S3 state bucket created with versioning + encryption
- [ ] ECR repository created with lifecycle policy
- [ ] OIDC provider created
- [ ] IAM deploy role created with OIDC trust
- [ ] Pulumi passphrase generated, stored in Secrets Manager + GitHub repo secret
- [ ] Pulumi stack initialized (`pulumi stack init`)
- [ ] GitHub Environment configured (AWS_DEPLOY_ROLE_ARN variable)
- [ ] Secrets populated (ts-authkey, es-url, api-auth-tokens)
- [ ] ECS service restarted and tasks healthy
- [ ] End-to-end verification passed

---

## Applying to other services

This same process applies to any new `<service>-<env>` account. Replace:
- `subject-data` → service name
- `sovraai/subject-data` → GitHub repo in the OIDC trust policy
- Bucket name, ECR repo name accordingly
