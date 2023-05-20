# CICD

Deploying dockerized app with aws.

## Git branch

We can create new branch `ft/docker` and check it out with

```bash
$ git branch -b ft/docker
```

or

```bash
$ git branch ft/docker
$ git checkout ft/docker
```

## Add app docker image

We can find many prebuilt images with [dockerhub](https://hub.docker.com/).

### Dockerfile

If we write

```dockerfile
FROM golang:1.18.10-alpine3.17 AS builder
```

then docker starts building with given image and pull it from dockerhub if no image on local.

There's several useful commands,
see [docker cli command](https://docs.docker.com/engine/reference/commandline/docker/) and
[dockerfile command](https://docs.docker.com/engine/reference/builder/) for more details.

- `WORKDIR /app`: change base working directory to `/app`.
- `COPY --from=builder /app/main .`: copy files or dirs from the image.
- `EXPOSE 8080`: specify what a port should be bounded to.
- `CMD [ "/app/main" ]`: command will be run when a container run. `ENTRYPOINT` is similar, but a bit different.

## Multistage build

Containers and Images have to be as small as possible.
We can isolate the final output image by other build-time only component like compilers, dependencies, tests, and so on with
[multistage build](https://docs.docker.com/build/building/multi-stage/).

## References
- https://hub.docker.com/
- https://docs.docker.com/engine/reference/commandline/docker/
- https://docs.docker.com/engine/reference/builder/
- https://kapeli.com/cheat_sheets/Dockerfile.docset/Contents/Resources/Documents/index
- https://docs.docker.com/build/building/multi-stage/
