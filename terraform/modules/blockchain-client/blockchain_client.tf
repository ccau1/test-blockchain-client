
# // BACKEND 

resource "aws_ecr_repository" "bc_client_api" {
  name                 = "${var.env}-bc-client-api"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "bc_client_api" {
  repository = aws_ecr_repository.bc_client_api.name

  policy = <<EOF
{
    "rules": [
        {
            "rulePriority": 1,
            "description": "Expire images older than 14 days",
            "selection": {
                "tagStatus": "untagged",
                "countType": "sinceImagePushed",
                "countUnit": "days",
                "countNumber": 14
            },
            "action": {
                "type": "expire"
            }
        }
    ]
}
EOF
}

resource "aws_lb_target_group" "bc_client_api" {
  name        = "${var.env}-bc-client-api-tg"
  port        = var.launch_type.container_port
  protocol    = "HTTP"
  vpc_id      = aws_vpc.default.id
  target_type = "ip"
}

resource "aws_security_group" "bc_client_api_task" {
  name   = "${var.env}-bc-client-api-task-sg"
  vpc_id = aws_vpc.default.id

  ingress {
    protocol        = "tcp"
    from_port       = 80
    to_port         = var.launch_type.container_port
    security_groups = [aws_security_group.lb.id]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_lb_listener" "bc_client_api" {
  load_balancer_arn = aws_lb.default.id
  port              = var.launch_type.container_port
  protocol          = "HTTP"

  default_action {
    target_group_arn = aws_lb_target_group.bc_client_api.id
    type             = "forward"
  }
}

resource "aws_ecs_service" "bc_client_api" {
  name            = "${var.env}-bc-client-api-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.bc_client_api.arn
  desired_count   = var.app_count
  launch_type     = var.launch_type.type

  network_configuration {
    security_groups = [aws_security_group.bc_client_api_task.id]
    subnets         = aws_subnet.private.*.id
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.bc_client_api.id
    container_name   = "${var.env}-bc-client-api"
    container_port   = var.launch_type.container_port
  }

  depends_on = [aws_lb_listener.bc_client_api, aws_iam_role_policy_attachment.ecs_tasks_execution_role]
}

# //

data "aws_iam_policy_document" "ecs_tasks_execution_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_tasks_execution_role" {
  name               = "${var.env}-ecs-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.ecs_tasks_execution_role.json
}

resource "aws_iam_role_policy_attachment" "ecs_tasks_execution_role" {
  role       = aws_iam_role.ecs_tasks_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_task_definition" "bc_client_api" {
  family                   = "${var.env}-bc-client-api"
  network_mode             = var.network_mode
  requires_compatibilities = [var.launch_type.type]
  cpu                      = var.launch_type.cpu
  memory                   = var.launch_type.memory
  execution_role_arn       = aws_iam_role.ecs_tasks_execution_role.arn

  container_definitions = <<DEFINITION
    [
      {
        "image": "029319782524.dkr.ecr.ap-southeast-1.amazonaws.com/bc-client-api:latest",
        "cpu": ${var.launch_type.cpu},
        "memory": ${var.launch_type.memory},
        "name": "${var.env}-bc-client-api",
        "networkMode": "${var.network_mode}",
        "portMappings": [
          {
            "containerPort": ${var.launch_type.container_port},
            "hostPort": ${var.launch_type.container_port}
          }
        ],
        "logConfiguration": {
          "logDriver": "awslogs",
          "options": {
            "awslogs-group": "${var.env}-bc-client-api-container",
            "awslogs-region": "ap-southeast-1",
            "awslogs-create-group": "true",
            "awslogs-stream-prefix": "${var.env}-bc-client-api"
          }
        }
      }
    ]
    DEFINITION
}
