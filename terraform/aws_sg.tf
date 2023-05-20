resource "aws_security_group" "no_limit" {
  name        = "public"
  description = "public sg of the vpc"
  vpc_id      = aws_vpc.vpc.id
  egress {
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    to_port     = 0
  }
  ingress {
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    to_port     = 0
  }
}

resource "aws_security_group" "private" {
  name        = "private"
  description = "private sg of the vpc"
  vpc_id      = aws_vpc.vpc.id
  egress {
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    to_port     = 0
  }
  ingress {
    protocol    = "-1"
    cidr_blocks = [var.vpc_cidr]
    from_port   = 0
    to_port     = 0
  }
}
