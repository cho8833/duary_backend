variable "function_name" {
    description = "lambda function name"
    type = string
}

variable "memory_size" {
    description = "lambda memory size"
    type = string
    default = "128"
}

variable "timeout" {
    description = "lambda call time out"
    type = string
    default = "10"
}

variable "handler" {
    description = "lambda handler"
    type = string
    default = "hello.handler"
}

variable "env_vars" {
    description = "lambda environment variables"
    type = map(string)
    default = {}
}

variable "file_path" {
    description = "lambda zip file path"
    type = string
}

variable "attach_policy_arns_map" {
    description = "lambda additinal policy arns, key is logical name, value is arn"
    type = map(string)
    default = {}
}

variable "description" {
    description = "lambda function description"
    type = string
    default = ""
}

variable "runtime" {
    description = "lambda runtime"
    type = string
    default = "provided.al2023"
}