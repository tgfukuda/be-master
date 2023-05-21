resource "aws_ecr_repository" "ecr_repo" {
  name = "${var.environment_name}-ecr-repo"
}

resource "aws_db_instance" "rds_instance" {
  identifier            = "${var.environment_name}-rds-instance"
  engine                = "postgres"
  engine_version        = "12"
  instance_class        = "db.t2.micro"
  allocated_storage     = 20
  storage_type          = "gp2"
  publicly_accessible  = false
  skip_final_snapshot   = true
  username              = "root"
  password              = "${var.db_password}"
}

resource "aws_secretsmanager_secret" "secret_manager" {
  name = "${var.environment_name}-secret-manager"
}
