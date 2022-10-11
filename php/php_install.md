
### php源码安装 configure
./configure --prefix=/usr/local/server/php \
--with-config-file-path=/usr/local/server/php/etc \
--enable-inline-optimization \
--disable-debug --disable-rpath \
--enable-shared \
--enable-opcache \
--enable-fpm \
--with-fpm-user=www \
--with-fpm-group=www \
--with-mysqli=mysqlnd \
--with-pdo-mysql=mysqlnd \
--with-gettext \
--enable-mbstring \
--with-iconv \
--with-mcrypt \
--with-mhash \
--with-openssl \
--enable-bcmath \
--enable-soap \
--with-libxml-dir \
--enable-pcntl \
--enable-shmop \
--enable-sysvmsg \
--enable-sysvsem \
--enable-sysvshm \
--enable-sockets \
--with-curl \
--with-zlib \
--enable-zip \
--with-bz2 \
--with-readline \
--without-sqlite3 \
--without-pdo-sqlite \
--with-pear


### 添加到环境变量
vim /etc/profile

###### 在文件结尾处添加
$PATH=$PATH:/usr/local/server/php
export PATH

###### 更改立即生效
source /etc/profile

### 问题

1、nginx FastCGI错误Primary script unknown解决办法

###### nginx.config 配置
location / {
   root   /usr/local/nginx/html; \
   index  index.php index.html index.htm; \
   fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name; \
}

## 参考

1、[PHP编译安装时常见错误解决办法](http://www.poluoluo.com/jzxy/201505/364819.html)

2、[nginx安装php和php-fpm](http://www.cnblogs.com/zhja/p/3978870.html)