# CICD

Deploying dockerized app with aws.
docker is very simple and powerful environment and instance management tool. (Who can't use it is not a developer)

## Git branch

We can create new branch `ft/docker` and check it out with `git branch -b ft/docker` or

```bash
$ git branch ft/docker
$ git checkout ft/docker
```

## Add app docker image

We can find many prebuilt images with [dockerhub](https://hub.docker.com/).
See [Dockerfile](./Dockerfile).

### Dockerfile

If we write

```dockerfile
FROM golang:1.18.10-alpine3.17 AS builder
```

then docker starts building with given image and pull it from dockerhub if no image on local.

There's several useful commands,
see [dockerfile command](https://docs.docker.com/engine/reference/builder/) for more details.

- `WORKDIR /app`: change base working directory to `/app`.
- `COPY --from=builder /app/main .`: copy files or dirs from the image.
- `EXPOSE 8080`: specify what a port should be bounded to.
- `CMD [ "/app/main" ]`: command will be run when a container run. `ENTRYPOINT` is similar, but a bit different.
    If Both `CMD` and `ENTRYPOINT` defined at the same time, `CMD` acts as additional parameters of `ENTRYPOINT`.

### Multistage build

Containers and Images have to be as small as possible.
We can isolate the final output image by other build-time only component like compilers, dependencies, tests, and so on with
[multistage build](https://docs.docker.com/build/building/multi-stage/).

### build the image

See [docker cli command](https://docs.docker.com/engine/reference/commandline/docker/).

`sudo docker build -t simplebank:latest .` build an image of Dockerfile in the current directory with the name `simplebank` and the tag `latest`.

### Run a container

```bash
$ sudo docker --name simplebank -p 8080:8080 -e GIN_MODE=release simplebank:latest
```

run a container named `simplebank` and bounded the container port `8080` to the host port `8080` and the environment `GIN_MODE`
with the image `simplebank:latest`.

### Inspect containers

`sudo docker container inspect <container_name>` show you the setting of `<container_name>` and
`sudo docker logs <container_name>` emit the logs of `<container_name>`.

### Dig deeper docker network

`sudo docker network ls` list the networks like

```bash
NETWORK ID     NAME                   DRIVER    SCOPE
0e8f55845d03   bridge                 bridge    local
1eb55a1a7bad   host                   host      local
28841caf08d1   none                   null      local
```

and we can find a network detail with `sudo docker network inspect <network_name>`.
Normally, containers in the same network can discover each other by name instead of IP address though it's not the case of the default `bridge` network.

We can create a network with `sudo docker network create <network_name>` and
put a container in it with `sudo docker network connect <network_name> <container_name>`.
`sudo docker network rm <network_name>` removes `<network_name>`.

Additionally, we can put networks with `--network` option of `docker run` and

```bash
$ sudo docker network connect simplebank-network postgres12
$ sudo docker \
    --name simplebank \
    --name simplebank-network
    -p 8080:8080 \
    -e GIN_MODE=release \
    -e "DB_SOURCE=postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable" \ <- we can use the container name for the network instance
    simplebank:latest
```

## Docker compose

docker-compose can organize multiple containers.
See [docker-compose](https://docs.docker.com/compose/) and [docker-compose.yml](./docker-compose.yml).

### What's will docker compose do?

Basically, docker-compose can do the same thing as we do with docker cli in shell.

If we have `docker-compose.yml` of

```yml
version: "3.9"  # docker compose version
services:   # services to launch
  postgres12:   # service name
    image: postgres:12-alpine   # base image
    environment:    # variables
      - POSTGRES_USER=root
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=simple_bank
    
  api:
    build:  # build with
      context: .
      dockerfile: Dockerfile
    ports:  # bounded ports
      - "8080:8080"
    environment:
      - DB_SOURCE=postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable
    depends_on: # make sure other images to be ready
      - postgres
    entrypoint: ["/app/wait-for.sh", "postgres:5432", "--", "/app/start.sh"]
    command: ["/app/main"]  # if we override entrypoint, dockerfile's CMD will ignored. see https://docs.docker.com/compose/compose-file/compose-file-v3/#entrypoint
```

then docker-compose do ...

- Build or pull the images. If it builds one, the name is prefixed with the root directory name such as `be-master-api`.
- Create and bind the network for the services if not exists. The name is prefixed with the root directory name such as `be-master_default`.
- Run containers of the services with prefix of the root directory such as `be-master-api-1`.

### wait-for script

[wait-for](https://github.com/eficode/wait-for) is a script designed to synchronize services like docker containers.
It's very simple, see [compose file](./docker-compose.yml), [start.sh](./start.sh).

## Deployment with AWS

Prepare for your [aws account](https://aws.amazon.com/jp/account/).

### ECR

[ECR](https://ap-northeast-1.console.aws.amazon.com/ecr/repositories) is a container management service provided by aws.
We use it for a example deployment workflow.
New ECR repository can be created by `Create repository` of top right bottun.

At the repository window we can see `View push command`, but basically we won't use it directly from local.
Instead, github actions performs build, tag and push.

### Setup workflow

We use https://github.com/marketplace/actions/amazon-ecr-login-action-for-github-actions for deploy.

~The workflow is (it seems to be the old version)~

```yml
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v1
  with:
  aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
  aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
  aws-region: ap-northeast-1

- name: Login to Amazon ECR
  id: login-ecr
  uses: aws-actions/amazon-ecr-login@v1

- name: Build, tag, and push image to Amazon ECR
  env:
  ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
  ECR_REPOSITORY: ecr_repository_name # the name of repo created in the previous step
  IMAGE_TAG: ${{ github.sha }}
  run: |
    docker build -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG -t $ECR_REGISTRY/$ECR_REPOSITORY:latest .
    docker push -a $ECR_REGISTRY/$ECR_REPOSITORY
```

~We use `Settings -> Secrets -> Actions` tab to set `${{ secrets.AWS_SECRET_ACCESS_KEY }}` and `${{ secrets.AWS_ACCESS_KEY_ID }}`.~

**WARN**:

The above workflow need to configure https://docs.aws.amazon.com/AmazonECR/latest/userguide/security_iam_id-based-policy-examples.html.
(PowerUser isn't enough to perform login in the case of private ECR)

```
Run aws-actions/amazon-ecr-login@v1
  
(node:1669) NOTE: We are formalizing our plans to enter AWS SDK for JavaScript (v2) into maintenance mode in 2023.

Please migrate your code to use AWS SDK for JavaScript (v3).
For more information, check the migration guide at https://a.co/7PzMCcy
(Use `node --trace-warnings ...` to show where the warning was created)
Error: User: arn:aws:iam::***:user/github-ci is not authorized to perform: ecr:GetAuthorizationToken on resource: * because no identity-based policy allows the ecr:GetAuthorizationToken action
```

will happen at `Login to Amazon ECR`.

### Production DB: RDS

We can setup RDBs (MySQL, PostgreSQL, Oracle, ...) with RDS and NoSQLs (GraphQL) with DynamoDB.

Migration will performs by app, but we can do by ourselves with some modifications to `make migrateup`

```bash
migrate -path db/migration -database "postgresql://<aws_exported_user>:<aws_exported_pass>@<RDS_Endpoint>:5432/simple_bank?sslmode=disable" -verbose up
```

## Secret Management

[Secret Manager](https://medium.com/awesome-cloud/aws-secrets-manager-overview-introduction-to-secrets-manager-getting-started-641bc722cd1a)
and
[Parameter Store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html)
allows us to store credentials for cheap price.
We use `Secret Manager` for storing app.env values of production.

### AWS cli

Install it by https://aws.amazon.com/cli/ and configure it.
[aws-mfa](https://github.com/broamski/aws-mfa) helps us login with cli.

```bash
$ aws configure # configuration will create at ~/.aws/credentials
:
```

### Ritrieve Secrets with cli

Add IAM setting of SecretManager to corresponding user, 

```bash
$ aws secretmanager get-secret-value --secret-id <secret-id or arn>
{
    "ARN": "arn:aws:secretmanager:...",
    "Name": "...",
    "SecretString": "{\"DB_SOURCE\":\"...\", ..stored secrets..}"
    "VersionStages": [
        "..."
    ],
    "CreatedDate": "..."
} 
```

returns stored ones. `"SecretString"` is a JSON for them.
To get human-readable data only for it, `aws secretmanager get-secret-value --secret-id <secret-id or arn> --query SecretString --output text`.

### Parse it with jq

[jq](https://stedolan.github.io/jq/) is a powerful cli to process JSON string.
We used it to make `SecretString` into environment variables in production.

- to_entries: make a json into key-value pairs array i.e. `{"a": 1, "b": 2}` into `{ {"key":"a", "value": 1}, {"key":"b", "value": 2}}`.
- map: apply functions to all entries.
- string interpolation of \(`op`): enable operations in a string.
- array object iterator: remove array/object brackets from output
- raw -r: remove string's quote.

```bash
$ jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env # write them to env
```

### Pull production Image from private registry

See https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login-password.html#examples and
run `sudo docker pull <aws_account_id>.dkr.ecr.<region>.amazonaws.com`.

## References
- https://hub.docker.com/
- https://docs.docker.com/engine/reference/commandline/docker/
- https://docs.docker.com/engine/reference/builder/
- https://kapeli.com/cheat_sheets/Dockerfile.docset/Contents/Resources/Documents/index
- https://docs.docker.com/build/building/multi-stage/
- https://docs.docker.com/compose/
- https://docs.aws.amazon.com/AmazonECR/latest/userguide/security_iam_id-based-policy-examples.html
- https://medium.com/awesome-cloud/aws-difference-between-secrets-manager-and-parameter-store-systems-manager-f02686604eae
- https://medium.com/awesome-cloud/aws-secrets-manager-overview-introduction-to-secrets-manager-getting-started-641bc722cd1a
