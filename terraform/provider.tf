variable "aws_profile" {}

provider "aws" {
  region = "us-east-1"
  profile = var.aws_profile
}