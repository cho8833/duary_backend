data "aws_caller_identity" "current" {}
data "aws_region" "current" {}


locals {
    route_key = var.api_type == "HTTP" ? "${var.http_method} ${var.http_path}" : var.ws_route_key

    source_arn_path = var.api_type == "HTTP" ? "/*/*${var.http_path}" : "/*${var.ws_route_key}"
}

resource "aws_apigatewayv2_route" "this" {
    api_id = var.apigw_id
    api_key_required = "false"
    route_key = local.route_key
    authorization_type = var.authorizer_id != null ? "CUSTOM" : "NONE"
    authorizer_id = var.authorizer_id
    target = "integrations/${aws_apigatewayv2_integration.this.id}"
}

resource "aws_apigatewayv2_integration" "this" {
    api_id = var.apigw_id
    integration_type = "AWS_PROXY"
    connection_type = "INTERNET"
    integration_method = "POST"
    payload_format_version = "2.0"
    integration_uri = var.lambda_invoke_arn
}

resource "aws_lambda_permission" "this" {
    action = "lambda:InvokeFunction"
    function_name = var.lambda_arn
    principal = "apigateway.amazonaws.com"
    region = "ap-northeast-2"
    source_arn = "arn:aws:execute-api:${data.aws_region.current.region}:${data.aws_caller_identity.current.account_id}:${var.apigw_id}${local.source_arn_path}"
    statement_id = "AllowExecutionFromAPIGateway_${var.function_name}"
}

