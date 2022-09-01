package minecraft

import (
	"fmt"
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/net/packet"
	"github.com/XiaoXianHW/ATProxy/config"
	"time"
)

func generateKickMessage(s *config.ConfigProxyService, name packet.String) chat.Message {
	return chat.Message{
		Color: chat.White,
		Extra: []chat.Message{
			{Bold: true, Color: chat.yellow, Text: "ATProxy \n"},
			{Text: "\n"},
			{Bold: false, Color: chat.Red, Text: "您已被代理服务器踢出！\n"},
			{Color: chat.Gray, Text: "原因:您的ID不在白名单范围内！\n"},
			{Text: "\n"},
			{
				Color: chat.Gray,
				Text: fmt.Sprintf("玩家ID: %s | 代理服务器: %s\n", name, s.Name),
			},
		},
	}
}

func generatePlayerNumberLimitExceededMessage(s *config.ConfigProxyService, name packet.String) chat.Message {
	return chat.Message{
		Color: chat.White,
		Extra: []chat.Message{
			{Bold: true, Color: chat.yellow, Text: "ATProxy \n"},
			{Text: "\n"},
			{Bold: false, Color: chat.Red, Text: "代理服务器拒绝了您的连接请求！\n"},
			{Color: chat.Gray, Text: "原因:当前代理服务器玩家已达到上限！\n"},
			{
				Color: chat.Gray,
				Text: fmt.Sprintf("玩家ID: %s | 代理服务器: %s\n", name, s.Name),
			},
		},
	}
}
