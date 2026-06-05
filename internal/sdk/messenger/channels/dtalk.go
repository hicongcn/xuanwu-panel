package channels

import "github.com/hicongcn/xuanwu-panel/internal/sdk/message"

type DtalkChannel struct{ *BaseChannel }

func NewDtalkChannel() Channel {
	return &DtalkChannel{NewBaseChannel(ChannelDtalk, []string{FormatTypeMarkdown, FormatTypeText})}
}

func (c *DtalkChannel) Send(config ChannelConfig, msg *Message) (*Result, error) {
	accessToken := config.GetString("access_token")
	secret := config.GetString("secret")

	if accessToken == "" {
		return SendError("dtalk config missing: access_token is required"), nil
	}

	contentType, formattedContent := c.FormatContent(msg)
	atMobiles := msg.GetAtMobiles()
	if msg.AtAll {
		atMobiles = append(atMobiles, "all")
	}

	cli := message.Dtalk{AccessToken: accessToken, Secret: secret}
	var res []byte
	var err error

	switch contentType {
	case FormatTypeText:
		res, err = cli.SendMessageText(formattedContent, atMobiles...)
	case FormatTypeMarkdown:
		res, err = cli.SendMessageMarkdown(msg.Title, formattedContent, atMobiles...)
	default:
		return SendError("未知的钉钉发送内容类型：%s", contentType), nil
	}

	if err != nil {
		return ErrorResult(string(res), err), nil
	}
	return SuccessResult(string(res)), nil
}
