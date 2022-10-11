## Linux下Nginx的安装

### 安装编译环境

yum -y install gcc automake autoconf libtool make

### 下载地址：

wget http://nginx.org/download/nginx-1.10.1.tar.gz?_ga=2.158863866.1516274824.1494978789-1527359295.1494947241

### 下载相应依赖包

1、安装PCRE库

包路径 ftp://ftp.csx.cam.ac.uk/pub/software/programming/pcre/

wget ftp://ftp.csx.cam.ac.uk/pub/software/programming/pcre/pcre2-10.23.tar.gz

2、安装zlib库

包路径 http://zlib.net

wget http://zlib.net/zlib-1.2.11.tar.gz

cd /usr/local/src
wget http://nginx.org/download/nginx-1.4.2.tar.gz
tar -zxvf nginx-1.4.2.tar.gz
cd nginx-1.4.2

### 安装
    ./configure --sbin-path=/usr/local/server/nginx/nginx \
    --prefix=/usr/local/server/nginx \
    --modules-path=/usr/local/server/nginx/modules \
    --error-log-path=/usr/local/server/nginx/logs/error.log \
    --http-log-path=/usr/local/server/nginx/logs/access.log \
    --conf-path=/usr/local/server/nginx/nginx.conf \
    --pid-path=/usr/local/server/nginx/nginx.pid \
    --with-http_ssl_module \
    --with-pcre=/usr/local/src/pcre-8.39 \
    --with-zlib=/usr/local/src/zlib-1.2.11 \
    --with-openssl=/usr/local/src/openssl-1.0.1t


    make

    make install

    --with-pcre=/usr/local/src/pcre-8.39 \    指的是pcre-8.39 的源码路径。

    --with-zlib=/usr/local/src/zlib-1.2.11 \    指的是zlib-1.2.11 的源码路径。

    --with-openssl=/usr/local/src/openssl-1.0.1t   指的是openssl-1.0.1t 的源码路径。



### 参考
 - [nginx安装错误汇总](http://blog.csdn.net/qq_29461259/article/details/52609775)
 - [nginx的启动、停止与重启](http://www.cnblogs.com/codingcloud/p/5095066.html)
