output "dynamodb_dev_event_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_dev_event_crud_policy.arn
}

output "dynamodb_dev_couple_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_dev_couple_crud_policy.arn
}

output "dynamodb_dev_member_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_dev_member_crud_policy.arn
}

output "dynamodb_dev_ws_connection_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_dev_ws_connection_crud_policy.arn
}