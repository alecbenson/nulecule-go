# Nulecule-go
####  (A Golang implementation of the Nulecule spec:  https://github.com/projectatomic/nulecule)

## Dependencies
- Version 1 of [go-yaml](https://github.com/go-yaml/yaml/tree/v1)
  - Also available through Yum: `yum install golang-gopkg-yaml`
- A version of docker that supports the `--format flag` for `docker -V`.
  -  [Included since PR #14194](https://github.com/docker/docker/pull/14194)
- You will likely want a working version of [Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes) installed on your system:
    -  [Fedora Instructions](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/getting-started-guides/fedora/fedora_manual_config.md)
    - [CentOS instructions](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/getting-started-guides/centos/centos_manual_config.md)
- [Sirupsen's Logrus](https://github.com/Sirupsen/logrus)
- [Codegangsta's CLI](https://github.com/codegangsta/cli)

## Installation
Clone the repository to within your `$GOPATH` and build the binary using the included makefile
Then, move the resulting binary somewhere within your `$PATH`
The program can now be invoked with `atomicgo`

## Example Usage
### Installing a Nulecule application
 * Make a directory in which your application will be installed in and `cd` into it.
 * install a valid projectatomic app -- for example, any of the following will work:
   * `atomicgo install projectatomic/guestbookgo-app`
   * `atomicgo install projectatomic/helloapache --destination=/home/alecbenson/Desktop/testproject`
     * If no `--destination` flag is provided, the current working directory is implicitly used

### Running a Nulecule application
 Simply deploy the application: `atomicgo run`

 You may also specify where to run the project from by specifying a directory after `run`:
   * `atomicgo run /home/alecbenson/Desktop/testproject`

  Before running your application, you will notice that there is now an `answers.conf.sample` file. It contains default values for all parameters provided in the `Nulecule` file that is also within your installation directory. You may edit any of the values within this file. By renaming the sample file to `answers.conf`, these values will be implicitly provided when the application is run.

By running the project with the `--ask` flag, the program will prompt the user for any parameters that are not specified in the `answers.conf` file (if it exists):
  * `atomicgo run /home/alecbenson/Desktop/testproject --ask`

You may also provide the `--write` flag to tell the program where to look for your answers.conf file. This is useful if you already have an answers file somewhere on your system. For example, both of the following are valid:
  * `atomicgo run . --write=/home/abenson/Desktop/`
  * `atomicgo run . --write=/home/abenson/Desktop/answers.conf`
If no `--write` flag is provided, the program looks for the answers file in the installation directory by default.

Verify that your application is running: `kubectl get pods`

### Un-deploying a Nulecule application
When you are done with your application, simply run `atomicgo stop` in your installation directory

### Running tests
You may run the project tests with `make test`

### Supported Providers
The following providers are currently supported:
  * Docker
  * Kubernetes

  Openshift will be available to use as a provider soon.
