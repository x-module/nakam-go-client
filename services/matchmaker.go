/**
 * Created by PhpStorm.
 * @file   matchmaker.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/5 10:47
 * @desc   matchmaker.go
 */

package services

import (
	"context"
	"github.com/ascii8/nakama-go"
	"go-client/params"
	"log"
)

type MatchmakerService struct {
	client *ClientService
}

func NewMatchmakerService() *MatchmakerService {
	return &MatchmakerService{}
}

// func (m *MatchmakerService) PartyMatchmaker(conn *nakama.Conn, params params.AddMatchMakerParams) error {
// 	_, err := conn.MatchmakerAdd(context.Background(), &nakama.MatchmakerAddMsg{
// 		MinCount:          int32(params.MinCount),
// 		MaxCount:          int32(params.MaxCount),
// 		Query:             params.Query,
// 		StringProperties:  params.StringProperties,
// 		NumericProperties: params.NumericProperties,
// 	})
// 	if err != nil {
// 		log.Println("MatchmakerAdd error,err:", err.Error())
// 		return err
// 	}
// 	return err
// }

func (*MatchmakerService) PartyMatchmaker(conn *nakama.Conn, partyId string, params params.AddMatchMakerParams) error {
	_, err := nakama.PartyMatchmakerAdd(partyId, params.Query, params.MinCount, params.MaxCount).WithStringProperties(params.StringProperties).
		WithNumericProperties(params.NumericProperties).Send(context.Background(), conn)
	if err != nil {
		log.Println("start matchmaker err:", err.Error())
		return err
	}
	// log.Println("matchmaker result:", result.String())
	return err
}
