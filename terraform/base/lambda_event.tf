module "delete_event_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "delete_event_api"
  file_path     = "${var.root_path}/build/package/delete_event.zip"

  api_type      = "HTTP"
  apigw_id      = aws_apigatewayv2_api.duary_apigw.id
  http_method = "DELETE"
  http_path = "/event"

  env_vars = {
    DEV_WEBSOCKET_URL = local.dev_websocket_post_url
    PROD_WEBSOCKET_URL = local.prod_websocket_post_url
  }
  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "dev_ws_connection_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_ws_connection_crud_policy_arn
    "prod_ws_connection_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_ws_connection_crud_policy_arn
    "dev_event_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "execute_ws_api" = aws_iam_policy.allow_execute_ws_api.arn
  }
}

module "create_event_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "create_event_api"
  file_path = "${var.root_path}/build/package/create_event.zip"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/event"

  env_vars = {
    DEV_WEBSOCKET_URL = local.dev_websocket_post_url
    PROD_WEBSOCKET_URL = local.prod_websocket_post_url
  }

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "dev_ws_connection_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_ws_connection_crud_policy_arn
    "prod_ws_connection_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_ws_connection_crud_policy_arn
    "dev_event_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "execute_ws_api" = aws_iam_policy.allow_execute_ws_api.arn
    "execute_send_fcm" = aws_iam_policy.allow_execute_send_fcm.arn
  }
}

module "update_event_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "update_event_api"
  file_path = "${var.root_path}/build/package/update_event.zip"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "PUT"
  http_path = "/event"

  env_vars = {
    DEV_WEBSOCKET_URL = local.dev_websocket_post_url
    PROD_WEBSOCKET_URL = local.prod_websocket_post_url
  }

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "dev_ws_connection_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_ws_connection_crud_policy_arn
    "prod_ws_connection_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_ws_connection_crud_policy_arn
    "dev_event_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "execute_ws_api" = aws_iam_policy.allow_execute_ws_api.arn
  }
}

module "get_event_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "get_event_api"
  file_path = "${var.root_path}/build/package/get_event.zip"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "GET"
  http_path = "/event"

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_event_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
  }
}
