data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

terraform {
  required_version = ">= 1.4.0"
  backend "s3" {
    key     = "terraform.tfstate"
    encrypt = true
  }
}

provider "aws" {}

locals {
  app = "dmm_schedule_checker"
}

locals {
  envs = { for tuple in regexall("(.*)=(.*)", file(".env.local")) : tuple[0] => sensitive(tuple[1]) }
}

output "endpoint" {
  value = aws_apprunner_service.default.service_url
}
