package main

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type LogContentDingTalk struct {
	Token        string
	ProxyUrl     string
	Secret       string
	BusinessName string
	StartTime    time.Time
	EndTime      time.Time
	FileName     string
	FileSize     string
	Status       bool
}

func (dt LogContentDingTalk) ResultToDingTalkGroup() {
	var MsgContent []string

	if dt.ProxyUrl != "" {
		err := os.Setenv("http_proxy", dt.ProxyUrl)
		if err != nil {
			logrus.Errorf("Can not set Proxy[%s] url error: %v", dt.ProxyUrl, err)
			return
		}
	}

	cli := dingtalk.InitDingTalkWithSecret(dt.Token, dt.Secret)

	MsgTitle := fmt.Sprintf("- MySQL备份-%s", dt.BusinessName)
	MdContent1 := fmt.Sprintf("#### 数据库备份【MySQL-%s】", dt.BusinessName)
	MdContent2 := fmt.Sprintf("---")
	MsgContent = append(MsgContent, MdContent1, MdContent2)

	if dt.Status {
		MsgContent = append(MsgContent, "- <font color=#148F77 size=6>备份任务：Success</font>")
	} else {
		MsgContent = append(MsgContent, "- <font color=#B25642 size=6>备份任务：Failed</font>")
	}

	MdContent3 := fmt.Sprintf("- 备份开始时间：%s", dt.StartTime.Format("2006-01-02 15:04:05"))
	MdContent4 := fmt.Sprintf("- 备份结束时间：%s", dt.EndTime.Format("2006-01-02 15:04:05"))
	MdContent5 := fmt.Sprintf("- 备份文件：%s", dt.FileName)
	MdContent6 := fmt.Sprintf("- 备份文件大小：%s", dt.FileSize)

	MsgContent = append(MsgContent, MdContent3, MdContent4, MdContent5, MdContent6)

	err := cli.SendMarkDownMessageBySlice(MsgTitle, MsgContent)
	if err != nil {
		logrus.Fatal(err)
	}

}

func main() {
	const DingTalkToken = "a6fa9bd1b631f2caf7ccbb76eff7db05d193c58f777cdda79913508c6a29a1f3"
	const Secret = "SEC0a7cfe8ca36ab44d17a54ec6e8904235c90b33c52969ee8329a0735135c94198"
	dt := LogContentDingTalk{
		Token:        DingTalkToken,
		ProxyUrl:     "http://192.168.55.208:1080",
		BusinessName: "德安备份通知测试",
		StartTime:    time.Now(),
		EndTime:      time.Now(),
		FileName:     "abc.sql",
		FileSize:     "12G",
		Status:       false,
		Secret:       Secret,
	}

	dt.ResultToDingTalkGroup()
}
