### PHP扩展编写

### 1 为什么编写PHP扩展
使用php代码编写 请求一个数组元素的平方和

    <?php
        function array_square_sum($data){
            $sum = 0;
            foreach($data as $value){
                $sum += $value * $value;
            }
            return $sum;
        }

实际执行的时候，php zend引擎会把这段话翻译成C语言，每次都需要进行内存分配。所以性能比较差。进而，基于性能上的考虑，我们可以编写一个扩展来做这个事情。

### 2 扩展编写 
> 一个是Configulator文件，它会告诉编译器编译这个扩展至少需要哪些依赖库；第二个是实际执行的文件  
> 第二个是实际执行的文件  

#### 2.1 通过源码中的 ext_skel 可以生成扩展框架
PHP 扩展由几个文件组成，这些文件对所有扩展来说都是通用的。不同扩展之间，这些文件的很多细节是相似的，只是要费力去复制每个文件的内容。幸运的是，有脚本可以做所有的初始化工作，名为 ext_skel，自 PHP 4.0 起与其一起分发。  
位置：/php-src/ext/ext_skel

    [root@iZ2ze1yhgn9t43zixdpjvcZ ext]# ./ext_skel 
    ./ext_skel --extname=module [--proto=file] [--stubs=file] [--xml[=file]]
               [--skel=dir] [--full-xml] [--no-help]

      --extname=module   module is the name of your extension
      --proto=file       file contains prototypes of functions to create
      --stubs=file       generate only function stubs in file
      --xml              generate xml documentation to be added to phpdoc-svn
      --skel=dir         path to the skeleton directory
      --full-xml         generate xml documentation for a self-contained extension
                         (not yet implemented)
      --no-help          don't try to be nice and create comments in the code
                         and helper functions to test if the module compiled


#### 2.2 执行扩展
    $ ./ext_skel --extname=array_square_sum
    Creating directory array_square_sum
    Creating basic files: config.m4 config.w32 .svnignore array_square_sum.c php_array_square_sum.h CREDITS EXPERIMENTAL tests/001.phpt array_square_sum.php [done].

    To use your new extension, you will have to execute the following steps:

    1.  $ cd ..
    2.  $ vi ext/array_square_sum/config.m4
    3.  $ ./buildconf
    4.  $ ./configure --[with|enable]-array_square_sum
    5.  $ make
    6.  $ ./php -f ext/array_square_sum/array_square_sum.php
    7.  $ vi ext/array_square_sum/array_square_sum.c
    8.  $ make

    Repeat steps 3-6 until you are satisfied with ext/array_square_sum/config.m4 and
    step 6 confirms that your module is compiled into PHP. Then, start writing
    code and repeat the last two steps as often as necessary.

#### 2.3 执行命令后生成的文件
执行命令之后，终端告诉我们怎么去生产新的扩展。查看一下文件内容，会发现多了几个比较重要的文件config.m4,php_array_square_sum.h，array_square_sum.c。

    [root@iZ2ze1yhgn9t43zixdpjvcZ array_square_sum]$ ll
    total 32
    -rw-r--r-- 1 root root 5582 May 11 08:26 array_square_sum.c
    -rw-r--r-- 1 root root  532 May 11 08:26 array_square_sum.php
    -rw-r--r-- 1 root root 2408 May 11 08:26 config.m4
    -rw-r--r-- 1 root root  425 May 11 08:26 config.w32
    -rw-r--r-- 1 root root   16 May 11 08:26 CREDITS
    -rw-r--r-- 1 root root    0 May 11 08:26 EXPERIMENTAL
    -rw-r--r-- 1 root root 2476 May 11 08:26 php_array_square_sum.h
    drwxr-xr-x 2 root root 4096 May 11 08:26 tests

##### 2.3.1 配置文件config.m4
> 1、告诉unix系统哪些扩展configure选项是支持的？  
> 2、需要哪些扩展库？  
> 3、哪些源文件要编译成它的一部分？  

config.m4 文件使用 GNU autoconf 语法编写  
注释：dnf  
字符串:[]  

*改动*
    
    dnl PHP_ARG_WITH(array_square_sum, for array_square_sum support,
    dnl Make sure that the comment is aligned:
    dnl [  --with-array_square_sum             Include array_square_sum support])
    
    修改之后  
    
    PHP_ARG_WITH(array_square_sum, for array_square_sum support,
    Make sure that the comment is aligned:
    [  --with-array_square_sum             Include array_square_sum support])

*注意*
在调用 configure 时，不管选项在命令行中的顺序如何，都会按在 config.m4 中指定的顺序进行检测。


##### 2.3.2 php_array_square_sum.h
当将扩展作为静态模块构建并放入 PHP 二进制包时，构建系统要求用 php_ name（扩展的名称命名）的头文件包含一个对扩展模块结构的指针定义。就象其他头文件，此文件经常包含附加的宏、原型和全局量。

*改动*  
添加 PHP_FUNCTION(array_square_sum);  

    PHP_FUNCTION(array_square_sum);
    
##### 2.3.3 array_square_sum.c
主要的扩展源文件。按惯例，此文件名就是扩展的名称，但不是必需的。此文件包含模块结构定义、INI 条目、管理函数、用户空间函数和其它扩展所需的内容。  

添加以下函数  

    const zend_function_entry array_square_sum_functions[] = { 
        PHP_FE(array_square_sum,    NULL)       /* For testing, remove later. */
        PHP_FE_END  /* Must be the last line in array_square_sum_functions[] */
    }; 

    PHP_FUNCTION(array_square_sum)
    {
        zval* array_data;
        HashTable *ht_data;
        int ret;
        char* key;
        uint index;
        zval **pdata;
        double sum = 0;

        if (zend_parse_parameters(ZEND_NUM_ARGS() TSRMLS_CC, "a", &array_data) == FAILURE) {
            return;
        }   

        ht_data = Z_ARRVAL_P(array_data);
        zend_hash_internal_pointer_reset(ht_data);
        while ( HASH_KEY_NON_EXISTANT != (ret = zend_hash_get_current_key(ht_data, &key, &index, 0)) ) { 
            ret = zend_hash_get_current_data(ht_data, &pdata);
        
            if( Z_TYPE_P(*pdata) == IS_LONG){
                sum +=  Z_LVAL_P(*pdata) *  Z_LVAL_P(*pdata);
            }else {
                RETURN_FALSE;
            }   
            zend_hash_move_forward(ht_data);
        }   
        zend_hash_internal_pointer_end(Z_ARRVAL_P(array_data));
        RETVAL_DOUBLE(sum);
    } 

#### 3 安装扩展
    
### 参考
> [一步步入门编写PHP扩展][1]  
> [深入 理解PHP内核][2]  
> [与 UNIX 构建系统交互: config.m4][3]  

[1]: http://www.open-open.com/lib/view/open1392188698114.html  
[2]: http://www.php-internals.com/book/?p=index  
[3]: http://php.net/manual/zh/internals2.buildsys.configunix.php
