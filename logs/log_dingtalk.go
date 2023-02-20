package logs

import (
	"fmt"
	"github.com/blinkbean/dingtalk"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type LogContentDingTalk struct {
	Instance     string
	Token        string
	ProxyUrl     string
	Secret       string
	BusinessName string
	StartTime    time.Time
	EndTime      time.Time
	CostTime     string
	FileName     string
	FileSize     string
	Status       bool
}

func (dt LogContentDingTalk) ResultToDingTalkGroup() error {
	var MsgContent []string

	fmt.Println("------------")
	fmt.Println(dt.ProxyUrl)

	if dt.ProxyUrl != "" {
		err := os.Setenv("https_proxy", dt.ProxyUrl)
		if err != nil {
			logrus.Errorf("Can not set Proxy[%s] url error: %v", dt.ProxyUrl, err)
			return err
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

	MdContent3 := fmt.Sprintf("- 数据库实例：%s", dt.Instance)
	MdContent4 := fmt.Sprintf("- 备份开始时间：%s", dt.StartTime.Format("2006-01-02 15:04:05"))
	MdContent5 := fmt.Sprintf("- 备份结束时间：%s", dt.EndTime.Format("2006-01-02 15:04:05"))
	MdContent6 := fmt.Sprintf("- 备份耗时：%s", dt.CostTime)
	MdContent7 := fmt.Sprintf("- 备份文件：%s", dt.FileName)
	MdContent8 := fmt.Sprintf("- 备份文件大小：%s", dt.FileSize)

	MsgContent = append(MsgContent, MdContent3, MdContent4, MdContent5, MdContent6, MdContent7, MdContent8)

	err := cli.SendMarkDownMessageBySlice(MsgTitle, MsgContent)
	if err != nil {
		logrus.Errorf("Send DingTalk message error: %v", err)
		return err
	}

	return nil
}
