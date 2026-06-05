package channels

import "github.com/hicongcn/xuanwu-panel/internal/sdk/message"

type QyWeiXinChannel struct{ *BaseChannel }

func NewQyWeiXinChannel() Channel {
	return &QyWeiXinChannel{NewBaseChannel(ChannelQyWeiXin, []string{FormatTypeMarkdown, FormatTypeText})}
}

func (c *QyWeiXinChannel) Send(config ChannelConfig, msg *Message) (*Result, error) {
	accessToken := config.GetString("access_token")

	if accessToken == "" {
		return SendError("qyweixin config missing: access_token is required"), nil
	}

	contentType, formattedContent := c.FormatContent(msg)
	atList := []string{}
	atList = append(atList, msg.GetAtUserIds()...)
	atList = append(atList, msg.GetAtMobiles()...)
	if msg.AtAll {
		atList = append(atList, "@all")
	}

	cli := message.QyWeiXin{AccessToken: accessToken}
	var res []byte
	var err error

	switch contentType {
	case FormatTypeText:
		res, err = cli.SendMessageText(formattedContent, atList...)
	case FormatTypeMarkdown:
		res, err = cli.SendMessageMarkdown(msg.Title, formattedContent, atList...)
	default:
		return SendError("未知的企业微信发送内容类型：%s", contentType), nil
	}

	if err != nil {
		return ErrorResult(string(res), err), nil
	}
	return SuccessResult(string(res)), nil
}
