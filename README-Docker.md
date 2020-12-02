# Docker README

A `Dockerfile` has been provided in order to make this code more convenient to run. To get started via `Docker`, you'll need:

* A working Docker environment
* A `git checkout` of this repository
* A copy of `GeoLite2-Country.mmdb` in the root of this repository

## Building Your Docker Container

After you have checked out this repository and added `GeoLite2-Country.mmdb` to the root directory of the repository, run this command from the root directory:

```bash
docker build . -t mmdb-from-go
```

This will build your Docker container, build your `Go` code and also run the example code for you. To log in to your container:

```bash
docker run -it --volume $PWD:/project mmdb-from-go:latest /bin/bash
```

Once you have logged in, you can test out your freshly created database:

```bash
mmdbinspect -db GeoLite2-Country.mmdb \
-db GeoLite2-Country-with-Department-Data.mmdb \
56.0.0.1 56.1.0.0/24 56.2.0.54 56.3.0.1 | less
```

You can freely edit the code outside of the Docker container and then re-run it from inside the container:

```bash
./mmdb-from-go-blogpost
```