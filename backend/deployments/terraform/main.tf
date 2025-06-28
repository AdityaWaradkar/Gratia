# Get default VPC and subnets
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "default-for-az"
    values = ["true"]
  }
}

# EKS Module
module "eks" {
  source          = "terraform-aws-modules/eks/aws"
  version         = "20.8.4"

  cluster_name    = "gratia-cluster"
  cluster_version = "1.29"

  subnet_ids      = data.aws_subnets.default.ids
  vpc_id          = data.aws_vpc.default.id

  enable_irsa     = true

  eks_managed_node_group_defaults = {
    instance_types = ["t3.small"]
  }

  eks_managed_node_groups = {
    default = {
      desired_size   = 1
      max_size       = 2
      min_size       = 1
      instance_types = ["t3.small"]
      capacity_type  = "SPOT"
    }
  }

  tags = {
    environment = "dev"
    project     = "gratia"
  }
}