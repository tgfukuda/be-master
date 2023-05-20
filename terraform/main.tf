data "aws_key_pair" "ssh" {
  key_pair_id = var.keypair_id
}

resource "aws_instance" "public" {
  ami                         = var.ami
  instance_type               = var.instance_type
  availability_zone = var.public_subnet_region
  key_name                    = data.aws_key_pair.ssh.key_name
  ebs_optimized               = false
  monitoring                  = false
  vpc_security_group_ids = [ aws_security_group.no_limit.id ]
  associate_public_ip_address = true
  subnet_id                   = aws_subnet.public.id
  iam_instance_profile = aws_iam_instance_profile.oracle_instance_profile.id
  ebs_block_device {
    device_name           = "/dev/xvdh"
    volume_type           = "gp3"
    volume_size           = var.volume_size
    encrypted             = false
    delete_on_termination = true
  }
  tags = {
    "Name" = "${var.environment_name}-public-oracle"
  }
}

resource "aws_instance" "oracle_1" {
  ami                         = var.ami
  instance_type               = var.instance_type
  availability_zone = var.private_subnet_region_1
  key_name                    = data.aws_key_pair.ssh.key_name
  ebs_optimized               = false
  monitoring                  = false
  vpc_security_group_ids = [ aws_security_group.private.id ]
  associate_public_ip_address = false
  subnet_id                   = aws_subnet.private_1.id
  iam_instance_profile = aws_iam_instance_profile.oracle_instance_profile.id
  ebs_block_device {
    device_name           = "/dev/xvdh"
    volume_type           = "gp3"
    volume_size           = var.volume_size
    encrypted             = false
    delete_on_termination = true
  }
  tags = {
    "Name" = "${var.environment_name}-private-oracle-1"
  }
}

resource "aws_instance" "oracle_2" {
  ami                         = var.ami
  instance_type               = var.instance_type
  availability_zone = var.private_subnet_region_2
  key_name                    = data.aws_key_pair.ssh.key_name
  ebs_optimized               = false
  monitoring                  = false
  vpc_security_group_ids = [ aws_security_group.private.id ]
  associate_public_ip_address = false
  subnet_id                   = aws_subnet.private_2.id
  iam_instance_profile = aws_iam_instance_profile.oracle_instance_profile.id
  ebs_block_device {
    device_name           = "/dev/xvdh"
    volume_type           = "gp3"
    volume_size           = var.volume_size
    encrypted             = false
    delete_on_termination = true
  }
  tags = {
    "Name" = "${var.environment_name}-private-oracle-2"
  }
}

resource "aws_s3_bucket" "log" {
  bucket_prefix = "${var.environment_name}-oracle-logs"
  
  # lifecycle {
  #   id = "log"
  #   enabled = true
  #   prefix= "logs/"

  #   expiration = {
  #     days = 30
  #   }
  # }

  tags = {
    "Name" = "${var.environment_name}-oracle-logs"
  }
}

# resource "aws_s3_bucket_acl" "log_acl" {
#   bucket = aws_s3_bucket.log.id
#   acl = "private"
# }

# resource "aws_s3_bucket_object" "endpoint_public" {
#   bucket = aws_s3_bucket.log.id
#   key = aws_instance.public.private_ip
#   content = aws_instance.public.tags["Name"]
# }

# resource "aws_s3_bucket_object" "endpoint_private_1" {
#   bucket = aws_s3_bucket.log.id
#   key = aws_instance.oracle_1.private_ip
#   content = aws_instance.oracle_1.tags["Name"]
# }

# resource "aws_s3_bucket_object" "endpoint_private_2" {
#   bucket = aws_s3_bucket.log.id
#   key = aws_instance.oracle_2.private_ip
#   content = aws_instance.oracle_2.tags["Name"]
# }
