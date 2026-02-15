# 客户端行为测试

### 1. 编译测试工具

安装 Go 工具链, 确保系统中已安装 `Wireshark` 软件用于后期分析

```shell
go build -ldflags="-s -w -checklinkname=0" .
```

### 2. 流量捕获

运行软件, 并使用 `tcpdump` 或 `Wireshark` 实时捕获数据包

```shell
# tcpdump -i <网卡名称> host <客户端地址> -w <捕获文件.cap>

# 示例: 捕获 eth0 网卡上来自 127.0.0.1 的流量
tcpdump -i eth0 host 127.0.0.1 -w eth0.cap
```

### 3. 执行测试连接

使用 SSL/TLS 客户端尝试连接到服务器配置的每一个测试端口, 确保触发握手

### 4. 导入捕获数据

测试完成后停止捕获, 使用 `Wireshark` 打开生成的 `.cap` 文件

### 5. 配置 TLS 解密首选项

为了查看加密的内容, 需要设置 Wireshark 解密选项

![Wireshark - Open TLS Option](../../assets/guide_wireshark_open_tls_option.png)

### 6. 设置密钥日志文件

![Wireshark - Set Key Log File](../../assets/guide_wireshark_set_key_log_file.png)

### 7. 行为分析与数据统计

通过 Wireshark 的过滤功能和统计工具分析以下内容

- 计数器
- 响应数据
- 响应行为