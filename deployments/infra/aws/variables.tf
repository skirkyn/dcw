#  general
variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}
variable "access_key" {
  description = "access key"
  type        = string
  default     = "<access key>"
}
variable "secret_key" {
  description = "secret key"
  type        = string
  default     = "<secret key>"
}


#  vpc
variable "vpc_name" {
  description = "name of the vpc"
  type        = string
  default     =  "dcw_vpc"
}


variable "cidr" {
  description = "cidr"
  type        = string
  default     =  "10.0.0.0/16"
}



variable "machine_image" {
  description = "Amazon AMI"
  type        = string
  default     =  "AL2_x86_64"
}

# controller
variable "controller_instance_type" {
  description = "Controller machine type"
  type        = string
  default     = "t3.small"
}
variable "controller_min_instances" {
  description = "Controller min instances"
  type        = number
  default     = 1
}
variable "controller_max_instances" {
  description = "Controller max instances"
  type        = number
  default     = 1
}
variable "controller_desired_instances" {
  description = "Controller desired instances"
  type        = number
  default     = 1
}

variable "controller_eks_node_group" {
  description = "Controller node group name"
  type        = string
  default     = "controller_nodes"
}

# worker
variable "worker_instance_type" {
  description = "Worker machine type"
  type        = string
  default     = "t3.small"
}

variable "worker_min_instances" {
  description = "Worker min instances"
  type        = number
  default     = 3
}
variable "worker_max_instances" {
  description = "Worker max instances"
  type        = number
  default     = 10
}
variable "worker_desired_instances" {
  description = "Worker desired instances"
  type        = number
  default     = 3
}

variable "worker_eks_node_group" {
  description = "Worker node group name"
  type        = string
  default     = "Worker"
}