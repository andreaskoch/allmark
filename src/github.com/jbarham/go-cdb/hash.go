package cdb

import "hash"

const (
	start = 5381 // Initial cdb checksum value.
)

// digest represents the partial evaluation of a checksum.
type digest struct {
	h uint32
}

func (d *digest) Reset() { d.h = start }

// New returns a new hash computing the cdb checksum.
func cdbHash() hash.Hash32 {
	d := new(digest)
	d.Reset()
	return d
}

func (d *digest) Size() int { return 4 }

func update(h uint32, p []byte) uint32 {
	for i := 0; i < len(p); i++ {
		h = ((h << 5) + h) ^ uint32(p[i])
	}
	return h
}

func (d *digest) Write(p []byte) (int, error) {
	d.h = update(d.h, p)
	return len(p), nil
}

func (d *digest) Sum32() uint32 { return d.h }

func (d *digest) Sum(in []byte) []byte {
	s := d.Sum32()
	in = append(in, byte(s>>24))
	in = append(in, byte(s>>16))
	in = append(in, byte(s>>8))
	in = append(in, byte(s))
	return in
}

func (d *digest) BlockSize() int { return 1 }

func checksum(data []byte) uint32 { return update(start, data) }
