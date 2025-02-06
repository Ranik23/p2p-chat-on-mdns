package main


import (
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

type Notify struct {
	PeerChan chan peer.AddrInfo
}

func (n *Notify) HandlePeerFound(pi peer.AddrInfo) {
	n.PeerChan <- pi
}

func initMDNS(peerHost host.Host, rendezvous string) chan peer.AddrInfo {
	notify := &Notify{}
	ch := make(chan peer.AddrInfo)

	notify.PeerChan = ch

	service := mdns.NewMdnsService(peerHost, rendezvous, notify)

	if err := service.Start(); err != nil {
		panic(err)
	}

	return notify.PeerChan

}

