package main

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

const difficulty = 6

type ProofOfWork struct {
	*Block
	Target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-difficulty))

	pow := &ProofOfWork{
		Block:  block,
		Target: target,
	}
	return pow
}

func (pow *ProofOfWork) PrepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			StringToHex(pow.Block.PreviousHash),
			StringToHex(pow.Block.Data),
			IntToHex(int(pow.Block.TimeStamp)),
			IntToHex(difficulty),
			IntToHex(nonce),
		}, []byte{},
	)
	return data
}

func StringToHex(s string) []byte {
	return []byte(s)
}

func IntToHex(i int) []byte {
	tmp := int32(i)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &tmp)
	return bytesBuffer.Bytes()
}
