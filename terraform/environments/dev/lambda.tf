module "dev_lambda_delete_event" {
    source = "../../modules/lambda"
    function_name = "dev_delete_event"
    
    file_path = "../../../build/package/delete_event.zip"

    attach_policy_arns_map = {
        "dynamodb_event_table_crud" : aws_iam_policy.dynamodb_event_crud_policy.arn
    }
}