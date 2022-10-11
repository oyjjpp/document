## 问题：

    1. 扩展的运行原理
    2. 扩展的作用
    3. php源码安装后php.ini 找不到相应扩展的配置？
    4. phpize工具是什么？
    5. 如何加载原有扩展？
    6. php7.1.5 编写扩展后 调用函数找不到？



### mbstring通过phpize安装

    cd /usr/local/src/php-7.1.6/ext/mbstring  
    ./configure --with-php-config=/usr/local/server/php/bin/php-config  
    make  
    make install  
    生成mbstring.so 扩展文件  
    copy到php扩展目录下 extension_dir  
    修改php.ini 添加extension = mbstring.so  
    重启nginx和php-fpm

### 通过ext_skel脚本创建say扩展

#### 1、创建say扩展目录
    ./ext_skel --extname=say  
注释：不会使用ext_skel 可以ext_skel --help 查看相关帮助文档

#### 2、进入目录中
    cd say

#### 3、修改config.m4文件
注释：针对扩展的congfigure配置项 在这里设置

    vim config.m4  

    //源文件  
    dnl   PHP_ARG_ENABLE(say, whether to enable say support,  
    dnl   Make sure that the comment is aligned:  
    dnl   [  --enable-say           Enable say support])  

    //修改后的文件
    PHP_ARG_ENABLE(say, whether to enable say support,  
    Make sure that the comment is aligned:  
    [  --enable-say           Enable say support])  

#### 4、创建sayhello()函数
vim say.c

    添加函数  
    PHP_FUNCTION(sayhello)
    {
        char *arg = NULL;
        size_t arg_len, len;
        zend_string *strg;

        if (zend_parse_parameters(ZEND_NUM_ARGS(), "s", &arg, &arg_len) == FAILURE) {
            return;
        }

        strg = strpprintf(0, "hello word");

        RETURN_STR(strg);
    }

    修改  
    const zend_function_entry say_functions[] = {
        PHP_FE(sayhello, NULL)
        PHP_FE(confirm_say_compiled,    NULL)       /* For testing, remove later. */
        PHP_FE_END  /* Must be the last line in say_functions[] */
    }; 

#### 5、在say目录下执行phpize
    [root@iZ2ze1yhgn9t43zixdpjvcZ say]# phpize
    Configuring for:
    PHP Api Version:         20160303
    Zend Module Api No:      20160303
    Zend Extension Api No:   320160303

#### 6、执行configure
    [root@iZ2ze1yhgn9t43zixdpjvcZ say]# ./configure --with-php-config=/usr/local/server/php/bin/php-config

#### 7、make 

#### 8、make install
    [root@iZ2ze1yhgn9t43zixdpjvcZ say]# make install
    Installing shared extensions:     /usr/local/server/php/lib/php/extensions/no-debug-non-zts-20160303/

> **以上为创建say扩展的过程**

#### 9、加载到PHP中

    say.so copy到php扩展目录下 extension_dir  
    修改php.ini 添加extension = say.so  
    重启nginx和php-fpm

#### php版本信息
    [root@iZ2ze1yhgn9t43zixdpjvcZ say]# php -v
    PHP 7.1.6 (cli) (built: Jun 26 2017 07:32:45) ( NTS )
    Copyright (c) 1997-2017 The PHP Group
    Zend Engine v3.1.0, Copyright (c) 1998-2017 Zend Technologies


#### 通过程序加以验证
    vim say.php
    <?php
    $str = sayhello('ouyangjun');
    var_dump($str);

#### 最终输出  
> string(10) "hello word"
    
### 参考
> [PHP扩展开发与内核应用][1]  
> [PHP扩展开发与内核应用][2]  

[1]: https://github.com/walu/phpbook  
[2]: http://www.cunmou.com/phpbook/preface.md  

