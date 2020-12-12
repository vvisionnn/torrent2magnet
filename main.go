package main

import (
	"fmt"
	"github.com/jackpal/bencode-go"
	"io"
	"os"
	"reflect"
)

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

// Open parses a torrent file
func Open(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
}

func main() {
	filePath := ""
	f, _ := os.OpenFile(
		filePath,
		os.O_RDONLY,
		0755)
	btoP, err := Open(f)
	if err != nil {
		fmt.Println(err)
	}

	t := reflect.TypeOf(*btoP)
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Tag.Get("bencode"))
	}
}
