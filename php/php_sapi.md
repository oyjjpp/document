### 问题（带着疑问学习）

    1. 源码的sapi目录
    2. 什么是SAPI
    3. 作用
    4. 安全、性能、效率、易用性透析不同的SAPI


### 1、环境

    [root@iZ2ze1yhgn9t43zixdpjvcZ sapi]# php -v
    PHP 7.1.1 (cli) (built: May 11 2017 07:45:21) ( NTS )
    Copyright (c) 1997-2017 The PHP Group
    Zend Engine v3.1.0, Copyright (c) 1998-2017 Zend Technologies

### 2、源码的sapi目录
sapi 包含了各种服务器抽象层的代码，例如apache的mod_php，cgi，fastcgi以及fpm等等接口

    [root@iZ2ze1yhgn9t43zixdpjvcZ php-src]# ll sapi/
    total 32
    drwxr-xr-x 2 root root 4096 May 10 08:30 apache2handler
    drwxr-xr-x 4 root root 4096 Jun  9 08:17 cgi
    drwxr-xr-x 4 root root 4096 Jun  9 08:17 cli
    drwxr-xr-x 2 root root 4096 May  9 08:28 embed
    drwxr-xr-x 5 root root 4096 Jun  9 08:17 fpm
    drwxr-xr-x 2 root root 4096 May 10 08:30 litespeed
    drwxr-xr-x 5 root root 4096 May 21 21:01 phpdbg
    drwxr-xr-x 2 root root 4096 May  9 08:28 tests

### 3、简介SAPI
SAPI是Server Application Programming Interface（服务器应用编程接口）的缩写。PHP通过SAPI提供了一组接口，供其他应用和PHP内核之间进行数据交互。


### 4、PHP生命周期
模块初始化 Init（MINT ），  
激活请求Init（RINT），  
处理  
关闭请求ShutDown(RSHUTDOWN)，  
关闭模块Shutdown（MSHUTDOWN）  

### 5、常见的PHP函数查看sapi信息

//查看PHP一次生命周期加载的模块  
CLI：php -m  
get_loaded_extensions()  

//查看当前SAPI接口类型  
string php_sapi_name(void)  

### 6、安装PHP时指定SAPI
通过configure进行配置  
SAPI modules:  

    --with-apxs2=FILE       Build shared Apache 2.0 Handler module. FILE is the optional
                          pathname to the Apache apxs tool apxs
    --disable-cli           Disable building CLI version of PHP
                          (this forces --without-pear)
    --enable-embed=TYPE     EXPERIMENTAL: Enable building of embedded SAPI library
                          TYPE is either 'shared' or 'static'. TYPE=shared
    --enable-fpm            Enable building of the fpm SAPI executable
    --with-fpm-user=USER    Set the user for php-fpm to run as. (default: nobody)
    --with-fpm-group=GRP    Set the group for php-fpm to run as. For a system user, this
                          should usually be set to match the fpm username (default: nobody)
    --with-fpm-systemd      Activate systemd integration
    --with-fpm-acl          Use POSIX Access Control Lists
    --with-litespeed        Build PHP as litespeed module
    --enable-phpdbg            Build phpdbg
    --enable-phpdbg-webhelper  Build phpdbg web SAPI support
    --enable-phpdbg-debug      Build phpdbg in debug mode
    --disable-cgi           Disable building CGI version of PHP

### 7、按应用分类

嵌入式：embed  
命令行：cli  
web服务脚本：fpm、apache2handler、cgi  

### 8、定义SAPI、首先需要定义sapi_module_struct
每个SAPI实现都是一个_sapi_module_struct结构体变量，结构中包含 name 、pretty_name、等属性及 startup() shutdown()等函数

