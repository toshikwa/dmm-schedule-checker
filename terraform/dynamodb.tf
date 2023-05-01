resource "aws_dynamodb_table" "teacher_table" {
  name           = "${local.app}_teacher_table"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "id"
  attribute {
    name = "id"
    type = "S"
  }
}

resource "aws_dynamodb_table" "schedule_table" {
  name           = "${local.app}_schedule_table"
  billing_mode   = "PROVISIONED"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "teacherId"
  range_key      = "dateTime"
  attribute {
    name = "teacherId"
    type = "S"
  }
  attribute {
    name = "dateTime"
    type = "S"
  }
}
