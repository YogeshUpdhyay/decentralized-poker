package p2p

import (
	"context"
	"encoding/json"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"github.com/YogeshUpdhyay/ypoker/internal/db"
	"github.com/YogeshUpdhyay/ypoker/internal/eventbus"
	eventModels "github.com/YogeshUpdhyay/ypoker/internal/eventbus/models"
	p2pModels "github.com/YogeshUpdhyay/ypoker/internal/p2p/models"
	"github.com/YogeshUpdhyay/ypoker/internal/utils"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	HandleMessage(ctx context.Context, msg *p2pModels.Envelope) error
}

type DefaultHandler struct{}

func (h *DefaultHandler) HandleMessage(ctx context.Context, msg *p2pModels.Envelope) error {
	switch msg.Type {
	case p2pModels.MsgTypeHandshake:
		var hs p2pModels.HandShake
		json.Unmarshal(msg.Payload, &hs)
		handleHandshake(ctx, msg.From, hs)

	case p2pModels.MsgTypeHandshakeAck:
		var hsAck p2pModels.HandShake
		json.Unmarshal(msg.Payload, &hsAck)
		handleHandshakeAck(ctx, msg.From, hsAck)

	case p2pModels.MsgTypeHandshakeReject:
		handleRejection(ctx, msg.From)

	case p2pModels.MsgTypeChat:
		var chat p2pModels.ChatPayload
		json.Unmarshal(msg.Payload, &chat)
		handleChatMessage(msg.From, chat)
	}
	return nil
}

func handleChatMessage(peerID string, chat p2pModels.ChatPayload) {
	panic("unimplemented")
}

func handleRejection(ctx context.Context, peerID string) {
	tx := db.Get().
		Model(&db.ConnectionRequests{}).
		Where("peer_id = ? AND status = ?", peerID, constants.RequestStatusSent).
		Update("status", constants.RequestStatusRejected)
	if tx.Error != nil {
		log.
			WithContext(ctx).
			WithField(constants.PeerID, peerID).
			Infof("error updating request status to rejected: %s", tx.Error.Error())
		return
	}
}

func handleHandshakeAck(ctx context.Context, peerID string, hsAck p2pModels.HandShake) {
	if hsAck.Version != utils.GetAppConfig().Version {
		log.WithContext(ctx).Info("version mismatch")
		return
	}

	var req db.ConnectionRequests
	tx := db.Get().
		Where(&db.ConnectionRequests{
			PeerID: peerID,
			Status: constants.RequestStatusSent,
		}).
		First(&req)
	if tx.Error != nil {
		log.
			WithContext(ctx).
			WithField(constants.PeerID, peerID).
			Infof("error getting connection request: %s", tx.Error.Error())
		return
	}

	req.Status = constants.RequestStatusAccepted
	req.Username = hsAck.UserInfo.Username
	req.AvatarUrl = hsAck.UserInfo.AvatarUrl

	if err := db.Get().Save(&req).Error; err != nil {
		log.
			WithContext(ctx).
			WithField(constants.PeerID, peerID).
			Infof("error updating connection request: %s", err.Error())
		return
	}

	peerInfo := db.PeerInfo{
		PeerID:    peerID,
		Username:  hsAck.UserInfo.Username,
		AvatarUrl: hsAck.UserInfo.AvatarUrl,
		Status:    constants.ConnectionStateInactive,
		Address:   req.Address,
	}
	tx = db.Get().Create(&peerInfo)
	if tx.Error != nil {
		log.
			WithContext(ctx).
			WithField(constants.PeerID, peerID).
			Infof("error creating peer info %s", tx.Error.Error())
		return
	}
}

func handleHandshake(ctx context.Context, peerID string, hs p2pModels.HandShake) {
	peer := GetServer().GetPeerFromPeerID(ctx, peerID)
	if peer == nil {
		log.WithContext(ctx).Infof("no peer connection found for id: %s", peerID)
		return
	}

	if utils.GetAppConfig().Version != hs.Version {
		log.WithContext(ctx).Infof("version mismatch rejecting request %s", peerID)
		peer.Send(&p2pModels.Envelope{
			Type: p2pModels.MsgTypeHandshakeReject,
		})
	}

	// storing the connection request awaiting approval
	connectionReq := db.ConnectionRequests{
		PeerID:    peerID,
		Status:    constants.RequestStatusAwaitingDecision,
		Username:  hs.UserInfo.Username,
		AvatarUrl: hs.UserInfo.AvatarUrl,
		Address:   peer.GetMultiAddrs(ctx).String(),
	}
	tx := db.Get().Create(&connectionReq)
	if tx.Error != nil {
		log.WithContext(ctx).Infof("error adding connection request to the db: %s", tx.Error.Error())
		return
	}

	// emitting event
	eventbus.Get().Publish(constants.EventNewConnectionRequest, eventModels.NewConnectionRequestEventData{
		PeerID:    peerID,
		Username:  hs.UserInfo.Username,
		AvatarUrl: hs.UserInfo.AvatarUrl,
	})
}
