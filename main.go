package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

func StreamHandler(stream network.Stream) {
	fmt.Println("connected")

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			panic(err)
		}

		if str == "" {
			return
		}
		if str != "\n" {
			fmt.Printf("\x1b[32m%s\x1b[0m>", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}

		_, err = rw.WriteString(fmt.Sprintf("%s\n", sendData))
		if err != nil {
			fmt.Println("Error writing to buffer")
			panic(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Error flushing buffer")
			panic(err)
		}
	}
}

func main() {

	cfg := parseFlags()

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.listenHost, cfg.listenPort))

	r := rand.Reader

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}

	host, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}

	host.SetStreamHandler(protocol.ID(cfg.ProtocolID), StreamHandler)

	peerChan := initMDNS(host, cfg.rendevouz)

	for {
		peer := <-peerChan
		if peer.ID > host.ID() {
			fmt.Println("Found peer:", peer, " id is greater than us, wait for it to connect to us")
			continue
		}

		fmt.Println("New Peer - ", peer.ID)

		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

		if err := host.Connect(ctx, peer); err != nil {
			fmt.Println("Error connecting to peer")
			continue
		}

		stream, err := host.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))
		if err != nil {
			fmt.Println("Error creating stream")
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

			go writeData(rw)
			go readData(rw)
			fmt.Println("Connected to:", peer)
		}
	}

}
