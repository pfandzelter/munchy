provider "aws" {
  profile = var.aws_cli_profile
  region  = var.aws_region
}

resource "aws_lambda_function" "munchy" {
  function_name    = "munchy"
  filename         = "munchy.zip"
  handler          = "munchy"
  source_code_hash = filebase64sha256("munchy.zip")
  role             = aws_iam_role.munchy-role.arn
  runtime          = "go1.x"
  memory_size      = var.lambda_memory_size
  timeout          = var.lambda_timeout

  environment {
    variables = {
      WEBHOOK_URL       = var.webhookurl,
      DYNAMODB_TABLE    = var.table_name,
      DYNAMODB_REGION   = var.aws_region,
      MENSA_TIMEZONE    = var.mensa_timezone,
      DEEPL_TARGET_LANG = var.deepl_target_lang,
      DEEPL_URL         = var.deepl_url,
      DEEPL_KEY         = var.deepl_key,
    }
  }
}

resource "aws_iam_role" "munchy-role" {
  name               = var.project_name
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
  name        = "munchy-lambda_logging"
  path        = "/"
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

resource "aws_iam_policy" "munchy-dynamo" {
  name        = "munchy-dynamo"
  path        = "/"
  description = "IAM policy for DynamoDB access from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "Stmt1582485790003",
      "Action": [
        "dynamodb:Scan"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:dynamodb:*:*:*"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "munchy-lambda_logs" {
  role       = aws_iam_role.munchy-role.name
  policy_arn = aws_iam_policy.munchy-lambda_logging.arn
}

resource "aws_iam_role_policy_attachment" "munchy-dynamo" {
  role       = aws_iam_role.munchy-role.name
  policy_arn = aws_iam_policy.munchy-dynamo.arn
}

resource "aws_cloudwatch_event_rule" "munchy-cron" {
  name                = "munchy-cron"
  schedule_expression = "cron(0 11 ? * 2-6 *)"
}

resource "aws_cloudwatch_event_rule" "munchy-cron-dst" {
  name                = "munchy-cron-dst"
  schedule_expression = "cron(0 10 ? * 2-6 *)"
}

resource "aws_cloudwatch_event_target" "munchy-lambda" {
  target_id = "runLambda"
  rule      = aws_cloudwatch_event_rule.munchy-cron.name
  arn       = aws_lambda_function.munchy.arn
}

resource "aws_cloudwatch_event_target" "munchy-lambda-dst" {
  target_id = "runLambda"
  rule      = aws_cloudwatch_event_rule.munchy-cron-dst.name
  arn       = aws_lambda_function.munchy.arn
}

resource "aws_lambda_permission" "munchy-cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.munchy.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.munchy-cron.arn
}

resource "aws_lambda_permission" "munchy-cloudwatch-dst" {
  statement_id  = "AllowExecutionFromCloudWatchDST"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.munchy.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.munchy-cron-dst.arn
}
