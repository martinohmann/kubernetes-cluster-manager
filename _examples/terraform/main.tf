# This is just an example which shows how to pass kubernetes credentials from
# terraform to kcm. For that purpose we assume to already have an AWS EKS
# cluster from which we read data to create a kubeconfig. In a real world
# example we would also create the cluster and all the necessary parts like the
# VPC, Security Groups, Roles and Node Pools here, but this is out of scope as
# it totally depends on your requirements and setup, and would make the example
# quite large.
#
# The example will retrieve data about the cluster from AWS and then builds a
# kubeconfig and makes kcm aware.
#
# If you just want to play around with that, go to the AWS console and create
# an EKS cluster manually and adjust the region and clustername in
# `config.auto.tfvars` if necessary.
variable "clustername" {
  default = "cluster"
}

variable "profile" {
  default = "default"
}

variable "region" {
  default = "eu-west-1"
}

locals {
  context = "eks-${var.clustername}"
}

provider "aws" {
  shared_credentials_file = "~/.aws/credentials"
  region                  = "${var.region}"
  profile                 = "${var.profile}"
}

data "aws_eks_cluster" "cluster" {
  name = "${var.clustername}"
}

data "template_file" "kubeconfig" {
  template = "${file("${path.module}/kubeconfig.tpl")}"

  vars = {
    endpoint                   = "${data.aws_eks_cluster.cluster.endpoint}"
    certificate_authority_data = "${data.aws_eks_cluster.cluster.certificate_authority.0.data}"
    clustername                = "${var.clustername}"
    context                    = "${local.context}"
  }
}

resource "local_file" "kubeconfig" {
  content  = "${data.template_file.kubeconfig.rendered}"
  filename = "${path.module}/kubeconfig"
}

output "kubeconfig" {
  value = "${local_file.kubeconfig.filename}"
}

output "context" {
  value = "${local.context}"
}

# # kcm also supports kubernetes credentials in the form of a server-token pair.
# # Note that these will be ignored if a kubeconfig is set:
#
# output "server" {
#   value = "https://localhost:6443"
# }
#
# output "token" {
#   value = "kubernetes-serviceaccount-token"
# }
#
# # You can also pass arbitrary outputs to kcm. These are then made available
# # during kubernetes manifest rendering:
#
# output "mycustomvariable" {
#   value = "myvalue"
# }
#
# This will make it available as {{ .Values.mycustomvariable }} in helm templates.

