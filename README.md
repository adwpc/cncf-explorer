# cncf-explorer [[中文]](https://github.com/adwpc/cncf-explorer/blob/main/README_CN.md)
explore cncf repo language percentage

# Quick Start

```
go build
./cncf-explorer
```
you will see cncf project detail list and a result
```
Name                                    Relation            Language            Percent   Git Repo
================================================================================================
Cloud Custodian                         sandbox             Python              96.6      https://github.com/cloud-custodian/cloud-custodian
KubeEdge                                incubating          Go                  94.3      https://github.com/kubeedge/kubeedge
OpenYurt                                sandbox             Go                  97.9      https://github.com/openyurtio/openyurt
......
......

percentage of languages used in cncf projects
filter: sandbox incubating graduated
====================
Go             58        0.69
Python         4         0.05
TypeScript     4         0.05
C++            3         0.04
Java           3         0.04
Makefile       3         0.04
Rust           3         0.04
unknown        2         0.02
Shell          2         0.02
Ruby           1         0.01
Lua            1         0.01
Total	 84
```

# Usage
```
Usage of ./cncf-explorer:
  -a	calc all cncf project, otherwise only: sandbox incubating graduated
  -c int
    	spider cycle (default 500)
  -o string
    	output file name (default "output.csv")
```
Note:
* will get error if cycle is too small
* will skip some project that didn't have a repo url
* will take a long time if using -a to get a result, please wait
