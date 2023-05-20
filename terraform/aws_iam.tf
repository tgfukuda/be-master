resource "aws_iam_instance_profile" "oracle_instance_profile" {
  name = "oracle-role"
  role = aws_iam_role.oracle_role.name
}

resource "aws_iam_role" "oracle_role" {
  name = "oracle-role"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
               "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF

  tags = {
    "Name" = "oracle-role"
  }
}

data "aws_iam_policy" "ec2_role_for_ssm" {
  arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "ssm_policy_attachment" {
  role = aws_iam_role.oracle_role.name
  policy_arn = data.aws_iam_policy.ec2_role_for_ssm.arn
}

resource "aws_iam_role_policy" "logger" {
  name = "oracle-logger"
  role = aws_iam_role.oracle_role.id

  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Action": [
          "s3:*Object"
        ],
        "Resource": [
          "${aws_s3_bucket.log.arn}",
          "${aws_s3_bucket.log.arn}/*"
        ]
      }
    ]
  })
}
