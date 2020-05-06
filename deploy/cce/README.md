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
The generated yaml is not for final usage. You need provide secret `bot-secret` which contains the gitee and github token as well.

```
$ cd $GOPATH/src/gitee.com/openeuler/ci-bot/deploy/cce
$ kustomize build .
```
