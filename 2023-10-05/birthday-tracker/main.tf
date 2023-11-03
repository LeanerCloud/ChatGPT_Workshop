provider "aws" {
  region = "eu-west-1"
}

locals {
  lambda_name                 = "BirthdayTrackerFunction"
  lambda_execution_role_name  = "BirthdayLambdaExecutionRole"
  dynamo_birthday_table_name  = "BirthdayTable"
  dynamo_usergroup_table_name = "UserGroupTable"
  cloudfront_origin_id        = "BirthdayLambdaOrigin"
}

data "archive_file" "birthday_lambda_zip" {
  type        = "zip"
  source_dir  = "${path.module}/lambda"
  output_path = "${path.module}/birthday_lambda.zip"
}

resource "aws_iam_role" "birthday_lambda_exec_role" {
  name = local.lambda_execution_role_name

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

module "birthday_lambda" {
  source = "terraform-aws-modules/lambda/aws"

  function_name = "BirthdayTrackerFunction"
  description   = "Lambda function for Birthday Tracker"

  # Note: Adjust architecture, timeout, memory size, and runtime as needed.
  architectures = ["x86_64"]
  timeout       = 300
  memory_size   = 1024

  runtime     = "go1.x"
  source_path = "${path.module}/lambda"
  handler     = "main"

  # If you have any environment variables, define them here.
  # environment_variables = {}

  attach_policy_statements = true
  policy_statements = {
    dynamodb_full_access = {
      effect    = "Allow",
      actions   = ["dynamodb:*"],
      resources = ["*"] # This grants full DynamoDB access. Adjust as needed.
    }
  }
}

resource "aws_lambda_function_url" "birthday_lambda_url" {
  function_name      = module.birthday_lambda.lambda_function_name
  authorization_type = "NONE"
}


module "lambda_url" {
  source = "matti/urlparse/external"

  url = module.birthday_lambda.lambda_function_arn
}

resource "aws_dynamodb_table" "birthday_table" {
  name         = local.dynamo_birthday_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserId"
  range_key    = "BirthdayDate"

  attribute {
    name = "UserId"
    type = "S"
  }

  attribute {
    name = "BirthdayDate"
    type = "S"
  }
}

resource "aws_dynamodb_table" "user_group_table" {
  name         = local.dynamo_usergroup_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "UserId"

  attribute {
    name = "UserId"
    type = "S"
  }
}

resource "aws_iam_policy" "birthday_lambda_dynamodb_policy" {
  name        = "BirthdayLambdaDynamoDBFullAccess"
  description = "Policy granting Birthday Lambda function full access to DynamoDB"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action   = ["dynamodb:*"],
        Effect   = "Allow",
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "birthday_lambda_dynamodb_attachment" {
  policy_arn = aws_iam_policy.birthday_lambda_dynamodb_policy.arn
  role       = aws_iam_role.birthday_lambda_exec_role.name
}

module "cloudfront" {
  source  = "terraform-aws-modules/cloudfront/aws"
  version = "3.0.1"

  comment             = "CloudFront distribution for API and Frontend"
  enabled             = true
  is_ipv6_enabled     = true
  price_class         = "PriceClass_All"
  retain_on_delete    = false
  wait_for_deployment = false
  http_version        = "http3"

  create_monitoring_subscription = false
  create_origin_access_identity  = true

  origin = {
    lambda = {
      domain_name = module.lambda_url.host
      custom_origin_config = {
        http_port              = 80
        https_port             = 443
        origin_protocol_policy = "https-only"
        origin_ssl_protocols   = ["TLSv1", "TLSv1.1", "TLSv1.2"]
      }
    }
    s3_frontend = {
      domain_name = aws_s3_bucket.frontend_bucket.bucket_regional_domain_name
      s3_origin_config = {
        origin_access_identity = module.cloudfront.origin_access_identity_path
      }
    }
  }

  ordered_cache_behavior = [
    {
      path_pattern           = "/api/*"
      target_origin_id       = "lambda"
      viewer_protocol_policy = "redirect-to-https"
      allowed_methods        = ["HEAD", "DELETE", "POST", "GET", "OPTIONS", "PUT", "PATCH"]
      cached_methods         = ["GET", "HEAD"]
      compress               = true
      use_forwarded_values   = true
      query_string           = false
      headers                = ["Authorization"]
      cookies = {
        forward = "none"
      }
    },
    {
      path_pattern           = "/*"
      target_origin_id       = "s3_frontend"
      viewer_protocol_policy = "redirect-to-https"
      allowed_methods        = ["GET", "HEAD"]
      cached_methods         = ["GET", "HEAD"]
      compress               = true
      use_forwarded_values   = false
      query_string           = false
    }
  ]

  default_cache_behavior = {
    target_origin_id       = "lambda" # Or set this to "s3_frontend" if you want S3 to be the default origin
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD", "OPTIONS"]
    cached_methods         = ["GET", "HEAD"]
    compress               = true
    query_string           = true
  }

  depends_on = [aws_lambda_function_url.birthday_lambda_url]
}

resource "aws_s3_bucket" "frontend_bucket" {
  bucket = "birthday-tracker-s3-bucket"
}


