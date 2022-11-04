# 网络

## 输入URL到浏览器发生了什么？

[键入网址到网页显示，期间发生了什么？](https://xiaolincoding.com/network/1_base/what_happen_url.html)

```network
1、浏览器会调用网络组件，通过DNS获取当前域名的IP地址  
2、拿到IP地址，与服务段建立TCP链接  
3、建立链接之后向服务端发送请求数据  
4、服务端处理数据之后返回给浏览器
5、浏览器针对相应数据进行渲染
```

## DNS基于TCP还是UDP，为什么？

```network
同时使用TCP和UDP协议;
DNS在区域传输的时候使用TCP协议，为了保证数据传输时的准确性；其他时候使用UDP协议，保证快速的域名解析
```

## 域名解析的流程

![image](./image/domain-name-resolution.gif)

- 1、客户端通过浏览器访问域名为 www.baidu.com (<http://www.baidu.com>) 的网站，发起查询该域名的IP地址的DNS请求。该请求发送到了本地DNS服务器上。本地DNS服务器会首先查询它的缓存记录，如果缓存中有此条记录，就可以直接返回结果。如果没有，本地DNS服务器还要向DNS根服务器进行查询。  
- 2、本地DNS服务器向根服务器发送DNS请求，请求域名为 www.baidu.com (<http://www.baidu.com>)的IP地址。  
- 3、根服务器经过查询，没有记录该域名及IP地址的对应关系。但是会告诉本地DNS服务器，可以到域名服务器上继续查询，并给出域名服务器的地址（.com 服务器）。  

- 4、本地 DNS 服务器向 .com 服务器发送 DNS 请求，请求域名www.baidu.com (<http://www.baidu.com>) 的 IP 地址。  
- 5、.com 服务器收到请求后，不会直接返回域名和IP地址的对应关系，而是告诉本地 DNS 服务器，该域名可以在 baidu.com 域名服务器上进行解析获取 IP 地址，并告诉 baidu.com 域名服务器的地址。  
- 6、本地 DNS 服务器向 baidu.com 域名服务器发送 DNS 请求，请求域名 www.baidu.com (<http://www.baidu.com>) 的 IP 地址。  
- 7、baidu.com 服务器收到请求后，在自己的缓存表中发现了该域名和 IP 地址的对应关系，并将IP地址返回给本地 DNS 服务器。  
- 8、本地 DNS 服务器将获取到与域名对应的 IP 地址返回给客户端，并且将域名和 IP 地址的对应关系保存在缓存中，以备下次别的用户查询时使用。

## 描述session和cookie

### 相同之处

- cookie和session都是用来跟踪浏览器用户身份的会话方式。

### 区别之处

- 1、cookie数据存放在客户的浏览器上，session数据放在服务器上；  
- 2、cookie不是很安全，别人可以分析存放在本地的COOKIE并进行COOKIE欺骗,如果主要考虑到安全应当使用session；  
- 3、session会在一定时间内保存在服务器上。当访问增多，会比较占用你服务器的性能，如果主要考虑到减轻服务器性能方面，应当使用cookie；  
- 4、cookie有大小限制

### cookie传输流程

- 1、由客户端首次请求时服务端生成  
- 2、服务端生成后，随着响应返回给客户端  
- 3、之后的每次客户端与服务端交互中携带cookie  

## 如何告诉客户端，响应的是什么数据类型？

```network
通过http的响应头Content-Type告知客户端响应的数类型
```

## 网络分层(OSI)

[网络模型](https://xiaolincoding.com/network/1_base/tcp_ip_model.html)

```network
应用层-》传输层-》网络层-》链路层-》物理层

应用层：HTTP、FTP、SMTP、DNS  
传输层：TCP、UDP  
网络层：IP  
链路层：DOCSIS  
物理层：网卡
```

![image](./image/202211011935.jpeg)

## TCP与UDP差异

- 1、TCP基于面向连接的协议，数据传输可靠，传输速度慢，适用于传输大量数据，可靠性要求高的场合。  
- 2、UDP协议面向非连接协议，数据传输不可靠，传输速度快，适用于一次只传送少量数据、对可靠性要求不高的应用环境。  

## io多路复用，select\poll\epoll的实现和区别

[IO 多路复用是什么意思？](https://www.zhihu.com/question/32163005/answer/1802684879)

## 网络协议

## 网络安全

## http和tcp有什么区别

## http和https的区别

## http和http2区别

## 用netstat看tcp连接的时候有关注过time_wait和close_wait吗？

## tcp如何保证稳定性

## https的连接过程

## 三次握手和四次挥手

## 长连接和短链接(怎么实现的、区别以及应用场景)

## 网络七层模型和五层模型

## 输入URL发生的整个网络过程

## TCP流量控制、拥塞控制

## TCP半连接队列

## TCP半关闭状态

## TCP滑动窗口

## tcp与udp区别，udp优点，适用场景

## tcp握手，为什么是四次挥手

## 挥手time_wait阶段是什么时候

## time_wait时间太长会怎样？

## http restful的定义规范

## 常见的http状态码

## tcp和udp有什么区别，在基础时间和传输层有什么区别

## tcp三次握手四次挥手

## tcp和udp区别

## 对连接池有了解吗？

## 信令用wss还是ws？

## wss是基于tcp的，tcp有个半连接队列，有没有遇到发了信令但是服务器没收到的情况？

## ntp协议

## 时间同步问题？

## tcp如何实现稳定性

## https如何实现的，证书是怎么获取的？

## 服务器上有大量time_wait状态，如何分析一下为什么会这样？

## 说下websocket是哪一层协议？

## 常见http状态码
