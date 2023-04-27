data "aws_availability_zones" "available_zones" {
  state = "available"
}

# // 

resource "aws_iam_group" "deploy" {
  name = "${var.env}-deploy-group"
}

resource "aws_iam_policy" "deploy_policy" {
  name        = "${var.env}-deploy-policy"
  description = "A deploy policy"

  tags = {
    Name = "${var.env}-iam-policy-deploy"
  }

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "imagebuilder:GetComponent",
        "imagebuilder:GetContainerRecipe",
        "logs:CreateLogGroup",
        "iam:PassRole",
        "ecs:DescribeServices",
        "ecs:RegisterTaskDefinition",
        "ecs:UpdateService",
        "ecr:GetAuthorizationToken",
        "ecr:BatchGetImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:PutImage"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "kms:Decrypt"
      ],
      "Resource": "*",
      "Condition": {
        "ForAnyValue:StringEquals": {
          "kms:EncryptionContextKeys": "aws:imagebuilder:arn",
          "aws:CalledVia": [
            "imagebuilder.amazonaws.com"
          ]
        }
      }
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::ec2imagebuilder*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:CreateLogGroup",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:log-group:/aws/imagebuilder/*"
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "deploy_policy_attachment" {
  name       = "${var.env}-deploy-policy-attachment"
  groups     = [aws_iam_group.deploy.name]
  policy_arn = aws_iam_policy.deploy_policy.arn
}

resource "aws_iam_user" "deploy" {
  name = "${var.env}-deploy"
  path = "/system/"
}

resource "aws_iam_user_login_profile" "deploy" {
  user                    = "${aws_iam_user.deploy.name}"
  password_reset_required = true
  pgp_key="keybase:terraform_user"
}

resource "aws_iam_access_key" "deploy" {
  user = aws_iam_user.deploy.name
}

resource "aws_iam_group_membership" "deploy" {
  name = "${var.env}-group-membership-deploy"

  users = [
    aws_iam_user.deploy.name,
  ]

  group = aws_iam_group.deploy.name
}

# // setup secrets manager

resource "aws_secretsmanager_secret" "deploy_access_key" {
  name = "${var.env}-secrets"
}

resource "aws_secretsmanager_secret_version" "deploy_access_key" {
  secret_id     = aws_secretsmanager_secret.deploy_access_key.id
  # define secret object
  secret_string = jsonencode({})

  depends_on = [
    aws_iam_access_key.deploy
  ]
}

# //

# resource "aws_s3_bucket" "main_bucket" {
#   bucket = "main-bucket"

#   tags = {
#     Name = "${var.env}-bucket"
#   }
# }

# //

resource "aws_vpc" "default" {
  cidr_block           = "10.32.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "${var.env}-vpc-default"
  }
}

# // 

resource "aws_subnet" "public" {
  count                   = var.az_count
  vpc_id                  = aws_vpc.default.id
  cidr_block              = cidrsubnet(aws_vpc.default.cidr_block, 8, var.az_count + count.index)
  map_public_ip_on_launch = true
  availability_zone       = data.aws_availability_zones.available_zones.names[count.index]

  tags = {
    Name = "${var.env}-public-subnet"
  }
}

resource "aws_subnet" "private" {
  count             = var.az_count
  vpc_id            = aws_vpc.default.id
  cidr_block        = cidrsubnet(aws_vpc.default.cidr_block, 8, count.index)
  availability_zone = data.aws_availability_zones.available_zones.names[count.index]

  tags = {
    Name = "${var.env}-private-subnet"
  }
}

# // 

resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.default.id

  tags = {
    Name = "${var.env}-igw"
  }
}

# // internet_access
resource "aws_route" "internet_access" {
  route_table_id         = aws_vpc.default.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.gateway.id
}

resource "aws_eip" "gateway" {
  count      = 2
  vpc        = true
  depends_on = [aws_internet_gateway.gateway]

  tags = {
    Name = "${var.env}-eip"
  }
}

resource "aws_nat_gateway" "gateway" {
  count         = 2
  subnet_id     = element(aws_subnet.public.*.id, count.index)
  allocation_id = element(aws_eip.gateway.*.id, count.index)

  tags = {
    Name = "${var.env}-nat-gateway"
  }
}

resource "aws_route_table" "private" {
  count  = 2
  vpc_id = aws_vpc.default.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = element(aws_nat_gateway.gateway.*.id, count.index)
  }

  tags = {
    Name = "${var.env}-route-table-private"
  }
}

resource "aws_route_table_association" "private" {
  count          = 2
  subnet_id      = element(aws_subnet.private.*.id, count.index)
  route_table_id = element(aws_route_table.private.*.id, count.index)
}

# // load balancer in public subnet

resource "aws_security_group" "lb" {
  name   = "${var.env}-main-alb-security-group"
  vpc_id = aws_vpc.default.id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "${var.env}-sg-lb"
  }
}

resource "aws_lb" "default" {
  name            = "${var.env}-default-lb"
  subnets         = aws_subnet.public.*.id
  security_groups = [aws_security_group.lb.id]

  tags = {
    Name = "${var.env}-lb-default"
  }
}

# // ECS cluster

resource "aws_ecs_cluster" "main" {
  name = "${var.env}-main-cluster"

  tags = {
    Name = "${var.env}-cluster-main"
  }
}
