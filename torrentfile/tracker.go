package torrentfile

import (
	"github.com/jackpal/bencode-go"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"torrent2magnet/peers"
)

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil { return "", err }

	params := url.Values{
		// info hash identify the torrent file we're trying to download
		"info_hash": []string{string(t.InfoHash[:])},
		// peer_id identify ourself to tracker server, random bytes
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (t *TorrentFile) requestPeers(peerID [20]byte, port uint16) ([]peers.Peer, error) {
	// build tracker url from torrent file
	trackerUrl, err := t.buildTrackerURL(peerID, port)
	if err != nil { return nil, err }

	c := http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(trackerUrl)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	// parse response to get peers data
	trackerResp := bencodeTrackerResp{}
	if err := bencode.Unmarshal(resp.Body, &trackerResp); err != nil {
		return nil, err
	}
	return peers.Unmarshal([]byte(trackerResp.Peers))
}
