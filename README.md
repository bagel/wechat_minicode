# wechat_minicode
微信小程序码生成

#### 使用方法

```go
import (
    "github.com/bagel/wechat_minicode"
)

func main() {
    appId := "xxxxx"
    appSecret := "xxxxx"
    scene := "参数"
    width := 430

    wx := utils.NewWechatCode(appId, appSecret, scene, width)

    body, err := wx.GetWechatCode()
}
```
