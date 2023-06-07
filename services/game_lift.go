/**
 * Created by goland.
 * @file   game_lift.go
 * @author 李锦 <Lijin@cavemanstudio.net>
 * @date   2023/2/17 15:02
 * @desc   game_lift.go
 */

package services

import (
	"fmt"
	utils2 "github.com/x-module/utils/utils"
	"github.com/x-module/utils/utils/xlog"
	"go-client/pkg/gamelift"
	"go-client/pkg/gamelift/handler"
	"go-client/pkg/gamelift/log"
	"go-client/pkg/gamelift/proto/pbuffer"
)

type GameLiftService struct {
}

func NewGameLiftService() *GameLiftService {
	return new(GameLiftService)
}

func (*GameLiftService) StartGameliftServer() *handler.Handler {
	conn, port, err := utils2.OpenFreeUDPPort(9000, 100)
	if err != nil {
		xlog.Logger.Panic(err)
	}
	defer conn.Close()
	logger := &log.StandardLogger{}
	c := gamelift.NewClient(logger)
	h := &handler.Handler{C: c}
	c.Handle(h)
	if err := c.Open(); err != nil {
		xlog.Logger.Panic(err)
	}
	if err := c.ProcessReady(&pbuffer.ProcessReady{
		LogPathsToUpload: []string{},
		Port:             int32(port),
		// MaxConcurrentGameSessions: 0, // not set in original ServerSDK
	}); err != nil {
		xlog.Logger.Panic(err)
	}
	res, err := c.GetInstanceCertificate(&pbuffer.GetInstanceCertificate{})
	fmt.Println("=============================================================")
	xlog.Logger.Debug(res, err)
	return h
}
