module "update_member_api" {
  source = "../modules/api_lambda_endpoint"

  file_path = "../../build/package/update_member.zip"
  function_name = "update_member_api"

  api_type = "HTTP"
  apigw_id = aws_apigatewayv2_api.duary_apigw.id
  http_method = "POST"
  http_path = "/member"

  authorizer_id = aws_apigatewayv2_authorizer.http_jwt_authorizer.id

  attach_policy_arns_map = {
    "dev_couple_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_couple_crud_policy_arn
    "prod_couple_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_couple_crud_policy_arn
    "dev_ws_connection_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_ws_connection_crud_policy_arn
    "prod_ws_connection_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_ws_connection_crud_policy_arn
    "dev_event_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_event_crud_policy_arn
    "prod_event_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_event_crud_policy_arn
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
    "execute_ws_api" = aws_iam_policy.allow_execute_ws_api.arn
    "event_bridge_scheduler_crud" = aws_iam_policy.allow_crud_event_bridge_scheduler.arn
  }
}