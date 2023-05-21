# Configure aws cli

## Setup account and access key

[aws-cli](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html)

1. Create Account (with required permission like ECR full access).
2. Set MFA device like [Authy](https://authy.com/).
3. Create Access Key.
4. `aws configure --profile <your-long-term-user>` and set 3 keys.
5. (Optional) Add `aws_mfa_device = arn:aws:iam:...:mfa/<device_name>` to `~/.aws/credentials`.


## Get temporary access key with `aws-mfa`

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
