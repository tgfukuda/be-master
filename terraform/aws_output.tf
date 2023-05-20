output "vpc" {
  description = "A reference to the created VPC"
  value       = aws_vpc.vpc.id
}

output "public_subnet" {
  description = "A list of the public subnets"
  value       = aws_subnet.public.id
}

output "private_subnets" {
  description = "A list of the private subnets"
  value       = join(",", [aws_subnet.private_1.id, aws_subnet.private_2.id])
}

output "sg_public" {
  description = "Security group with no limit"
  value       = aws_security_group.no_limit.id
}

output "sg_private" {
  description = "Security group only connected with public instance"
  value       = aws_security_group.private.id
}

output "attached_role" {
  value = aws_iam_role.oracle_role.arn
}

output "public_instance" {
  value = aws_instance.public.id
}

output "oracle_1" {
  value = aws_instance.oracle_1.id
}

output "oracle_2" {
  value = aws_instance.oracle_2.id
}

output "log_bucket" {
  value = aws_s3_bucket.log.id
}
