// public subnet
resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.vpc.id
  availability_zone       = var.public_subnet_region
  cidr_block              = var.public_subnet_cidr
  map_public_ip_on_launch = true
  tags = {
    "Name" = "${var.environment_name}-public-subnet"
  }
}

resource "aws_route_table_association" "subnet_rtb_association_pub" {
  route_table_id = aws_route_table.public.id
  subnet_id      = aws_subnet.public.id
}

// private subnet
resource "aws_subnet" "private_1" {
  vpc_id                  = aws_vpc.vpc.id
  availability_zone       = var.private_subnet_region_1
  cidr_block              = var.private_subnet_cidr_1
  map_public_ip_on_launch = false
  tags = {
    "Name" = "${var.environment_name}-private-subnet-1"
  }
}

resource "aws_route_table_association" "subnet_rtb_association_prv_1" {
  route_table_id = aws_route_table.private.id
  subnet_id      = aws_subnet.private_1.id
}

resource "aws_subnet" "private_2" {
  vpc_id                  = aws_vpc.vpc.id
  availability_zone       = var.private_subnet_region_2
  cidr_block              = var.private_subnet_cidr_2
  map_public_ip_on_launch = false
  tags = {
    "Name" = "${var.environment_name}-private-subnet-2"
  }
}

resource "aws_route_table_association" "subnet_rtb_association_prv_2" {
  route_table_id = aws_route_table.private.id
  subnet_id      = aws_subnet.private_2.id
}
