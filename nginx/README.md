# nginx常见问题

## location的优先级是什么样的

精准匹配： 相等（=）  
字符串匹配：字符串匹配（空格）、匹配开头（^~）  
正则匹配： 区分大小写匹配（~）、不区分大小写匹配（~*）、区分大小写不匹配（！~）、不区分大小写不匹配（！~*）  

优先级：  
精准匹配  > 字符串匹配（长 > 短，^~匹配是最长匹配则停止匹配） > 正则匹配（先后顺序）

精准匹配只能匹配一个  
字符串匹配使用匹配最长的为匹配结果  
正则匹配按照location定义的顺序进行匹配，先定义具有最高优先级  

注：字符串匹配优先搜索，但是只是记录下最长的匹配（如果^~是最长匹配，则会直接命中，停止搜索正则），然后继续正则匹配，如果有正则匹配，则命中正则匹配，如果没有正则匹配，则命中最长的字符串匹配

## nginx负载均衡策略

nginx负载均衡是一种将网络请求分配到多个服务器上的方法，以提高性能和可靠性。在nginx中，负载均衡可以使用不同的策略来实现。

下面是一些常用的nginx负载均衡策略：

轮询（Round Robin）：默认情况下，nginx使用轮询策略。它会将每个新的请求依次发送到下一个可用的服务器上，以确保所有服务器获得相同数量的请求。这是一种简单而均衡的负载均衡策略。

IP哈希（IP Hash）：这种策略基于客户端的IP地址来进行负载均衡。nginx会根据客户端的IP地址计算一个哈希值，然后将请求发送到相应的服务器上。这种策略可以确保相同IP地址的客户端总是被发送到同一个服务器上，这对于某些应用程序可能是必要的。

最小连接数（Least Connections）：该策略会将请求发送到当前连接数最少的服务器上。这种策略可以确保服务器的负载更均衡，因为它总是选择连接数最少的服务器。

通配符（Wildcard）：该策略会根据请求的URL匹配一个通配符表达式，然后将请求发送到匹配的服务器上。这种策略可以根据请求的特定属性来选择服务器，以确保更好的性能和可靠性。

这些策略可以单独使用或组合使用，以满足特定应用程序的需求。例如，您可以使用IP哈希和最小连接数组合来确保客户端总是被发送到连接数最少的服务器上。

## nginx用过吗？

## 大致了解nginx的哪些功能？

## nginx的负载均衡是在第几层？

## 除了nginx的负载均衡还了解过其他负载均衡吗？

## 反向代理和正向代理有什么差别吗？

## 如何统计nginx日志里面的访问量最多的十个IP地址？
