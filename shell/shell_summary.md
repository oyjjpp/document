### 常用Shell命令总结

#### 包管理工具

###### 查看包的版本
>rpm -qf /usr/bin/autoconf

###### 查看软件路径
>whereis autoconf  


#### 文件、目录、内容查找

###### 只显示当前目录下的文件
>ls -al | grep "^-"

###### 查看指定字符串
>grep  -ri "string" *


#### 网络相关

###### 查看端口使用
>netstat -anl | grep "80"
