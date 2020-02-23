data "aws_region" "current" {}

variable "project-name" {
  type = string
  default = "munchy"
}

variable "table-name" {
  type = string
  description = "Name of DynamoDB table where food items are stored."
  default = "go-eat"
}

variable "webhookurl" {
  description = "Slack webhook url to post to."
  type = string
}


resource "aws_lambda_function" "munchy" {
  function_name    = "munchy"
  filename         = "munchy.zip"
  handler          = "munchy"
  source_code_hash = filebase64sha256("munchy.zip")
  role             = aws_iam_role.munchy-role.arn
  runtime          = "go1.x"
  memory_size      = 128
  timeout          = 1

  environment {
    variables = {
      WEBHOOK_URL = var.webhookurl,
      DYNAMODB_TABLE = var.table-name,
      DYNAMODB_REGION = data.aws_region.current.name
    }
  }
}

resource "aws_iam_role" "munchy-role" {
  name               = var.project-name
  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": {
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    },
    "Effect": "Allow"
  }
}
POLICY
}

resource "aws_iam_role_policy_attachment" "munchy-basic-exec-role" {
  role       = aws_iam_role.munchy-role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "munchy-lambda_logging" {
  name = "munchy-lambda_logging"
  path = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "munchy-lambda_logs" {
  role = aws_iam_role.munchy-role.name
  policy_arn = aws_iam_policy.munchy-lambda_logging.arn
}

resource "aws_cloudwatch_event_rule" "munchy-cron" {
  name                = "munchy-cron"
  schedule_expression = "cron(0 11 ? * 2-6 *)"
}

resource "aws_cloudwatch_event_target" "munchy-lambda" {
  target_id = "runLambda"
  rule      = aws_cloudwatch_event_rule.munchy-cron.name
  arn       = aws_lambda_function.munchy.arn
}

resource "aws_lambda_permission" "munchy-cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.munchy.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.munchy-cron.arn
}
