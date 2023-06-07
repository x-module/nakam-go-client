/**
 * Created by goland.
 * @file   health_report.go
 * @author 李锦 <Lijin@cavemanstudio.net>
 * @date   2023/2/17 16:36
 * @desc   health_report.go
 */

package services

import (
	"errors"
	"fmt"
	"github.com/x-module/utils/utils/request"
	"github.com/x-module/utils/utils/xlog"
	"go-client/params"
	"go-client/pkg/gamelift/handler"
	"time"
)

type ReportService struct {
}

func NewReportService() *ReportService {
	return new(ReportService)
}

type CreateLobbyConfig struct {
	LobbyType            int    `json:"LobbyType"` // Match Or Coop
	ModLobby             bool   `json:"ModLobby"`
	MaxPlayerNum         int    `json:"MaxPlayerNum"`
	MapTag               string `json:"MapTag"`
	ModeTag              string `json:"ModeTag"`
	LoadoutTag           string `json:"LoadoutTag"`
	Region               string `json:"Region"`
	Password             int    `json:"Password"`
	GameType             string `json:"GameType"`
	ServerMessage        string `json:"ServerMessage"`
	ModType              int    `json:"ModType"`
	IsLocalServer        bool   `json:"IsLocalServer"`
	OculusLobbySessionId string `json:"OculusLobbySessionId"`
	MemberCount          int    `json:"MemberCount"`
	TimeLength           int64  `json:"TimeLength"`
}

func (*ReportService) HealthReport(handler *handler.Handler) error {
	time.Sleep(time.Second * 10)
	session := handler.GetSession()
	// utils.JsonDisplay(session.GameProperties)
	if len(session.GameProperties) < 2 {
		return nil
		// xlog.Logger.Info("game properties is empty")
		// return errors.New("game properties is empty")
	}

	if len(session.GameProperties) < 2 || session.GameProperties[1].Value == "" {
		xlog.Logger.Warn("game session is empty")
		return errors.New("game session is empty")
	}
	url := fmt.Sprintf("%s:%d/v2/rpc/report/healthy?http_key=defaulthttpkey&unwrap", params.Nakama.Host, params.Nakama.Port)
	res, err := request.NewRequest().Debug(true).Post(url, map[string]string{
		"MatchId": session.GameProperties[0].Value,
		"Secret":  session.GameProperties[1].Value,
	})
	if err != nil {
		fmt.Println("reportHealth error:", err.Error())
	} else {
		result, _ := res.Content()
		fmt.Println("response:", result)
	}
	// if params.Config[0].TimeLength+BeginTime < time.Now().Unix() {
	// 	ReportTermination(params, session)
	// 	break
	// }
	return nil
}
