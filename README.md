# CoreOS rkt Authentication with AWS ECR via Golang

Writes to stdout JSON for rkt authentication with AWS ECR registery. The JSON object should be written to `/etc/rkt/auth.d/<registery-name>` for use by CoreOS rkt.

This should be run on a cron/timer because the AWS ECR token will expire every 12 hours.

**Example:**

```
AWS_REGION=us-east-1 AWS_REGISTERY_ID=384322436518 rkt-aws-ecr-auth > /etc/rkt/auth.d/ecr.json
```
