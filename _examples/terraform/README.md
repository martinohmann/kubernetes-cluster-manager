Terraform provisioner example
=============================

This example is very basic and assumes that you already have an EKS cluster on
AWS. See comments in [`main.tf`] for further explainations.

Copy tfvars file and adjust to your needs:

```sh
$ cp config.auto.tfvars.example config.auto.tfvars
```

Initialize terraform:

```sh
$ terraform init
```

Run `kcm provision` in dry-run mode to see what happens.

```sh
$ kcm provision --config config.yaml --dry-run
```

kcm will retrieve the `kubeconfig` and `context` from the terraform output.
Check [`main.tf`](main.tf) for an example how to pass these values from
terraform to kcm. You can also provide the kubernetes credentials to kcm via
the `--cluster-*` command line flags.
