output "http_apigw_id" {
    description = "ID of HTTP api gateway"
    value = aws_apigatewayv2_api.duary_apigw.id
}

output "http_jwt_authorizer_id" {
    description = "ID of HTTP jwt authorizer"
    value = module.http_jwt_authorizer.function_id
}

output "ws_apigw_id" {
    description = "ID of WEBSOCKET api gateway"
    value = aws_apigatewayv2_api.duary_ws_apigw.id
}

output "ws_jwt_authorizer_id" {
    description = "ID of WEBSOCKET jwt authorizer"
    value = aws_apigatewayv2_api.duary_ws_apigw.id
}