#### 8.1 main/SAPI.h 文件夹下
    218 _struct _sapi_module_struct {
    219     char *name;
    220     char *pretty_name;
    221 
    222     int (*startup)(struct _sapi_module_struct *sapi_module);
    223     int (*shutdown)(struct _sapi_module_struct *sapi_module);
    224 
    225     int (*activate)(void);
    226     int (*deactivate)(void);
    227 
    228     size_t (*ub_write)(const char *str, size_t str_length);
    229     void (*flush)(void *server_context);
    230     zend_stat_t *(*get_stat)(void);
    231     char *(*getenv)(char *name, size_t name_len);
    232 
    233     void (*sapi_error)(int type, const char *error_msg, ...) ZEND_ATTRIBUTE_FORMAT(printf, 2, 3);
    234 
    235     int (*header_handler)(sapi_header_struct *sapi_header, sapi_header_op_enum op, sapi_headers_struct *sapi_headers);
    236     int (*send_headers)(sapi_headers_struct *sapi_headers);
    237     void (*send_header)(sapi_header_struct *sapi_header, void *server_context);
    238 
    239     size_t (*read_post)(char *buffer, size_t count_bytes);
    240     char *(*read_cookies)(void);
    241 
    242     void (*register_server_variables)(zval *track_vars_array);
    243     void (*log_message)(char *message, int syslog_type_int);
    244     double (*get_request_time)(void);
    245     void (*terminate_process)(void);
    246 
    247     char *php_ini_path_override;
    248 
    249     void (*default_post_reader)(void);
    250     void (*treat_data)(int arg, char *str, zval *destArray);
    251     char *executable_location;
    252 
    253     int php_ini_ignore;
    254     int php_ini_ignore_cwd; /* don't look for php.ini in the current directory */
    255 
    256     int (*get_fd)(int *fd);
    257 
    258     int (*force_http_10)(void);
    259 
    260     int (*get_target_uid)(uid_t *);
    261     int (*get_target_gid)(gid_t *);
    262 
    263     unsigned int (*input_filter)(int arg, char *var, char **val, size_t val_len, size_t *new_val_len);
    264 
    265     void (*ini_defaults)(HashTable *configuration_hash);
    266     int phpinfo_as_text;
    267 
    268     char *ini_entries;
    269     const zend_function_entry *additional_functions;
    270     unsigned int (*input_filter_init)(void);
    271 };

#### 8.2 sapi/cgi/cgi_main.c
    static sapi_module_struct cgi_sapi_module = {
    933     "cgi-fcgi",                     /* name */
    934     "CGI/FastCGI",                  /* pretty name */
    935 
    936     php_cgi_startup,                /* startup */
    937     php_module_shutdown_wrapper,    /* shutdown */
    938 
    939     sapi_cgi_activate,              /* activate */
    940     sapi_cgi_deactivate,            /* deactivate */
    941 
    942     sapi_cgi_ub_write,              /* unbuffered write */
    943     sapi_cgi_flush,                 /* flush */
    944     NULL,                           /* get uid */
    945     sapi_cgi_getenv,                /* getenv */
    946 
    947     php_error,                      /* error handler */
    948 
    949     NULL,                           /* header handler */
    950     sapi_cgi_send_headers,          /* send headers handler */
    951     NULL,                           /* send header handler */
    952 
    953     sapi_cgi_read_post,             /* read POST data */
    954     sapi_cgi_read_cookies,          /* read Cookies */
    955 
    956     sapi_cgi_register_variables,    /* register server variables */
    957     sapi_cgi_log_message,           /* Log message */
    958     NULL,                           /* Get request time */
    959     NULL,                           /* Child terminate */
    960 
    961     STANDARD_SAPI_MODULE_PROPERTIES
    962 };


### 参考
> [深入PHP内核（二）——SAPI探究][1]  
> [深入理解PHP内核][2]  
> [深入理解Zend SAPIs(Zend SAPI Internals)][3]  

[1]: http://www.csdn.net/article/2014-09-26/2821885-exploring-of-the-php-2  
[2]: http://www.php-internals.com/book/?p=chapt02/02-02-00-overview  
[3]: http://www.laruence.com/2008/08/12/180.html