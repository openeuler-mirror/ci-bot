# ci-bot

## Introduction

This repository is used to address the code of openEuler ci bot.

## Architecture

<img src="./docs/images/architecture.png" />

## Prerequisites

You'll need to setup a MySQL Database before you are getting started.
This is an example to create Database instance.

* Setup MySQL instance by the Huawei Cloud Relational Database Service (RDS)
* Login in MySQL with your account and password
* Create database instance by running the following command
    ```
    CREATE DATABASE cibot;
    ```
The information of database instance will be used in the following Installation.

## Config
### environment variables
Some sensitive configurations items support reading from environment variables.
you can set the following environment variables in your OS:
* GITEE_TOKEN 
* WEBHOOK_SECRET
* DATABASE_HOST 
* DATABASE_PORT
* DATABASE_USERNAME
* DATABASE_PASSWORD
### label config
If you want to clear some tags when the pull request source branch changes,
 you can configure it in the configuration file(config.yaml).
 for example see config.yaml delLabels fields.Description:
 * kind,sig,openeuler-cla,priority Delete labels beginning with kind,sig,openeuler-cla or priority. 
 * lgtm Delete labels lgtm or beginning with lgtm-.
 * Except for the above description items, other labels will be judged as equal.
### extraLgtmCountRequired config
 If you want to set the number of lgtm tags for a separate repository or organization, 
 you can configure this configuration item.The configuration item is a list, and the 
 list element contains the following configuration items:
 * lcrType Indicates whether the configuration is for the repository or the organization
 * lcrName Configure the spatial address of the repository or organization 
 * lcrCount Number of lgtm tags
 
**If your configuration is for a specific repository, you should configure the full path instead of just the repository space address. For example: lcrName:openEuler/ci-bot**

## Getting Started

* [Getting Started on Locally](deploy/locally/README.md)
* [Getting Started on CCE](deploy/cce/README.md)

## Command Help

See the [Command Help](https://gitee.com/openeuler/community/blob/master/en/sig-infrastructure/command.md) file for details.
> For Chinese version, please refer to [命令帮助](https://gitee.com/openeuler/community/blob/master/zh/sig-infrastructure/command.md).

## License

See the [LICENSE](LICENSE) file for details.


