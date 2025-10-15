module "ws_connect" {
  source = "../modules/api_lambda_endpoint"

  api_type      = "WEBSOCKET"
  apigw_id      = aws_apigatewayv2_api.duary_ws_apigw.id

  file_path     = "${var.root_path}/build/package/ws_connect.zip"
  function_name = "ws_connect_api"
  ws_route_key = "$connect"

  authorizer_id = aws_apigatewayv2_authorizer.ws_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_ws_connection_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_ws_connection_crud_policy_arn
    "prod_ws_connection_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_ws_connection_crud_policy_arn
  }
}

module "ws_disconnect" {
  source = "../modules/api_lambda_endpoint"

  function_name = "ws_disconnect_api"
  file_path = "${var.root_path}/build/package/ws_disconnect.zip"
  ws_route_key = "$disconnect"

  api_type = "WEBSOCKET"
  apigw_id = aws_apigatewayv2_api.duary_ws_apigw.id

  attach_policy_arns_map = {
    "dev_ws_connection_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_ws_connection_crud_policy_arn
    "prod_ws_connection_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_ws_connection_crud_policy_arn
  }
}

