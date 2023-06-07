/**
 * Created by PhpStorm.
 * @file   handler.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/5 10:41
 * @desc   handler.go
 */

package services

import (
	"context"
	"github.com/ascii8/nakama-go"
	"github.com/x-module/utils/global"
	"github.com/x-module/utils/utils/xlog"
	"go-client/params"
	"sync"
	"time"
)

type HandlerService struct {
	ConnectHandler              func(context.Context)
	DisconnectHandler           func(context.Context, error)
	ErrorHandler                func(context.Context, *nakama.ErrorMsg)
	ChannelMessageHandler       func(context.Context, *nakama.ChannelMessageMsg)
	ChannelPresenceEventHandler func(context.Context, *nakama.ChannelPresenceEventMsg)
	MatchDataHandler            func(context.Context, *nakama.MatchDataMsg)
	MatchPresenceEventHandler   func(context.Context, *nakama.MatchPresenceEventMsg)
	// MatchmakerMatchedHandler    func(context.Context, *nakama.MatchmakerMatchedMsg)
	NotificationsHandler       func(context.Context, *nakama.NotificationsMsg)
	StatusPresenceEventHandler func(context.Context, *nakama.StatusPresenceEventMsg)
	StreamDataHandler          func(context.Context, *nakama.StreamDataMsg)
	StreamPresenceEventHandler func(context.Context, *nakama.StreamPresenceEventMsg)
}

func NewHandlerService() *HandlerService {
	return new(HandlerService)
}

var hasJoined sync.Map

// func (h *HandlerService) MatchmakerMatchedHandler(ctx context.Context, msg *nakama.MatchmakerMatchedMsg) {
// 	b, _ := json.Marshal(*msg)
// 	var matchMakeData MatchMakeData
// 	_ = json.Unmarshal(b, &matchMakeData)
// 	// fmt.Printf("recive MatchmakerMatchedHandler msg:%+v\n", string(b))
// 	conn := OnlinePlayerList[matchMakeData.Self.Presence.UserID]
// 	m2, _ := conn.Conn.MatchJoin(ctx, matchMakeData.ID.MatchID, nil)
// 	fmt.Println("join match id:", m2.MatchId)
// }

func (h *HandlerService) MatchmakerMatchedHandler(ctx context.Context, msg *nakama.MatchmakerMatchedMsg) {
	if _, exist := params.StartMatchDsService.Load(msg.GetMatchId()); !exist {
		params.StartMatchDsService.Store(msg.GetMatchId(), true)
		go func() {
			xlog.Logger.Info("开始启动GameLift服务...")
			dsHandler := NewGameLiftService().StartGameliftServer()
			xlog.Logger.Info("开始启动健康上报服务...")
			GameServerService := NewReportService()
			for {
				err := GameServerService.HealthReport(dsHandler)
				if err != nil {
					xlog.Logger.WithField(global.ErrField, err).Error("report health error")
				}
				time.Sleep(time.Second * 5)
			}
		}()
	}
	userList := msg.GetUsers()
	for _, user := range userList {
		xlog.Debugf("player start join match ...\n")
		val, _ := params.OnlinePlayerList.Load(user.Presence.UserId)
		conn := val.(params.Connect)
		xlog.Debugf("player:%s start join match \n", conn.UserId)
		if _, exist := hasJoined.Load(user.Presence.UserId); !exist {
			hasJoined.Store(user.Presence.UserId, true)
			_, err := conn.Conn.MatchJoin(ctx, msg.GetMatchId(), nil)
			if err != nil {
				xlog.WithField(global.ErrField, err).Error("join match fail.match id:", msg.GetMatchId())
			} else {
				xlog.Debugf("player:%s join match success\n", conn.UserId)
				// 开始启动健康上报
				go NewMatchService().SendStatus(conn, msg.GetMatchId())
			}
		} else {
			xlog.Debugf("player:%s already join match ,jump \n", conn.UserId)
		}
	}

}
