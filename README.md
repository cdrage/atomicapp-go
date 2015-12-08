# Atomic App in Go!
####  (A Go implementation of the [Nulecule](https://github.com/projectatomic/nulecule) specification)

## Code dependencies
This implemention is currently __100%__ in Go. `1.5.1` or higher is required. You can find the current Go dependencies we rely on within the `scripts/deps.sh` script.

## Other depedencies
Required: 
  - Version 1.8.3 of Docker or higher
Option:
  - Kubernetes provider: A version of Kubernetes using the v2 API. [Fedora Instructions](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/getting-started-guides/fedora/fedora_manual_config.md) [CentOS instructions](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/getting-started-guides/centos/centos_manual_config.md)

## Installation
```bash
go get github.com/cdrage/atomicapp-go
cd $GOPATH/src/github.com/cdrage/atomicapp-go
make build
```

```bash
▶ ./atomicapp 
NAME:
   atomicapp - A Nulecule implementation written in Go

USAGE:
   ./atomicapp [global options] command [command options] [arguments...]
   
VERSION:
   0.1.1

...
```

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

## Makefile development

99% of the development work-flow is within the Makefile. Develop. Use it. Contribute.

`make build`
Build an `atomicapp` binary blob for you.

`make deps`
This is ran on each build. Go checks what build dependencies you are missing and fixes them accordingly.

`make updatedeps`
Verbose check on each depedency if it's been updated.

`make test`
Runs `go fmt` in all directories.

`make clean`
Remove `atomicapp` binary blob.

`make format`
Runs `go fmt` in all directories.

## Supported Providers
The following providers currently supported by this implementation are:
  * Docker
  * Kubernetes
  * Openshift (soooonnn!)

## Communication channels

* IRC: #nulecule (On Freenode)
* Mailing List: [container-tools@redhat.com](https://www.redhat.com/mailman/listinfo/container-tools)

## Copyright

Copyright (C) 2015 Red Hat Inc.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

The GNU Lesser General Public License is provided within the file lgpl-3.0.txt.
