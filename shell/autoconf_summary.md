## autoconf总结

## autoconf 升级

一、查看autoconf路径

whereis autoconf

二、查看autoconf版本

rpm -qf /usr/bin/autoconf

三、卸载当前版本

rpm -e --nodeps autoconf-2.63-5.1.el6.noarch

四、下载新版本

1、路径http://ftp.gnu.org/gnu/autoconf/

2、挑选需要的版本下载

wget http://ftp.gnu.org/gnu/autoconf/autoconf-2.69.tar.gz

五、安装

./configure --prefix=/usr

make && make install

六、查看是否安装成功

/usr/bin/autoconf --help

/usr/bin/autoconf -V

### 参考

http://www.cnblogs.com/sunss/archive/2011/07/20/2111697.html




















