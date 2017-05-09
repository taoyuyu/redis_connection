#基于redigo的redis连接池

redis是单线程模型，使用redigo连接进行并发set，get操作时会报错：write tcp 127.0.0.1:51008->127.0.0.1:6379: use of closed network connection
因此对单连接set，get操作必须加锁

避免使用锁的方式可以使用redigo连接池，将同步操作交给redis底层实现

