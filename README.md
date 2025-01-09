# Smart Game Server Engine

> base on github.com/bytedance/netpoll

## Protocol 协议编码

0 1 2 3 4 5 6 7 8 9 a b c d e f 0 1 2 3 4 5 6 7 8 9 a b c d e f
+----------------------------------------------------------------+
| 0|                          LENGTH                             |
+----------------------------------------------------------------+
| 0|       HEADER MAGIC       |   COMPRESS      |      FLAGS     |
+----------------------------------------------------------------+
|                                                                |
|                              BODY                              |
|                                                                |
+----------------------------------------------------------------+
其中：

- LENGTH 字段 32bits，BODY 部分的字节大小
- HEADER MAGIC 字段 16bits，值为：0x0000，用于标识 协议(SmartProtocol)
- FLAGS 字段 8bits，为预留字段，暂未使用，默认值为 0x0000
- BODy 消息内容, message.ProtocolMessage 的字节数组
- COMPRESS 字段 8bits，BODY 压缩算法，暂未实现，默认值为 0x0000
    - ZLIB = 0x01，表示用 zlib 压缩数据；
    - SNAPPY = 0x02，表示用 snappy 压缩数据；
