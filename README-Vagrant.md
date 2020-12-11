# Vagrant README

<!-- vim-markdown-toc GFM -->

* [Getting Started](#getting-started)
* [Building Your Vagrant VM](#building-your-vagrant-vm)
* [Logging in to Your VM](#logging-in-to-your-vm)
* [Inspecting Your Database](#inspecting-your-database)
* [Refactoring Your Code](#refactoring-your-code)
  * [Re-creating the database](#re-creating-the-database)

<!-- vim-markdown-toc -->

## Getting Started

A `Vagrantfile` has been provided in order to make this code more convenient to
run. To get started via `Vagrant`, you'll need:

* A working `Vagrant` environment
* A `git checkout` of this repository
* A copy of `GeoLite2-Country.mmdb` in the root of this repository

## Building Your Vagrant VM

After you have checked out this repository and added a `GeoLite2-Country.mmdb` to
the root directory of the repository, run this command from the root directory.

```bash
vagrant up
```

## Logging in to Your VM

This will build your Vagrant VM, build your `Go` code and also run the example
code for you. To log in to your container, run this command:.

```bash
vagrant ssh
cd /vagrant
```

## Inspecting Your Database

Once you have logged in, you can test out your freshly created database:

```bash
mmdbinspect -db GeoLite2-Country.mmdb \
-db GeoLite2-Country-with-Department-Data.mmdb \
56.0.0.1 56.1.0.0/24 56.2.0.54 56.3.0.1 | less
```

## Refactoring Your Code

You can now freely edit the code outside of the Vagrant VM and re-run it from
inside the VM.

### Re-creating the database

After you have edited the files, log in to your VM using the instructions above
and run the following code:

```bash
go build && ./mmdb-from-go-blogpost
```
