### PHP7内核进阶之路

#### 测试

###### 1、PHP的纯CPU基准测试
 - 测试文件位置  
源码/Zend/bench.php,源码/Zend/micro_bench.php  

 - 测试环境  
>[root@Zend]$lsb_release -a  
LSB Version:	:core-4.1-amd64:core-4.1-noarch  
Distributor ID:	CentOS  
Description:	CentOS Linux release 7.4.1708 (Core)  
Release:	7.4.1708  
Codename:	Core  


 - 测试命令  
time php bench.php

 - 输出结果  
>real	0m0.665s  
user	0m0.653s  
sys	0m0.008s  
