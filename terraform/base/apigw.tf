resource "aws_apigatewayv2_api" "duary_apigw" {
    api_key_selection_expression = "$request.header.x-api-key"
    disable_execute_api_endpoint = "false"
    name                         = "duary"
    protocol_type                = "HTTP"
    route_selection_expression   = "$request.method $request.path"
}

resource "aws_apigatewayv2_api" "duary_ws_apigw" {
    api_key_selection_expression = "$request.header.x-api-key"
    disable_execute_api_endpoint = "false"
    name                         = "duary-ws"
    protocol_type                = "WEBSOCKET"
    route_selection_expression   = "$request.body.action"
}

resource "aws_apigatewayv2_stage" "duary_apigw_dev_stage" {
    api_id = aws_apigatewayv2_api.duary_apigw.id
    name   = "dev"
    stage_variables = {
        stage = "dev"
    }
    auto_deploy = true
}

resource "aws_apigatewayv2_stage" "duary_apigw_prod_stage" {
    api_id = aws_apigatewayv2_api.duary_apigw.id
    name = "prod"

    stage_variables = {
        stage = "prod"
    }
}

resource "aws_apigatewayv2_stage" "duary_ws_apigw_dev_stage" {
    api_id = aws_apigatewayv2_api.duary_ws_apigw.id
    name = "dev"

    stage_variables = {
        stage = "dev"
    }
    default_route_settings {
        throttling_burst_limit = 500
        throttling_rate_limit = 1000
    }
    auto_deploy = true
}

resource "aws_apigatewayv2_stage" "duary_ws_apigw_prod_stage" {
    api_id = aws_apigatewayv2_api.duary_ws_apigw.id
    name = "prod"

    default_route_settings {
        throttling_burst_limit = 500
        throttling_rate_limit = 1000
    }

    stage_variables = {
        stage = "prod"
    }
}