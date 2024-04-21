# Cirrostratus Cloud - Auth AWS

## Requirements

First create an dotenv file (`.env`), like this:

```shell
AWS_STAGE=prod
LOG_LEVEL=INFO
CIRROSTRATUS_AUTH_MODULE_NAME=cirrostratus-auth
CIRROSTRATUS_AUTH_USER_TABLE=users
AWS_DEFAULT_REGION=us-west-1
AWS_REGION=us-west-1
AWS_ACCESS_KEY_ID=<your-access-key>
AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
USER_MIN_PASSWORD_LENGTH=8
USER_UPPER_CASE_REQUIRED=true
USER_LOWER_CASE_REQUIRED=true
USER_NUMBER_REQUIRED=true
USER_SPECIAL_CHARACTER_REQUIRED=true
TOPIC_ARN_PREFIX=arn:aws:sns:us-west-1:*:cirrostratus-auth_
SES_EMAIL_FROM=cirrostratus@cloud.com
SES_EMAIl_ARN=arn:aws:ses:us-west-1:*:identity/cirrostratus@cloud.com
SES_CONFIGURATION_SET=arn:aws:ses:us-west-1:*:configuration-set/my-first-configuration-set
EMAIL_CONFIRMATION_URL=https://cloud.com/confirm-email
PRIVATE_KEY=s3cr3t
MAX_AGE_IN_SECONDS=3600
```

## Run locally

```shell
task serve_http
```

```shell
task serve_event
```

## Deploy

```shell
task deploy
```