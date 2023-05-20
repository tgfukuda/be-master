variable "environment_name" {
  description = "An environment name that is prefixed to resource names"
  type        = string
}

variable "vpc_cidr" {
  description = "IP range (CIDR notation) for this VPC"
  type        = string
  default     = "10.68.0.0/23"
}

// public subnet conf
variable "public_subnet_region" {
  description = "Public subnet region"
  type        = string
}

variable "public_subnet_cidr" {
  description = "IP range (CIDR notation) for the public subnet"
  type        = string
  default     = "10.68.0.0/27"
}

//private subnet conf
variable "private_subnet_region_1" {
  description = "Private subnet region for the first Availability Zone"
  type        = string
}

variable "private_subnet_cidr_1" {
  description = "IP range (CIDR notation) for the public subnet in the first Availability Zone"
  type        = string
  default     = "10.68.0.32/28"
}

variable "private_subnet_region_2" {
  description = "Private subnet region for the second Availability Zone"
  type        = string
}

variable "private_subnet_cidr_2" {
  description = "IP range (CIDR notation) for the public subnet in the second Availability Zone"
  type        = string
  default     = "10.68.0.48/28"
}

variable "instance_type" {
  description = "Instance Type"
  type        = string
  default     = "m1.mediam"
}

// default: Amazon Linux
variable "ami" {
  description = "AMI ID"
  type        = string
  default     = "ami-05fa00d4c63e32376"
}

variable "volume_size" {
  description = "Input Drive Volume(GB)"
  type        = number
  default     = 64
}

variable "keypair_id" {
  description = "ssh keypair to interract with the instances"
  type = string
}