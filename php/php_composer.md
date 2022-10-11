### composer包依赖管理工具

#### 可以解决的问题
>  
1、你有一个项目依赖于若干个库。  
2、其中一些库依赖于其他库。  
3、你声明你所依赖的东西。  
4、Composer 会找出哪个版本的包需要安装，并安装它们（将它们下载到你的项目中）。

#### composer.json组成
文件包含项目的依赖和一些原数据

##### 1、require
- 指定项目的依赖包  
  指定了项目的名称“monolog/monolog”及版本1.0.*  
>  
{ "require": {"monolog/monolog": "1.0.*" } }  

- 包名：monolog/monolog  
 包名称由供应商名称和其项目名称构成

- 版本1.0.*  
 指定依赖项目的版本，实际开发中可以指定项目的tag(oyj_1)或者开发分支（dev-master）

##### 2、repositories
- 指定包的来源  
  可以指定packagist、github、私有仓库  
>"repositories": {"packagist": {"type": "composer","url": "https://packagist.phpcomposer.com"}}


##### 3、autoload
- 指定自动加载的方式  
  指定psr-4自定加载规则
>"autoload":{"psr-4":{"Monolog\\" : "./src/"}}

- 自动加载
>在使用composer 包依赖管理工具安装后会按指定规则生成一个autoload.php，  
在项目的入口文件引入即可使用相应安装的包require 'vendor/autoload.php';

##### 4、name
- 指定包的名称
>"name":"oyjjpp/study",

#### composer命令行使用

##### 1、init
php composer init  
以交互的方式进行初始化

##### 2、install
php composer install  
从当前目录读取composer.json文件，处理依赖关系；   
如果当前目录下存在composer.lock，该命令会从此文件读取，确保每个使用者下载的依赖相同。

##### 3、update
php composer update  
获取依赖的最新版本，并且升级 composer.lock 文件

##### 4、require
php composer.phar require  
php composer.phar require vendor/package:2.* vendor/package2:dev-master

##### 5、search
php composer.phar search monolog  
为当前项目搜索依赖包，通常它只搜索 packagist.org 上的包

##### 6、show
php composer.phar show
//展示可用的包
php composer.phar show monolog/monolog
//展示指定包的信息  
php composer show --platform  
//平台软件包的列表

##### 7、self-update
php composer.phar slef-update  
对composer自身升级

##### 8、config
php composer.phat config key value  
可以对composer.json文件进行编辑，针对指定key编辑对应value    
php composer.phat config --list  
展示当前配置的选项列表

##### 9、dump-autoload
php composer.phat dump-autoload  
更新 autoloader

##### 10、licenses
php composer.phat licenses  
列出已安装的每个包的名称、版本、许可协议

##### 11、环境变量
###### COMPOSER
COMPOSER  可以为 composer.json 文件指定其它的文件名（适应不同的环境开发）  
COMPOSER=composer-pro.json composer install  

###### COMPOSER_ROOT_VERSION
COMPOSER_ROOT_VERSION  指定 root 包的版本

###### COMPOSER_VENDOR_DIR
COMPOSER_VENDOR_DIR 指定 composer 将依赖安装在 vendor 以外的其它目录中  
与config下的vendor-dir作用一致

#### 参考
 - [Composer文档](https://docs.phpcomposer.com/)
