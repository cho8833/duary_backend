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
    authorizer_payload_format_version = "1.0"
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