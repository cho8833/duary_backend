resource "aws_apigatewayv2_api" "duary_apigw" {
    api_key_selection_expression = "$request.header.x-api-key"
    disable_execute_api_endpoint = "false"
    name                         = "duary"
    protocol_type                = "HTTP"
    region                       = "ap-northeast-2"
    route_selection_expression   = "$request.method $request.path"
}

resource "aws_apigatewayv2_api" "duary_ws_apigw" {
    api_key_selection_expression = "$request.header.x-api-key"
    disable_execute_api_endpoint = "false"
    ip_address_type              = "ipv4"
    name                         = "duary-ws"
    protocol_type                = "WEBSOCKET"
    region                       = "ap-northeast-2"
    route_selection_expression   = "$request.body.action"
}