# Configure aws cli

[aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

[aws-mfa](https://qiita.com/ogady/items/c17ffe8f7c8e15b15f77)

```bash
pip3 install aws-mfa
aws configure --profile your-name-long-term
    <set your api key, secret, region, output>
aws configure --profile your-name
    <set region, output. no need to set api key, secret>
aws-mfa --profile your-name --device <mfa device id>
```

# Variables
use terraform.tfvars or cli
```
environment_name = "test" # prefix
keypair_id = "key-*********" # common keypair
public_subnet_region = "us-east-1a"
private_subnet_region_1 = "us-east-1a"
private_subnet_region_2 = "us-east-1b"
aws_profile = "******" # profile name configured by aws-mfa
```

# Deploy terraform

## Plan
```bash
terraform plan
```

## Apply
```bash
terraform apply
```

## Destruct
```bash
terraform destroy
```
