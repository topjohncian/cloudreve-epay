# Cloudreve 易支付网关

需要更新到 Pro 3.7.1 才可正常使用

## 点点 Star 不迷路 ❤ 有问题请发 Issue 😭

## 注意事项

1. 一定更新到 Pro 3.7.1
2. 最好启用 Redis，否则使用内存缓存的话，一旦程序终止，支付将永远无法回调 （后续改进）
3. 已修复，请使用 `-eject` 参数导出模板 ~~目前支付模板是硬编码字符串拼接，可能造成 XSS （后续改进）~~
4. 已支持，请设置 `CR_EPAY_EPAY_PURCHASE_TYPE` ~~目前只会默认调用支付宝，建议选择有自己可以选择支付方式的收银台的易支付~~

## 部署

1. 下载 Releases 中对应系统和架构类型的二进制可执行文件
2. 复制 .env.example 到 .env
3. 根据注释修改配置文件
4. 启动程序，以部署Cloudreve的相同方式部署本程序

```env
# 是否启用debug模式
CR_EPAY_DEBUG=true
# 监听端口，TLS 请使用其他服务器进行反代
CR_EPAY_LISTEN=:4560
# 后台 - 增值服务 - 通信密钥 建议随机生成uuid 请务必保密 https://www.uuidgenerator.net/
CR_EPAY_CLOUDREVE_KEY=
# 本站点的外部访问 URL
CR_EPAY_BASE=https://payment.cloudreve.dev
# 自定义订单名称
# CR_EPAY_CUSTOM_NAME=TESTTTTT
# 商家ID
CR_EPAY_EPAY_PARTNER_ID=1010
# 商家密钥
CR_EPAY_EPAY_KEY=SFDHSKHFJKDSHEUIFHU
# 更换成你的易支付网关
CR_EPAY_EPAY_ENDPOINT=https://payment.moe/submit.php
# 支付方式 wxpay 或 alipay
CR_EPAY_EPAY_PURCHASE_TYPE=alipay
# 是否启用redis 请务必启用
CR_EPAY_REDIS_ENABLED=true
CR_EPAY_REDIS_SERVER=localhost:6379
# CR_EPAY_REDIS_PASSWORD=
CR_EPAY_REDIS_DB=0
```

## 设置

1. 打开 Cloudreve 后台，打开 `参数设置` `增值服务`
2. 开启 `自定义付款渠道`
3. 填入 `付款方式名称`
4. `通讯密钥填入` 上一步 `CR_EPAY_CLOUDREVE_KEY` 的值
5. `支付接口地址` 填入 上一步 `CR_EPAY_BASE` 的值 加上 `/cloudreve/purchase`
6. 保存设置

## CHANGELOG

### 0.2

- 修复易支付自定义付款方式，`CR_EPAY_EPAY_PURCHASE_TYPE`
- 支持自定义商品名称
