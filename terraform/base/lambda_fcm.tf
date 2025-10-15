module "send_fcm" {
  source = "../modules/lambda"


  file_path     = "${var.root_path}/build/package/send_fcm.zip"
  function_name = "send_fcm"

  env_vars = {
    GOOGLE_APPLICATION_CREDENTIALS = var.google_application_credentials
  }

  attach_policy_arns_map = {
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
  }
}

module "send_anniversary_fcm" {
  source = "../modules/lambda"

  file_path = "${var.root_path}/build/package/send_anniversary_fcm.zip"
  function_name = "send_anniversary_fcm"

  env_vars = {
    GOOGLE_APPLICATION_CREDENTIALS = var.google_application_credentials
  }

  attach_policy_arns_map = {
    "dev_member_table_crud" = data.terraform_remote_state.dev.outputs.dynamodb_dev_member_crud_policy_arn
    "prod_member_table_crud" = data.terraform_remote_state.prod.outputs.dynamodb_prod_member_crud_policy_arn
  }
}