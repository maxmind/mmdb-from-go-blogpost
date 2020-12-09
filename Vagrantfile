# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://atlas.hashicorp.com/search.
  config.vm.box = 'bento/ubuntu-20.04'

  config.vm.provision 'shell', inline: <<-SHELL

      curl --silent --location \
      https://github.com/maxmind/mmdbinspect/releases/download/v0.1.1/mmdbinspect_0.1.1_linux_amd64.deb \
      -o /tmp/mmdbinspect.deb

      dpkg -i /tmp/mmdbinspect.deb

      apt-get update && apt-get install --no-install-recommends -y golang-go less

      cd /vagrant

      go build
      ./mmdb-from-go-blogpost
  SHELL
end
