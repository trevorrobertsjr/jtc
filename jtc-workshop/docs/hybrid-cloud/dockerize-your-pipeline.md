---
sidebar_position: 5
---

# Docker Overview

## What is Docker?
Docker is an open platform for developing, shipping, and running applications. Docker accomplishes this by using containers.

## What are containers?
A lightweight, portable, and isolated runtime environment that includes an application and its dependencies, running consistently across different computing environments.

## What are container images?
A compressed tar archive (.tar.gz) that contains the container filesystem layers and the information that a container runtime needs to know to run your application.

## Why do we need all of this?
How much do you enjoy managing python virtual environments? 🥴 Python developers have it relatively easy compared to C\C++ developers who regularly have to manage a slew of library dependencies to get their applications to run correctly. Sometimes the correct package will be installed on your server, but it is the wrong version for one application, while it is the right version for another application.

The following image depicts how you can use Docker containers (the characters running on their respective treadmills) running different versions of Python with their own separate packages that do not conflict with each other.

![Manage Dependencies with Docker](/img/Run-Different-Python-Versions-With-Docker_Watermarked.webp)

Let's see how to package our pipeline code with Docker.

