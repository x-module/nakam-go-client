/**
 * Created by goland.
 * @file   gamelift.go
 * @author 李锦 <Lijin@cavemanstudio.net>
 * @date   2022/12/13 19:18
 * @desc   gamelift.go
 */

package handler

import (
	"github.com/x-module/utils/utils/xlog"
	"go-client/pkg/gamelift"
	"go-client/pkg/gamelift/proto/pbuffer"
	"net/http"
)

type Handler struct {
	GameSession pbuffer.GameSession
	C           gamelift.Client
	port        int
}

func (h *Handler) GetSession() pbuffer.GameSession {
	return h.GameSession
}
func (h *Handler) StartGameSession(event *pbuffer.ActivateGameSession) {
	xlog.Logger.Debug("StartGameSession", event)
	h.GameSession = *event.GameSession
	if err := h.C.ActivateGameSession(&pbuffer.GameSessionActivate{
		GameSessionId: event.GetGameSession().GetGameSessionId(),
		MaxPlayers:    event.GetGameSession().GetMaxPlayers(),
		Port:          event.GetGameSession().GetPort(),
	}); err != nil {
		xlog.Logger.Panic(err)
	}
	xlog.Logger.Debug("ActivateGameSession complete. sessionID:", *h.C.GetGameSessionId())
}

func (h *Handler) UpdateGameSession(event *pbuffer.UpdateGameSession) {
	xlog.Logger.Debug("UpdateGameSession", event)
}

func (h *Handler) ProcessTerminate(event *pbuffer.TerminateProcess) {
	xlog.Logger.Debug("ProcessTerminate", event)
}

func (h *Handler) HealthCheck() bool {
	xlog.Logger.Debug("HealthCheck")
	return true
}

func (h *Handler) AcceptPlayerHandler(w http.ResponseWriter, r *http.Request) {
	psess := r.URL.Query().Get("psess")
	if err := h.C.AcceptPlayerSession(&pbuffer.AcceptPlayerSession{
		GameSessionId:   *h.C.GetGameSessionId(),
		PlayerSessionId: psess,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemovePlayerHandler(w http.ResponseWriter, r *http.Request) {
	psess := r.URL.Query().Get("psess")
	if err := h.C.RemovePlayerSession(&pbuffer.RemovePlayerSession{
		GameSessionId:   *h.C.GetGameSessionId(),
		PlayerSessionId: psess,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) TerminateHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.C.TerminateGameSession(&pbuffer.GameSessionTerminate{
		GameSessionId: *h.C.GetGameSessionId(),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// This is hidoi.
	xlog.Logger.Panic("terminate called")

	w.WriteHeader(http.StatusNoContent)
}
