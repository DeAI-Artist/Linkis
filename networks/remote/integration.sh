#!/usr/bin/env bash

# XXX: this script is intended to be run from a fresh Digital Ocean droplet

sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install -y jq unzip python3-pip software-properties-common make

# get golang
sudo snap install go --classic

#install make for building binaries of the project
apt install make        # version 4.2.1-1.2, or

#install gcc
sudo apt update
sudo apt install build-essential

## move binary and add to path
mv go /usr/local
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile

## create the goApps directory, set GOPATH, and put it on PATH
mkdir goApps
echo "export GOPATH=/root/goApps" >> ~/.profile
echo "export PATH=\$PATH:\$GOPATH/bin" >> ~/.profile
# **turn on the go module, default is auto. The value is off, if mintai source code
#is downloaded under $GOPATH/src directory
echo "export GO111MODULE=on" >> ~/.profile

source ~/.profile

mkdir -p $GOPATH/src/github.com/DeAI-Artist
cd $GOPATH/src/github.com/DeAI-Artist
# ** use git clone instead of go get.
# once go module is on, go get will download source code to
# specific version directory under $GOPATH/pkg/mod the make
# script will not work
git clone https://github.com/DeAI-Artist/Linkis.git
cd MintAI

## build
make build
#** need to install the package, otherwise mintai testnet will not execute
make install

# generate an ssh key
ssh-keygen -f $HOME/.ssh/id_rsa -t rsa -N ''
echo "export SSH_KEY_FILE=\"\$HOME/.ssh/id_rsa.pub\"" >> ~/.profile
source ~/.profile

# return to root and create /.mintai/ directory
cd
mkdir /root/.mintai

#Copy your config files to the directory, replace corresponding node file and ips needed
scp -r $HOME/MintAI/mytestnet/node0/* root@$IP0:/root/.mintai/

mintai node

echo "congratulations, your testnet is ready to run :)"
