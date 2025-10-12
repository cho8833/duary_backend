module "http_jwt_authorizer" {
    source = "../modules/lambda"
    file_path = "../../build/package/jwt_authorizer.zip"
    function_name = "http_jwt_authorizer"

    env_vars = {
        "secretKey" = var.jwt_secret_key
    }
}

module "ws_jwt_authorizer" {
    source = "../modules/lambda"
    file_path = "../../build/package/ws_authorizer.zip"
    function_name = "ws_jwt_authorizer"

    env_vars = {
        "secretKey" = var.jwt_secret_key
    }
}

resource "aws_apigatewayv2_authorizer" "http_jwt_authorizer" {
    api_id = aws_apigatewayv2_api.duary_apigw.id
    name = "http_jwt_authorizer"
    authorizer_type = "REQUEST"
    authorizer_uri = module.http_jwt_authorizer.function_invoke_arn
    authorizer_payload_format_version = "2.0"
    authorizer_result_ttl_in_seconds = "3600"
    identity_sources = ["$request.header.Authorization"]
    region = "ap-northeast-2"
}

resource "aws_apigatewayv2_authorizer" "ws_jwt_authorizer" {
    api_id = aws_apigatewayv2_api.duary_ws_apigw.id
    name = "ws_jwt_authorizer"
    authorizer_type = "REQUEST"
    authorizer_uri = module.ws_jwt_authorizer.function_invoke_arn
    identity_sources = ["route.request.querystring.Authorization"]
    region = "ap-northeast-2"
}

resource "aws_lambda_permission" "allow_http_api_invoke_authorizer" {
    statement_id = "AllowInvokeFromHTTPAPIGatewayAuthorizer"
    action = "lambda:InvokeFunction"
    function_name = module.http_jwt_authorizer.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_apigatewayv2_api.duary_apigw.execution_arn}/authorizers/${aws_apigatewayv2_authorizer.http_jwt_authorizer.id}"
}

resource "aws_lambda_permission" "allow_ws_api_invoke_authorizer" {
    statement_id = "AllowInvokeFromWSAPIGatewayAuthorizer"
    action = "lambda:InvokeFunction"
    function_name = module.ws_jwt_authorizer.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_apigatewayv2_api.duary_ws_apigw.execution_arn}/authorizers/${aws_apigatewayv2_authorizer.ws_jwt_authorizer.id}"
}