package block

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Block struct {
	Hash      string
	Data      string
	Timestamp string
	Previous  string
	UID       string
}

func (block *Block) BuildHash() {
	hash := sha256.New()

	var bin_buf bytes.Buffer
	binary.Write(&bin_buf, binary.BigEndian, block.Data)
	binary.Write(&bin_buf, binary.BigEndian, block.Timestamp)
	binary.Write(&bin_buf, binary.BigEndian, block.UID)

	hash.Write(bin_buf.Bytes())
	block.Hash = fmt.Sprintf("%x", hash.Sum(nil))
}

func (block *Block) SaveBlock(path string) error {
	json, _ := json.Marshal(block)
	filename := path + block.Timestamp + "-" + block.Hash + ".zip"
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	zipWriter := zip.NewWriter(out)
	defer zipWriter.Close()
	blockf, err := zipWriter.Create("block.json")
	if err != nil {
		return err
	}
	_, err = blockf.Write(json)
	return err
}

func LoadBlock(filename string) (Block, error) {
	var block Block
	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		return block, err
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		if f.Name == "block.json" {
			blockFile, err := f.Open()
			if err != nil {
				return block, err
			}
			buf := bytes.NewBuffer(nil)
			io.Copy(buf, blockFile)
			err = json.Unmarshal(buf.Bytes(), &block)
			if err != nil {
				return block, err
			}
			return block, err
		}
	}
	return block, errors.New("Error while loading block")
}
