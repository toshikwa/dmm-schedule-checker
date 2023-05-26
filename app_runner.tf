// ecr
resource "aws_ecr_repository" "default" {
  name                 = local.app
  image_tag_mutability = "MUTABLE"
  force_delete         = true
  image_scanning_configuration {
    scan_on_push = true
  }
}
data "aws_ecr_authorization_token" "default" {}
resource "null_resource" "image_push" {
  triggers = {
    code_diff = join("", [for file in fileset("app/", "**/*.go") : filebase64("app/${file}")])
  }
  provisioner "local-exec" {
    command = <<-EOF
      docker build -t ${aws_ecr_repository.default.repository_url}:latest app
      docker login -u AWS -p ${data.aws_ecr_authorization_token.default.password} ${data.aws_ecr_authorization_token.default.proxy_endpoint}
      docker push ${aws_ecr_repository.default.repository_url}:latest
    EOF
  }
  depends_on = [aws_ecr_repository.default]
}

// iam - access role
resource "aws_iam_role" "access_role" {
  name = "${local.app}_access_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "build.apprunner.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}
resource "aws_iam_role_policy_attachment" "access_role" {
  role       = aws_iam_role.access_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSAppRunnerServicePolicyForECRAccess"
}

// iam - instance role
resource "aws_iam_role" "instance_role" {
  name = "${local.app}_instance_role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = "sts:AssumeRole"
        Principal = {
          Service = ["tasks.apprunner.amazonaws.com"]
        }
      }
    ]
  })
}
data "aws_iam_policy_document" "instance_role" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:*"]
    resources = ["*"]
  }
}
resource "aws_iam_policy" "instance_role" {
  name   = "${local.app}_instance_role_policy"
  policy = data.aws_iam_policy_document.instance_role.json
}
resource "aws_iam_role_policy_attachment" "instance_role" {
  role       = aws_iam_role.instance_role.name
  policy_arn = aws_iam_policy.instance_role.arn
}

// wait for 5 secs to avoid error
resource "time_sleep" "wait" {
  create_duration = "5s"
  depends_on = [
    null_resource.image_push,
    aws_iam_role_policy_attachment.access_role,
    aws_iam_role_policy_attachment.instance_role,
  ]
}

// app runner service
resource "aws_apprunner_service" "default" {
  service_name = local.app
  source_configuration {
    authentication_configuration {
      access_role_arn = aws_iam_role.access_role.arn
    }
    image_repository {
      image_configuration {
        port = 8000
        runtime_environment_variables = {
          AWS_REGION               = data.aws_region.current.name
          PORT                     = 8000
          TEACHER_TABLE_NAME       = aws_dynamodb_table.teacher_table.name
          SCHEDULE_TABLE_NAME      = aws_dynamodb_table.schedule_table.name
          LINE_NOTIFY_ACCESS_TOKEN = local.envs["LINE_NOTIFY_ACCESS_TOKEN"]
        }
      }
      image_identifier      = "${aws_ecr_repository.default.repository_url}:latest"
      image_repository_type = "ECR"
    }
  }
  instance_configuration {
    cpu               = 1024
    memory            = 2048
    instance_role_arn = aws_iam_role.instance_role.arn
  }
  health_check_configuration {
    protocol = "HTTP"
    path     = "/check"
    timeout  = 20
    interval = 20
  }
  depends_on = [time_sleep.wait]
}
