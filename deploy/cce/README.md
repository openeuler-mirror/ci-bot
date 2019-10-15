# Getting Started on CCE

## Build

```
$ mkdir -p $GOPATH/src/gitee.com/openeuler
$ cd $GOPATH/src/gitee.com/openeuler
$ git clone https://gitee.com/openeuler/ci-bot
$ cd ci-bot
$ make ci-bot-image
```

## Usage

```
$ cd $GOPATH/src/gitee.com/openeuler/ci-bot/deploy/cce
$ kubectl create -f .
```
