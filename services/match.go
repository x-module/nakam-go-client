/**
 * Created by goland.
 * @file   match.go
 * @author 李锦 <Lijin@cavemanstudio.net>
 * @date   2023/2/22 19:50
 * @desc   match.go
 */

package services

import (
	"encoding/json"
	"github.com/x-module/utils/global"
	"github.com/x-module/utils/utils/xlog"
	"go-client/params"
	"time"
)

type MatchService struct {
}

func NewMatchService() *MatchService {
	return new(MatchService)
}

func (*MatchService) sendData(client params.Connect, matchId string, opCode int64, data []byte) error {
	err := client.Conn.MatchDataSend(client.Ctx, matchId, opCode, data, true)
	if err != nil {
		xlog.WithField(global.ErrField, err).Error("send match data error")
	} else {
		// xlog.Debug(client.Email + " send match data success")
	}
	return err
}

func (m *MatchService) SendStatus(client params.Connect, matchId string) {
	data, _ := json.Marshal(map[string]any{"data": "test", "id": 123})
	for _ = range time.Tick(10 * time.Second) {
		_ = m.sendData(client, matchId, 501, data)
	}
}
