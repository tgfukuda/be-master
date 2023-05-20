resource "aws_vpc" "vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_support   = true
  enable_dns_hostnames = true
  tags = {
    "Name" = "${var.environment_name}-vpc"
  }
}

// public
resource "aws_internet_gateway" "internet_gateway" {
  tags = {
    "Name" = "${var.environment_name}-igw"
  }
}

resource "aws_internet_gateway_attachment" "internet_gateway_attachment" {
  internet_gateway_id = aws_internet_gateway.internet_gateway.id
  vpc_id              = aws_vpc.vpc.id
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    "Name" = "${var.environment_name}-pub-rtb"
  }
}

resource "aws_route" "public" {
  route_table_id         = aws_route_table.public.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.internet_gateway.id
}

// private
resource "aws_eip" "natgateway_eip" {
  depends_on = [
    aws_internet_gateway_attachment.internet_gateway_attachment
  ]
}

resource "aws_nat_gateway" "natgateway" {
  subnet_id     = aws_subnet.public.id
  allocation_id = aws_eip.natgateway_eip.allocation_id
}

resource "aws_route_table" "private" {
  vpc_id = aws_vpc.vpc.id
  tags = {
    "Name" = "${var.environment_name}-prv-rtb"
  }
}

resource "aws_route" "private" {
  route_table_id         = aws_route_table.private.id
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = aws_nat_gateway.natgateway.id
}
