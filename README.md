# Smart Game Server Engine

> base on github.com/bytedance/netpoll

## Message Protocol
协议编码
0 1 2 3 4 5 6 7 8 9 a b c d e f 0 1 2 3 4 5 6 7 8 9 a b c d e f
+----------------------------------------------------------------+
| 0|                          LENGTH                             |
+----------------------------------------------------------------+
| 0|       HEADER MAGIC          |            FLAGS              |
+----------------------------------------------------------------+
|                         SEQUENCE NUMBER                        |
+----------------------------------------------------------------+
| 0|     HEADER SIZE        | ...
+---------------------------------

                  Header is of variable size:

                   (and starts at offset 14)

+----------------------------------------------------------------+
| PROTOCOL ID  |NUM TRANSFORMS . |TRANSFORM 0 ID (uint8)|
+----------------------------------------------------------------+
|  TRANSFORM 0 DATA ...
+----------------------------------------------------------------+
|         ...                              ...                   |
+----------------------------------------------------------------+
| INFO 0 ID (uint8)|       INFO 0  DATA ...
+----------------------------------------------------------------+
|         ...                              ...                   |
+----------------------------------------------------------------+
|                                                                |
|                              PAYLOAD                           |
|                                                                |
+----------------------------------------------------------------+
其中：

- LENGTH 字段 32bits，包括数据包剩余部分的字节大小，不包含 LENGTH 自身长度
- HEADER MAGIC 字段 16bits，值为：0x1000，用于标识 TTHeaderTransport
- FLAGS 字段 16bits，为预留字段，暂未使用，默认值为 0x0000
- SEQUENCE NUMBER 字段 32bits，表示数据包的 seqId，可用于多路复用，最好确保单个连接内递增
- HEADER SIZE 字段 16bits，等于头部长度字节数 /4，头部长度计算从第 14 个字节开始计算，一直到 PAYLOAD 前（备注：header 的最大长度为 64K）
- PROTOCOL ID 字段 uint8 编码，取值有：
    * ProtocolIDBinary = 0
    * ProtocolIDCompact = 2
- NUM TRANSFORMS 字段 uint8 编码，表示 TRANSFORM 个数
- TRANSFORM ID 字段 uint8 编码，具体取值参考下文
- INFO ID 字段 uint8 编码，具体取值参考下文
- PAYLOAD 消息内容
- PADDING 填充
- Header 部分长度 bytes 数必须是 4 的倍数，不足部分用 0x00 填充
- Transform IDs
    > 表示压缩方式，为预留字段，暂不支持，取值有：
    - ZLIB_TRANSFORM = 0x01，对应的 data 为空，表示用 zlib 压缩数据；
    - SNAPPY_TRANSFORM = 0x03，对应的 data 为空，表示用 snappy 压缩数据；
