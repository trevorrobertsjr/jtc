---
sidebar_position: 6
---

# Create your Docker Container

## The Container Image Blueprint

To create your container image, you need to write some instructions for the image builder in a file called `Dockerfile`.

Here is what our Dockerfile looks like for our script:
```Dockerfile
# Use an official Python image
FROM python:3.12-slim

# Set working directory
WORKDIR /app

# Copy the requirements file to the container
COPY requirements.txt ./

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy any remaining files that are needed for the application.
COPY . .

# Run the script
CMD ["python", "pipeline.py"]
```

Dockerfiles are mostly self-explanatory, but I'll explain a few items to be aware of.

`FROM` tells us which container image will be used as the foundation for our container image. Always use trusted images from official sources. In a corporate setting, access to public images is usually blocked to prevent an attack vector. Some companies may also go as far as to work with a vendor (ex: [Chainguard](www.chainguard.dev)) that specializes in providing hardened minimal containers to use with the `FROM` keyword

Multiple `COPY` commands: I do this to save time in building future versions of our container image. Each command in a Dockerfile roughly represents a layer that the image builder creates. If I have a single copy command (i.e. `COPY pipeline.py requirements.txt ./`), then every time I make a source code change, I need to re-run that entire layer including the pip package installation process.

`RUN` command uses the `--no-cache-dir` option with pip to save on file system space. The goal is to make the container image size as small as possible for ease of portability.

If you want to dive deeper on the topic of writing a Dockerfile, checkout this site: https://docs.docker.com/get-started/docker-concepts/building-images/writing-a-dockerfile/

## Building with the Blueprint

Now that we have our Dockerfile, we need to build our container image before we can use it.

The syntax is:

```bash
docker build -t dockerhub-username/pipeline .
```

For those of you with Macs that have the Apple CPU, you will need to use a special syntax to build an Intel-compatible image:

```bash
docker buildx build --platform linux/amd64 -t dockerhub-username/pipeline .
```

You will see output similar to the following:
```bash
❯ docker buildx build --platform linux/amd64 -t dockerhub-username/pipeline .
[+] Building 8.4s (8/10)                                    docker:desktop-linux
 => [internal] load build definition from Dockerfile                        0.0s
 => => transferring dockerfile: 398B                                        0.0s
 => [internal] load metadata for docker.io/library/python:3.12-slim         0.6s
 => [auth] library/python:pull token for registry-1.docker.io               0.0s
 => [internal] load .dockerignore                                           0.0s
 => => transferring context: 2B                                             0.0s
 => [1/5] FROM docker.io/library/python:3.12-slim@sha256:34656cd9045634904  0.0s
 => [internal] load build context                                           0.0s
 => => transferring context: 738B                                           0.0s
 => CACHED [2/5] WORKDIR /app                                               0.0s
 => [3/5] COPY requirements.txt ./                                          0.0s
 => [4/5] RUN pip install --no-cache-dir -r requirements.txt                7.8s
 => => # Downloading pyasn1-0.6.1-py3-none-any.whl (83 kB)
 => => # Downloading pyasn1_modules-0.4.1-py3-none-any.whl (181 kB)
 => => # Downloading python_dateutil-2.9.0.post0-py2.py3-none-any.whl (229 kB)
 => => # Downloading python_dotenv-1.0.1-py3-none-any.whl (19 kB)
 => => # Downloading pytz-2025.1-py2.py3-none-any.whl (507 kB)
 => => # Downloading requests-2.32.3-py3-none-any.whl (64 kB)
 => [5/5] COPY . .                                                          0.0s
 => exporting to image                                                      0.8s
 => => exporting layers                                                     0.8s
 => => writing image sha256:4d3b401bb4e0f2fa0e688c6a2063c61b1da618ed42a237  0.0s
 => => naming to docker.io/dockerhub-username/pipeline                         0.0s
```
For our image to be easily accessible by others on your team, you need to store it in a repository, similar to GitHub, that is called a container registry. In our lab, we are using Docker Hub, the default registry that Docker uses as the source and destination of container images.

We use the following command to push our image:
```bash
docker push dockerhub-username/pipeline
```

You should see output similar to the following:
```bash
❯ docker push dockerhub-username/pipeline
Using default tag: latest
The push refers to repository [docker.io/dockerhub-username/pipeline]
125035b11658: Pushed
3614b971016b: Pushing  123.5MB/325.7MB
23f4431f4f52: Pushed
626938010757: Layer already exists
05956879d134: Layer already exists
7dc17f2f9831: Layer already exists
c3844976239d: Layer already exists
7914c8f600f5: Layer already exists
```

Similar to Git, Docker is intelligent enough to push only the changes you have made since the last push to save on bandwidth consumption and to speed up operations.

Now, that we have designed and built our container image, let's run our container!.

