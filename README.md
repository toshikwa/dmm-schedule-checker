# DMM schedule checker

DMM schedule checker continuously monitors the schedule of your favorite teachers, and notifies via LINE whenever new slots are available.

## Installation

You need to have [Docker CLI](https://github.com/docker/cli), [AWS CLI](https://github.com/aws/aws-cli) and [Terraform](https://github.com/hashicorp/terraform) installed on your machine.

### Generate LINE Notify access token

Then, you need to generate [LINE Notify](https://notify-bot.line.me/) access token.

- go to [mypage](https://notify-bot.line.me/my/)
- click "Generate token"
- select "1-on-1 chat with LINE Notify" and generate

Let's save the generated token in the `.env.local` file.

```bash
echo "LINE_NOTIFY_ACCESS_TOKEN=[YOUR_ACCESS_TOKEN]" >> .env.local
```

### Deploy backend API

You can deploy backend API for DMM schedule checker to your AWS account as follows.

```bash
# set envs
AWS_REGION=ap-northeast-1
ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)

# create S3 bucket to store tfstate
aws s3 mb s3://dmm-schedule-checker-${ACCOUNT_ID} --region ${AWS_REGION}

# initialize terraform
terraform init \
    -backend-config="bucket=dmm-schedule-checker-${ACCOUNT_ID}" \
    -backend-config="region=${AWS_REGION}" -reconfigure

# deploy application
terraform apply -auto-approve
```

## Usage

To add/delete your favorite teacher, you can simple call the API as follows. Currently, we don't have frontend application for it.

```bash
# API endpoint
ENDPOINT_URL=https://$(terraform output --raw endpoint)

# add teacher
curl -X POST -H "Content-Type: application/json" \
    -d '{"id": "5_DIGIT_TEACHER_ID"}' ${ENDPOINT_URL}/teachers

# delete teacher
curl -X DELETE ${ENDPOINT_URL}/teachers/5_DIGIT_TEACHER_ID
```

The application checks the schedule of enrolled teachers every 20 seconds, and notifies via LINE whenever new slots are available.

## Clean up

You can clean up all resources as follows.

```bash
# set envs
AWS_REGION=ap-northeast-1
ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)

# clean up api resources
terraform destroy -auto-approve

# remove tfstate file and s3 bucket
aws s3 rm s3://dmm-schedule-checker-${ACCOUNT_ID}/terraform.tfstate --region ${AWS_REGION}
aws s3 rb s3://dmm-schedule-checker-${ACCOUNT_ID} --region ${AWS_REGION}
```
