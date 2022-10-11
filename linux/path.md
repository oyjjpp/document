### 简介
>环境变量PATH，设置程序的快捷访问方式


### 查看当前环境变量的值

echo $PATH

    [root@php]$echo $PATH
    /usr/local/go/bin:/usr/local/common/jdk/bin:/usr/local/common/jdk:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin:/data/server/git/bin:/root/bin:/root/bin
    [root@php]$


### 命令行编辑环境变量PATH

###### 1、直接命令修改

export PATH=$PATH:新添加的路径

    [root@php]$export PATH=$PATH:/data/server/php/bin
    [root@php]$echo $PATH
    /usr/local/go/bin:/usr/local/common/jdk/bin:/usr/local/common/jdk:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin:/data/server/git/bin:/root/bin:/root/bin:/data/server/php/bin
    [root@php]$type php
    php is /data/server/php/bin/php
    [root@php]$



###### 2、命令行删除

export PATH=不包含去掉的路径值

    [root@php]$echo $PATH
    /usr/local/go/bin:/usr/local/common/jdk/bin:/usr/local/common/jdk:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin:/data/server/git/bin:/root/bin:/root/bin:/data/server/php/bin
    [root@php]$export PATH=/usr/local/go/bin:/usr/local/common/jdk/bin:/usr/local/common/jdk:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin:/data/server/git/bin:/root/bin:/root/bin[root@php]$echo $PATH
    /usr/local/go/bin:/usr/local/common/jdk/bin:/usr/local/common/jdk:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/root/bin:/data/server/git/bin:/root/bin:/root/bin
    [root@php]$type php
    -bash: type: php: not found
    [root@php]$


###### 3、注意

 - 用命令行修改PATH,只针对当前shell有用,关了终端就失效了

### 文本文件修改

###### 1、文件类型

普通文本文件修改（修改PATH每次登陆需要source）  
系统文件修改（首次设置好会一直生效）

###### 2、案列

修改.bashrc文件
>export PATH=$PATH:新增路径  
export PATH=路径1:路径2:路径n


###### 3、注意
 - 要立即生效只需要source一下.bashrc文件,之后无需再source
 - 要恢复系统默认PATH,删掉.bashrc中的修改语句即可

### 参考

  - [Linux系统修改PATH环境变量方法](https://www.cnblogs.com/cursorhu/p/5806596.html)
  - [Ubuntu系统环境变量详解](https://blog.csdn.net/netwalk/article/details/9455893)
