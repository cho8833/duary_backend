output "dynamodb_prod_event_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_prod_event_crud_policy.arn
}

output "dynamodb_prod_cert_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_prod_cert_crud_policy.arn
}

output "dynamodb_prod_member_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_prod_member_crud_policy.arn
}

output "dynamodb_prod_couple_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_prod_couple_crud_policy.arn
}

output "dynamodb_prod_ws_connection_crud_policy_arn" {
  value = aws_iam_policy.dynamodb_prod_ws_connection_crud_policy.arn
}