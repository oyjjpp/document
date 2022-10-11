### 简介
UNIX及UNIX-like下的调试工具  
17页
### 安装
- 检查当前机器安装情况  
rpm -qa | grep gdb

- 如果存在则删除  
rpm -e --nodeps gdb-7.2-60.el6_4.1.i686  

- 安装文件系统，用于展示输出信息  
yum install ncurses-devel  
yum install texinfo

- 下载  
wget http://mirrors.ustc.edu.cn/gnu/gdb/gdb-8.0.tar.gz

- 解压  
tar -zxvf gdb-8.0.tar.gz

- 安装  
cd gdb-8.0  
./configure
make
make install

- 测试
gdb -v  


[安装包](http://mirrors.ustc.edu.cn/gnu/gdb/)

### 功能

 - 启动你的程序，可以按照你的自定义的要求随心所欲的运行程序。  
 - 可让被调试的程序在你所指定的调置的断点处停住。（断点可以是条件表达式）  
 - 当程序被停住时，可以检查此时你的程序中所发生的事。  
 - 你可以改变你的程序，将一个BUG产生的影响修正从而测试其他BUG。  


### 进入和离开

1、输入gdb开始GDB  
2、输入quit或者ctrl+d退出GDB

##### 1、常见启动方式  
gdb program//调用一个可自行的程序  
gdb program core//调用一个可自行的程序和core  
gdb program 1234//调用一个调试的进程

##### 2、Shell命令  
可以在gdb交互中直接执行shell命令而不用退出gdb去执行  
shell command string

##### 3、日志输出
set logging on //开启日志  
set logging off //关闭日志  
set logging file //修改输出日志文件  
show logging //展示目前日志设置状态  
set logging overwrite [on|off]//以覆盖形式输出日志

### 调试命令

###### help
展示某个分类下的帮助  
help status

###### complete
列出所有可能的不全结果（与tab补全类似）  
complete i

    //补全i开头的指令
    (gdb) complete i
    if
    ignore
    inferior
    info
    init-if-undefined
    interpreter-exec
    interrupt


###### info
查询程序的状态或者GDB本身的状态  
info breakpoints  

    //列出程序所有的断点
    (gdb) info breakpoints
    No breakpoints or watchpoints.

    //查看当前函数的参数值
    info args

    //查看当前函数栈上的信息
    info locals

    //查看寄存器的值
    info registers

###### set
设置环境变量的值

###### show
描述GDB本身的状态
show radix     

    //显示数值进制
    (gdb) show radix
    Input and output radices set to decimal 10, hex a, octal 12.

### GDB中运行程序

###### run(r)
开始一段程序

    (gdb) run
    Starting program: /data/script/c/heoo.out

    Breakpoint 1, main () at heoo.c:11
    11	    int n = 1;


###### attach
调试一个已经在运行的进程  
attach process-id  

###### detach
释放GDB对进程的控制；detach进城后，进程继续执行，进程与GDB就没有关系了

###### break(b)  
添加断点

b 函数名 在某函数入口处添加断点；  
b 行号 在指定行添加断点  
b 文件名：行号 在指定文件的指定行添加断点   
b 行号 if 条件 当条件为真时，指定行号处断点生效，例 b 5 if i=10 ，当i=10时，第5行断点生效

    (gdb) info breakpoints
    No breakpoints or watchpoints.
    (gdb) b main
    Breakpoint 1 at 0x400554: file heoo.c, line 11.
    (gdb) info breakpoints
    Num     Type           Disp Enb Address            What
    1       breakpoint     keep y   0x0000000000400554 in main
                                                       at heoo.c:11
    (gdb) b 13
    Breakpoint 2 at 0x40055b: file heoo.c, line 13.
    (gdb) b _add
    Breakpoint 3 at 0x40052b: file heoo.c, line 6.
    (gdb) b 7
    Breakpoint 4 at 0x400542: file heoo.c, line 7.


###### list(l)  
查看源码

    (gdb) list
    1	#include <stdio.h>
    2
    3	int g_var = 0;
    4
    5	static int _add(int a, int b) {
    6	    printf("_add callad, a:%d, b:%d\n", a, b);
    7	    return a+b;
    8	}
    9
    10	int main(void) {
    (gdb) list
    11	    int n = 1;
    12
    13	    printf("one n=%d, g_var=%d\n", n, g_var);
    14	    ++n;
    15	    --n;
    16
    17	    g_var += 20;
    18	    g_var -= 10;
    19	    n = _add(1, g_var);
    20	    printf("two n=%d, g_var=%d\n", n, g_var);
    (gdb) list
    21
    22	    return 0;
    23	}


###### next(n)
下一步

    (gdb) n
    17	    g_var += 20;
    (gdb) n
    18	    g_var -= 10;
    (gdb) n
    19	    n = _add(1, g_var);
    (gdb) n
    new-ui           nexti            nosharedlibrary


###### step(s)
下一步，可能会进入子函数

    (gdb) step

    Breakpoint 3, _add (a=1, b=10) at heoo.c:6
    6	    printf("_add callad, a:%d, b:%d\n", a, b);
    (gdb) step
    _add callad, a:1, b:10

    Breakpoint 4, _add (a=1, b=10) at heoo.c:7
    7	    return a+b;


###### print(p)
输出数据

###### continue(c)
继续后面的程序

    (gdb) continue
    Continuing.
    two n=11, g_var=10
    [Inferior 1 (process 18878) exited normally]


### 问题

###### Missing separate debuginfos?

    Missing separate debuginfos, use: debuginfo-install cyrus-sasl-lib-2.1.26-21.el7.x86_64 freetype-2.4.11-15.el7.x86_64 glibc-2.17-196.el7.x86_64 keyutils-libs-1.5.8-3.el7.x86_64 krb5-libs-1.15.1-8.el7.x86_64 libcom_err-1.42.9-10.el7.x86_64 libcurl-7.29.0-42.el7_4.1.x86_64 libgcrypt-1.5.3-14.el7.x86_64 libgpg-error-1.12-3.el7.x86_64 libidn-1.28-4.el7.x86_64 libpng-1.5.13-7.el7_2.x86_64 libselinux-2.5-12.el7.x86_64 libssh2-1.4.3-10.el7_2.1.x86_64 libxml2-2.9.1-6.el7_2.3.x86_64 libxslt-1.1.28-5.el7.x86_64 nspr-4.13.1-1.0.el7_3.x86_64 nss-3.28.4-12.el7_4.x86_64 nss-softokn-freebl-3.28.3-8.el7_4.x86_64 nss-util-3.28.4-3.el7.x86_64 openldap-2.4.44-5.el7.x86_64 openssl-libs-1.0.2k-8.el7.x86_64 pcre-8.32-17.el7.x86_64 xz-libs-5.2.2-1.el7.x86_64 zlib-1.2.7-17.el7.x86_64

*解决方案*  
debuginfo-install

###### debuginfo-install: command not found

*解决方案*  
yum update
yum install yum-utils

### 参考

 - [手把手教你玩GDB](https://www.cnblogs.com/skyofbitbit/p/3672848.html)
 - [gdb 调试PHP](https://phpor.net/blog/post/997)
 - [使用GDB调试PHP代码，解决PHP代码死循环](http://rango.swoole.com/archives/325)
 - [30分钟gdb调试快速突破](https://www.cnblogs.com/life2refuel/p/5396538.html)
