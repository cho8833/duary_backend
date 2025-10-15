module "update_couple_api" {
  source = "../modules/api_lambda_endpoint"

  file_path     = "${var.root_path}/build/package/update_couple.zip"
  function_name = "update_couple_api"

  api_type    = "HTTP"
  apigw_id    = aws_apigatewayv2_api.duary_apigw.id
  http_method = "PUT"
  http_path   = "/couple"

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_couple_table_crud"       = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "prod_couple_table_crud"      = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "dev_event_table_crud"        = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud"       = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "event_bridge_scheduler_crud" = aws_iam_policy.allow_crud_event_bridge_scheduler.arn
  }
}

module "connect_couple_api" {
  source = "../modules/api_lambda_endpoint"

  file_path     = "${var.root_path}/build/package/connect_couple.zip"
  function_name = "connect_couple_api"

  api_type    = "HTTP"
  apigw_id    = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path   = "/couple/connect"

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  env_vars = {
    secretKey = var.jwt_secret_key
    DEV_WEBSOCKET_URL  = local.dev_websocket_post_url
    PROD_WEBSOCKET_URL = local.prod_websocket_post_url
  }

  attach_policy_arns_map = {
    "dev_couple_table_crud"         = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "prod_couple_table_crud"        = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "dev_ws_connection_table_crud"  = data.terraform_remote_state.dev.outputs.dynamodb_dev_ws_connection_crud_policy_arn
    "prod_ws_connection_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_ws_connection_crud_policy_arn
    "dev_event_table_crud"          = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud"         = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "dev_member_table_crud"         = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_member_table_crud"        = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "execute_ws_api"                = aws_iam_policy.allow_execute_ws_api.arn
    "event_bridge_scheduler_crud"   = aws_iam_policy.allow_crud_event_bridge_scheduler.arn

  }
}

module "disconnect_couple_api" {
  source = "../modules/api_lambda_endpoint"

  function_name = "disconnect_couple_api"
  file_path     = "${var.root_path}/build/package/disconnect_couple.zip"

  api_type    = "HTTP"
  apigw_id    = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path   = "/couple/disconnect"

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  env_vars = {
    secretKey = var.jwt_secret_key
  }

  attach_policy_arns_map = {
    "dev_couple_table_crud"       = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "prod_couple_table_crud"      = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "dev_event_table_crud"        = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud"       = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "dev_member_table_crud"       = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_member_table_crud"      = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "event_bridge_scheduler_crud" = aws_iam_policy.allow_crud_event_bridge_scheduler.arn
  }
}