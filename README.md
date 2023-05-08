# DMM schedule checker

DMM schedule checker continuously monitors the schedule of your favorite teachers, and notifies via LINE whenever new slots are available.

## Usage

You need to have [Docker CLI](https://github.com/docker/cli) (or [Fintch](https://github.com/runfinch/finch)), [AWS CLI](https://github.com/aws/aws-cli) and [Terraform](https://github.com/hashicorp/terraform) installed on your machine.

You can deploy DMM schedule checker to your AWS account as follows.

```bash
# you need LINE access token to send messages
LINE_ACCESS_TOKEN="YOUR_TOKEN"

# set env
AWS_REGION=ap-northeast-1
ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text)
IMAGE=${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/dmm_schedule_checker

# create S3 bucket to store tfstate
aws s3 mb s3://dmm-schedule-checker-${ACCOUNT_ID} --region ${AWS_REGION}

# build image
docker build -t ${IMAGE} --target prod app

# initialize
cd terraform
terraform init \
    -backend-config="bucket=dmm-schedule-checker-${ACCOUNT_ID}" \
    -backend-config="region=${AWS_REGION}" -reconfigure

# create ECR repo
terraform apply --target=aws_ecr_repository.default \
    -var="line_access_token=${LINE_ACCESS_TOKEN}" -auto-approve

# push image
aws ecr get-login-password --region ${AWS_REGION} | \
    docker login --username AWS --password-stdin \
    ${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
docker push ${IMAGE}:latest

# deploy application
terraform apply -var="line_access_token=${LINE_ACCESS_TOKEN}" -auto-approve
```

To add/delete your favorite teacher, you can simple call the API as follows. (Currently, we don't have frontend application for it.)

```bash
# add teacher
curl -X POST -H "Content-Type: application/json" \
    -d '{"id": "5_DIGITS_TEACHER_ID"}' YOUR_APP_RUNNER_ENDPOINT_URL/teacher
# delete teacher
curl -X DELETE -H "Content-Type: application/json" \
    -d '{"id": "5_DIGITS_TEACHER_ID"}' YOUR_APP_RUNNER_ENDPOINT_URL/teacher
```

The application checks the schedule of enrolled teachers every 20 seconds, and notifies via LINE whenever new slots are available.
