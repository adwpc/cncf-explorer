# cncf-explorer
分析CNCF项目语言占比

# 用法

```
go build
./cncf-explorer
```
可以看到每个项目详细情况，最后还会打印整体结果
下面的结果表示，沙盒、孵化中、毕业项目的语言占比
```
Name                                    Relation            Language            Percent   Git Repo
================================================================================================
Cloud Custodian                         sandbox             Python              96.6      https://github.com/cloud-custodian/cloud-custodian
KubeEdge                                incubating          Go                  94.3      https://github.com/kubeedge/kubeedge
OpenYurt                                sandbox             Go                  97.9      https://github.com/openyurtio/openyurt
......
......

percentage of languages used in cncf projects
CNCF项目语言占比
filter: sandbox incubating graduated
类型：沙盒、孵化中、毕业
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
        加-a表示分析所有项目，不加-a表示只分析：沙盒、孵化中、毕业项目
  -c int
    	spider cycle (default 500)
        爬虫周期，每隔多少ms启动一个 
  -o string
    	output file name (default "output.csv")
        输出csv表格到文件
```
Note:
* 如果周期设置太小，会报错，建议默认的500
* 如果没有找到git repo 地址，爬虫会跳过该项目，实际上很多初级项目没有地址
* 分析全部项目结果需要一些时间，不要杀掉进程
