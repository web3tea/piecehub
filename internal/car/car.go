package car

import (
	"context"
	"fmt"
	"io"
	"os"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car/v2/blockstore"
	"github.com/multiformats/go-multihash"
)

func GenerateCar(size, chunksize int64) (string, cid.Cid, error) {
	reader := NewRepeatedReader(size, chunksize)
	ctx := context.TODO()

	firstBlock := make([]byte, chunksize)
	n, err := reader.Read(firstBlock)
	if err != nil && err != io.EOF {
		return "", cid.Cid{}, fmt.Errorf("failed to read first block: %w", err)
	}
	if n == 0 {
		return "", cid.Cid{}, fmt.Errorf("no data in input stream")
	}

	rootCid := createCID(firstBlock[:n])
	if err != nil {
		return "", cid.Cid{}, fmt.Errorf("failed to create root CID: %w", err)
	}

	carFile, err := os.CreateTemp(os.TempDir(), "car-")
	if err != nil {
		return "", cid.Cid{}, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer carFile.Close()

	bs, err := blockstore.OpenReadWriteFile(carFile, []cid.Cid{rootCid}, blockstore.WriteAsCarV1(true))
	if err != nil {
		return "", cid.Cid{}, fmt.Errorf("failed to create blockstore: %w", err)
	}
	defer bs.Finalize()

	blk, _ := blocks.NewBlockWithCid(firstBlock[:n], rootCid)
	err = bs.Put(ctx, blk)
	if err != nil {
		return "", cid.Cid{}, fmt.Errorf("failed to write first block: %w", err)
	}

	buf := make([]byte, chunksize)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return "", cid.Cid{}, fmt.Errorf("failed to read data: %w", err)
		}
		if n == 0 {
			break
		}

		blk, _ := blocks.NewBlockWithCid(buf[:n], createCID(buf[:n]))
		err = bs.Put(ctx, blk)
		if err != nil {
			return "", cid.Cid{}, fmt.Errorf("failed to write block: %w", err)
		}

		if err == io.EOF {
			break
		}
	}
	return carFile.Name(), rootCid, nil
}

func createCID(data []byte) cid.Cid {
	mh, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		return cid.Cid{}
	}
	return cid.NewCidV1(cid.Raw, mh)
}

type RepeatedReader struct {
	size     int64
	remain   int64
	block    []byte
	blockNum uint64
}

func NewRepeatedReader(size, blockSize int64) *RepeatedReader {
	return &RepeatedReader{
		size:   size,
		remain: size,
		block:  make([]byte, blockSize),
	}
}

func (r *RepeatedReader) Read(p []byte) (n int, err error) {
	if r.remain <= 0 {
		return 0, io.EOF
	}

	toRead := len(p)
	if int64(toRead) > r.remain {
		toRead = int(r.remain)
	}

	r.block[0] = byte(r.blockNum)
	r.block[1] = byte(r.blockNum >> 8)
	r.block[2] = byte(r.blockNum >> 16)
	r.block[3] = byte(r.blockNum >> 24)
	r.block[4] = byte(r.blockNum >> 32)
	r.block[5] = byte(r.blockNum >> 40)
	r.block[6] = byte(r.blockNum >> 48)
	r.block[7] = byte(r.blockNum >> 56)
	r.blockNum++

	copy(p[:toRead], r.block[:toRead])
	r.remain -= int64(toRead)

	return toRead, nil
}
