### Linux下Git的安装

>下载目录  
https://www.kernel.org/pub/software/scm/git/


>安装依赖库  
yum install curl-devel expat-devel gettext-devel \
  openssl-devel zlib-devel  
wget https://www.kernel.org/pub/software/scm/git/git-2.9.3.tar.gz  
tat -zxf git-2.9.3.tar.gz  
make prefix=/usr/local all  
sudo make prefix=/usr/local install


### 参考
 - [沉浸式学Git]http://igit.linuxtoy.org/contents.html
 - [git]https://git-scm.com/
