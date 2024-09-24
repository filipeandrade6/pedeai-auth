provider "aws" {
    region = "us-east-1"
}

data "aws_subnets" "cluster_subnets" {
    filter {
        name = "tag:Name"
        values = ["fiap44-private-us-east-1a", "fiap44-private-us-east-1b"]
    }
}

data "aws_vpc" "vpc" {
  filter {
    name   = "tag:Name"
    values = ["fiap44-vpc"]
  }
}

resource "aws_default_security_group" "default_security_group" {
  vpc_id = data.aws_vpc.vpc.id

    ingress {
        protocol  = -1
        self      = true
        from_port = 0
        to_port   = 0
    }

    egress {
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }

    tags = {
        Name = "fiap44-default-security-group"
    }
}

resource "aws_lambda_function" "terraform_lambda_func" {
    filename                       = "${path.module}/deployment.zip"
    function_name                  = "fiap44_auth"
    role                           = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/LabRole"
    handler                        = "handleRequest"
    runtime                        = "provided.al2023"

    vpc_config {
      subnet_ids = data.aws_subnets.cluster_subnets.ids
      security_group_ids = [ aws_default_security_group.default_security_group.id ]
    }

    environment {
        variables = {
            DB_USER = "pedeai"
            DB_PASS = "senha1ABC"
            DB_HOST = "pedeai-db.cfdtpga50euz.us-east-1.rds.amazonaws.com"
            DB_PORT = "5432"
            DB_NAME = "pedeai"
        }
    }
}

data "aws_caller_identity" "current" {}