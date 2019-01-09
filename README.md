## Overview
- Data domain file system GC
  ddcli -host 10.98.xx.xx -user sysadmin -password changeme -x=0
- Data domain restart ddboost
  ddcli -host 10.98.xx.xx -user sysadmin -password changeme -x=1
- Data domain delete all storage unit
  ddcli -host 10.98.xx.xx -user sysadmin -password changeme -x=2
- Data domain delete oldest storage unit
  ddcli -host 10.98.xx.xx -user sysadmin -password changeme -x=3
- Data movement to cloud
  ddcli -host 10.98.xx.xx -user sysadmin -password changeme -x=4
- Avamar GC
  ddcli -host 10.98.xx.xx -user admin -password changeme -x=5

## Usage
```bash
Usage of ddcli.exe:
  -host string
        Remote server host name or ip address
  -password string
        Remote server password
  -user string
        Remote server user name
  -v    prints current version and exits
  -x int
        0: Data domain file system GC
        1: Data domain restart ddboost
        2: Data domain delete all storage unit
        3: Data domain delete oldest storage unit
        4: Data movement to cloud
        5: Avamar GC
```
e.g.
### Data domain file system GC 
Include those dd commands:
- filesys clean start
- filesys clean watch

```bash
ddcli -host=<remote server host> -user=<remote server user name> -password=<remote server password> -x=0
```
![T](https://github.com/xiak/datadomain/blob/master/src/filesys-clean.PNG?raw=true)

### Data domain restart ddboost 
Include those dd commands:
- ddboost disable
- ddboost enable

```bash
ddcli -host=<remote server host> -user=<remote server user name> -password=<remote server password> -x=1
```
![T](http://github.com/xiak/datadomain/raw/master/src/ddboost-restart.PNG)

### Data domain delete all storage unit
Include those dd commands:
- ddboost storage-unit show
- ddboost storage-unit delete (Loop)

```bash
ddcli -host=<remote server host> -user=<remote server user name> -password=<remote server password> -x=2
```
![T](http://github.com/xiak/datadomain/raw/master/src/delete-all-storage-unit.PNG)

### Data domain delete oldest storage unit
Include those dd commands:
- ddboost storage-unit show
- ddboost storage-unit delete

```bash
ddcli -host=<remote server host> -user=<remote server user name> -password=<remote server password> -x=3
```
![T](http://github.com/xiak/datadomain/raw/master/src/delete-oldest-storage-unit.PNG)

### Data domain datamovement to cloud
Include those dd commands:
- data-movement start to-tier cloud
- data-movement watch

```bash
ddcli -host=<remote server host> -user=<remote server user name> -password=<remote server password> -x=4
```

### Avamar GC
Include those dd commands:
- avmaint sched stop --ava
- avmaint checkpoint --ava --wait
- avmaint --ava garbagecollect
- avmaint --ava gcstatus

```bash
ddcli -host=<remote server host> -user=<remote server user name> -password=<remote server password> -x=5
```
## Develop from master
If you want to build a package yourself, or contribute. Here is a guide for how to do that. You can always find
the latest master builds [here](http://github.com/xiak/datadomain/bin/)

### Dependencies
- Go 1.8.1

### Building
```bash
git clone http://github.com/xiak/datadomain.git
go run build.go setup
go run build.go build
```

## Contribute
If you have any idea for an improvement or found a bug do not hesitate to open an issue.
And if you have time clone this repo and submit a pull request and help me make Diablo
better that is all we dream about it!

## License
ddcli is distributed under Apache 2.0 License.
