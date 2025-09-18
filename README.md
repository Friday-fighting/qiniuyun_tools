# qiniuyun_tools
七牛云go工具库，上传文件、解析链接、获取上传凭证


# Usage

## get
```bash
go get github.com/zhengyuxiang/qiniuyun_tools
```

## use
```go
import (
    "github.com/zhengyuxiang/qiniuyun_tools"
)
ak := "your qiniu access key"
sk := "your qiniu secret key"
urlPrefix := "your url prefix"  # http://xxx.qiniudn.com/
bucket := "your bucket name"
expTimeKey := "your qiniu exp time key"
expTime := 3600

client := qiniuyun.NewClient(ak, sk, urlPrefix, bucket, expTimeKey, expTime)

### 操作
......

```