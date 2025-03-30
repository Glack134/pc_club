package remote

import (
	"github.com/pion/webrtc/v3"
)

type Peer struct {
	conn *webrtc.PeerConnection
}

func New() (*Peer, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}

	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	return &Peer{conn: pc}, nil
}

func (p *Peer) Connect(signalServer string) error {
	// Реализация подключения к signaling серверу
	return nil
}
