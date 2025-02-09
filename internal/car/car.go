package car

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/filecoin-project/go-commp-utils/v2/writer"
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
	size       int64
	remain     int64
	buffer     []byte
	bufferSize int
	offset     int
}

const (
	DefaultBufferSize = 16*1024*1024 + 1
)

func NewRepeatedReader(size int64, chunkSize int64) *RepeatedReader {
	buffer := make([]byte, DefaultBufferSize)
	rand.Read(buffer)

	return &RepeatedReader{
		size:       size,
		remain:     size,
		buffer:     buffer,
		bufferSize: DefaultBufferSize,
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

	written := 0
	for written < toRead {
		available := r.bufferSize - r.offset
		shouldRead := toRead - written
		if shouldRead > available {
			shouldRead = available
		}

		copy(p[written:written+shouldRead], r.buffer[r.offset:r.offset+shouldRead])
		written += shouldRead

		r.offset += shouldRead
		if r.offset >= r.bufferSize {
			r.offset = 0
		}
	}

	r.remain -= int64(written)
	return written, nil
}

func CommpReader(rdr io.Reader) (*writer.DataCIDSize, error) {
	w := &writer.Writer{}
	_, err := io.CopyBuffer(w, rdr, make([]byte, writer.CommPBuf))
	if err != nil {
		return nil, fmt.Errorf("copy into commp writer: %w", err)
	}

	cp, err := w.Sum()
	if err != nil {
		return nil, fmt.Errorf("computing commP failed: %w", err)
	}

	return &cp, nil
}

func Commp(inPath string) (*writer.DataCIDSize, error) {
	rdr, err := os.Open(inPath)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()

	return CommpReader(rdr)
}
