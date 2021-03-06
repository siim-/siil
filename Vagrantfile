# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure(2) do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://atlas.hashicorp.com/search.
  config.vm.box = "debian/jessie64"

  # Disable automatic box update checking. If you disable this, then
  # boxes will only be checked for updates when the user runs
  # `vagrant box outdated`. This is not recommended.
  # config.vm.box_check_update = false

  # Create a forwarded port mapping which allows access to a specific port
  # within the machine from a port on the host machine. In the example below,
  # accessing "localhost:8080" will access port 80 on the guest machine.
  # config.vm.network "forwarded_port", guest: 80, host: 8080

  # Create a private network, which allows host-only access to the machine
  # using a specific IP.
  config.vm.network "private_network", ip: "192.168.56.104"

  # Create a public network, which generally matched to bridged network.
  # Bridged networks make the machine appear as another physical device on
  # your network.
  # config.vm.network "public_network"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  config.vm.synced_folder "./", "/go/src/github.com/siim-/siil"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  config.vm.provider "virtualbox" do |vb|
  #   # Display the VirtualBox GUI when booting the machine
  #   vb.gui = true
  #
  #   # Customize the amount of memory on the VM:
    vb.memory = "4096"
  end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Define a Vagrant Push strategy for pushing to Atlas. Other push strategies
  # such as FTP and Heroku are also available. See the documentation at
  # https://docs.vagrantup.com/v2/push/atlas.html for more information.
  # config.push.define "atlas" do |push|
  #   push.app = "YOUR_ATLAS_USERNAME/YOUR_APPLICATION_NAME"
  # end

  # Enable provisioning with a shell script. Additional provisioners such as
  # Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
  # documentation for more information about their specific syntax and use.
  config.vm.provision "shell", inline: <<-SHELL

    cd /tmp
    wget --quiet https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz

    #Write the export directives to /etc/profile since this is a vagrant machine
    sudo echo "export PATH=$PATH:/usr/local/go/bin:/go/bin" >> /etc/profile
    export PATH=$PATH:/usr/local/go/bin:/go/bin
    sudo echo "export GOPATH=/go" >> /etc/profile
    export GOPATH=/go

    sudo echo "deb http://http.debian.net/debian jessie-backports main" > /etc/apt/sources.list.d/backports.list
    sudo apt-get -qq update && sudo apt-get install -qq -y mercurial git
    curl -sSL https://get.docker.com/ | sh
    sudo usermod -a -G docker vagrant

    sudo wget --quiet -O /usr/local/bin/docker-compose https://github.com/docker/compose/releases/download/1.4.1/docker-compose-`uname -s`-`uname -m`
    sudo chmod +x /usr/local/bin/docker-compose

    docker pull httpd:latest
    docker pull mysql:latest
    docker pull redis:latest

    sudo chown -R vagrant:vagrant /go
    sudo chown -R vagrant:vagrant /usr/local/go

    go get github.com/siim-/siil

    cd /go/src/github.com/siim-/siil

    # Generate self signed cert for local development
    mkdir certs
    cd certs
    sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt -subj "/C=EE/ST=Harjumaa/L=Tallinn/O=Siil/OU=Development/CN=www.siil.lan"
    sudo chmod +x revocationlists.sh
    sudo chmod +x eidcert.sh

    # Download eID cert
    ./eidcert.sh

    # Download revocation lists
    ./revocationlists.sh

    cd ..

    #Set the correct docker-compose file
    cp docker-compose.yml.dev docker-compose.yml

    #And apache2 conf
    cp conf/apache2-vhosts/siil.conf.dev conf/apache2-vhosts/active/siil.conf

    cp conf/apache2-vhosts/siil.conf.dev
    #Bring up supporting services
    docker-compose up -d

  SHELL

  config.vm.provision "shell", inline: "cd /go/src/github.com/siim-/siil && docker-compose up -d", run:"always"
end
