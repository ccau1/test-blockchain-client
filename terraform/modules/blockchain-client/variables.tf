variable "env" {
  type = string
  default = "dev"
}

variable "az_count" {
  type = number
  default = 2
}

variable "app_count" {
  type    = number
  default = 1
}

variable "network_mode" {
  type        = string
  description = "network mode"
  default     = "awsvpc"
}

variable "launch_type" {
  description = "launch type"
  type = object({
    type   = string
    cpu    = number
    memory = number
    container_port = number
  })
  default = {
    type   = "FARGATE"
    cpu    = 1024
    memory = 2048
    container_port = 3000
  }
}