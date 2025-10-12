output "function_arn" {
    description = "arn of lambda function"
    value = aws_lambda_function.this.arn
}

output "function_name" {
    description = "name of lambda function"
    value = aws_lambda_function.this.function_name
}

output "function_invoke_arn" {
    description = "invoke arn of lambda function"
    value = aws_lambda_function.this.invoke_arn
}

output "function_id" {
    description = "ID of lambda function"
    value = aws_lambda_function.this.id
}