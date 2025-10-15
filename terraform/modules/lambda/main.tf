data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# lambda Cloud watch Log Policy
data "aws_iam_policy_document" "basic_lambda_exec_doc" {
    statement {
      actions = [
        "logs:CreateLogGroup"
      ]
      effect = "Allow"
      sid = "AllowCreateLogGroup"
      resources = [
        "arn:aws:logs:${data.aws_region.current.region}:${data.aws_caller_identity.current.account_id}:*"
      ]
    }
    statement {
      actions = [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
      effect = "Allow"
      sid = "AllowPutLogEvents"
      resources = [
        "arn:aws:logs:${data.aws_region.current.region}:${data.aws_caller_identity.current.account_id}:log-group:/aws/lambda/${var.function_name}:*"
      ]
    }
    version = "2012-10-17"
}

data "aws_iam_policy_document" "assume_role" {
    statement {
      effect = "Allow"

      principals {
        type = "Service"
        identifiers = ["lambda.amazonaws.com"]
      }
      actions = ["sts:AssumeRole"]
    }
}

resource "aws_iam_role" "this" {
    name = "${var.function_name}-execution-role"
    path = "/service-role/"

    assume_role_policy = data.aws_iam_policy_document.assume_role.json
}



# cloud watch logs
resource "aws_iam_role_policy" "basic_execution" {
    name = "${var.function_name}-basic-execution-policy"
    role = aws_iam_role.this.id
    policy = data.aws_iam_policy_document.basic_lambda_exec_doc.json
}

# attach additional policy param
resource "aws_iam_role_policy_attachment" "additional_policies" {
    for_each = var.attach_policy_arns_map
    role = aws_iam_role.this.name
    policy_arn = each.value
}

# lambda function
resource "aws_lambda_function" "this" {

    logging_config {
      log_format = "Text"
      log_group = "/aws/lambda/${var.function_name}"
    }

    ephemeral_storage {
      size = "512"
    }

    memory_size = "128"

    package_type = "Zip"

    reserved_concurrent_executions = "-1"

    skip_destroy = "false"
    
    architectures = ["x86_64"]

    tracing_config {
      mode = "PassThrough"
    }

    filename = var.file_path

    # 파일 변경 감지 hash
    source_code_hash = filebase64sha256(var.file_path)

    function_name = var.function_name
    
    description = var.description
    
    handler = var.handler
    
    runtime = var.runtime

    timeout = var.timeout

    role = aws_iam_role.this.arn

    depends_on = [ 
        aws_iam_role.this,
        aws_iam_role_policy.basic_execution
    ]

    environment {
        variables = var.env_vars
    }
}

resource "aws_cloudwatch_log_group" "this" {
  name = "/aws/lambda/${var.function_name}"
}