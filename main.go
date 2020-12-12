package main

import (
	"encoding/binary"
	"fmt"
	"github.com/jackpal/bencode-go"
	"io"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
)

type TorrentFile struct {
	// announce server
	Announce	string
	// torrent file hash code, SHA-1
	InfoHash	[20]byte
	PieceHashes	[][20]byte
	PieceLength	int
	Length		int
	Name		string
}

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce     string      `bencode:"announce"`
	Comment      string      `bencode:"comment"`
	CreationDate int         `bencode:"creation date"`
	Info         bencodeInfo `bencode:"info"`
}

type Peer struct {
	IP		net.IP
	Port	uint16
}


// Unmarshal parse peer IP from buffer
func Unmarshal(peersBin []byte) ([]Peer, error) {
	const peerSize = 6	// 4 for IP, 2 for port
	numPeers := len(peersBin) / peerSize
	if len(peersBin) % peerSize != 0 {
		err := fmt.Errorf("received malformed peers")
		return nil, err
	}

	peers := make([]Peer, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = peersBin[offset: offset+4]
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4: offset+6])
	}
	return peers, nil
}

// Open parses a torrent file
func Open(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
}


func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port int) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil { return "", err }

	params := url.Values{
		// info hash identify the torrent file we're trying to download
		"info_hash":	[]string{ string(t.InfoHash[:]) },
		// peer_id identify ourself to tracker server, random bytes
		"peer_id":		[]string{ string(peerID[:]) },
		"port":			[]string{ strconv.Itoa(port) },
		"uploaded":		[]string{ "0" },
		"downloaded":	[]string{ "0" },
		"compact":		[]string{ "1" },
		"left": 		[]string{ strconv.Itoa(t.Length) },
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}


func main() {
	filePath := "test.torrent"
	f, _ := os.OpenFile(filePath, os.O_RDONLY, 0755)
	btoP, err := Open(f)
	if err != nil { fmt.Println(err) }

	t := reflect.TypeOf(*btoP)
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Tag.Get("bencode"))
	}
}
