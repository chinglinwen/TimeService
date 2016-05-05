# TimeService

提供单一时间源的时间服务

## 运行

```
cd bin
./timeServer
```

## 服务访问
可采用grpc支持的任一语言编写
以下通过go语言实现的客户端进行示例

```
cd bin
./timeClient
```

## 示例输出
```
# ./timeClient 
Time is: 2016-05-05 14:14:45.241
# 
```

更多信息可参考源代码
