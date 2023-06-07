/**
 * Created by PhpStorm.
 * @file   party.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/5 10:45
 * @desc   party.go
 */

package services

import (
	"context"
	"go-client/params"
	"log"
)

type PartyService struct {
}

func NewPartyService() *PartyService {
	return &PartyService{}
}
func (*PartyService) CreateParty(client params.Connect) (string, error) {
	party, err := client.Conn.PartyCreate(context.Background(), true, 10)
	if err != nil {
		log.Panic(err)
	}
	return party.PartyId, nil
}

func (*PartyService) JoinParty(client params.Connect, partyId string) error {
	err := client.Conn.PartyJoin(context.Background(), partyId)
	if err != nil {
		log.Panic(err)
	}
	return err
}
