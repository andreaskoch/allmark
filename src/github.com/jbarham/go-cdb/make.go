package cdb

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"strconv"
)

var BadFormatError = errors.New("bad format")

// Make reads cdb-formatted records from r and writes a cdb-format database
// to w.  See the documentation for Dump for details on the input record format. 
func Make(w io.WriteSeeker, r io.Reader) (err error) {
	defer func() { // Centralize error handling.
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	if _, err = w.Seek(int64(headerSize), 0); err != nil {
		return
	}

	buf := make([]byte, 8)
	rb := bufio.NewReader(r)
	wb := bufio.NewWriter(w)
	hash := cdbHash()
	hw := io.MultiWriter(hash, wb) // Computes hash when writing record key.
	rr := &recReader{rb}
	htables := make(map[uint32][]slot)
	pos := headerSize
	// Read all records and write to output.
	for {
		// Record format is "+klen,dlen:key->data\n"
		c := rr.readByte()
		if c == '\n' { // end of records
			break
		}
		if c != '+' {
			return BadFormatError
		}
		klen, dlen := rr.readNum(','), rr.readNum(':')
		writeNums(wb, klen, dlen, buf)
		hash.Reset()
		rr.copyn(hw, klen)
		rr.eatByte('-')
		rr.eatByte('>')
		rr.copyn(wb, dlen)
		rr.eatByte('\n')
		h := hash.Sum32()
		tableNum := h % 256
		htables[tableNum] = append(htables[tableNum], slot{h, pos})
		pos += 8 + klen + dlen
	}

	// Write hash tables and header.

	// Create and reuse a single hash table.
	maxSlots := 0
	for _, slots := range htables {
		if len(slots) > maxSlots {
			maxSlots = len(slots)
		}
	}
	slotTable := make([]slot, maxSlots*2)

	header := make([]byte, headerSize)
	// Write hash tables.
	for i := uint32(0); i < 256; i++ {
		slots := htables[i]
		if slots == nil {
			putNum(header[i*8:], pos)
			continue
		}

		nslots := uint32(len(slots) * 2)
		hashSlotTable := slotTable[:nslots]
		// Reset table slots.
		for j := 0; j < len(hashSlotTable); j++ {
			hashSlotTable[j].h = 0
			hashSlotTable[j].pos = 0
		}

		for _, slot := range slots {
			slotPos := (slot.h / 256) % nslots
			for hashSlotTable[slotPos].pos != 0 {
				slotPos++
				if slotPos == uint32(len(hashSlotTable)) {
					slotPos = 0
				}
			}
			hashSlotTable[slotPos] = slot
		}

		if err = writeSlots(wb, hashSlotTable, buf); err != nil {
			return
		}

		putNum(header[i*8:], pos)
		putNum(header[i*8+4:], nslots)
		pos += 8 * nslots
	}

	if err = wb.Flush(); err != nil {
		return
	}

	if _, err = w.Seek(0, 0); err != nil {
		return
	}

	_, err = w.Write(header)

	return
}

type recReader struct {
	*bufio.Reader
}

func (rr *recReader) readByte() byte {
	c, err := rr.ReadByte()
	if err != nil {
		panic(err)
	}

	return c
}

func (rr *recReader) eatByte(c byte) {
	if rr.readByte() != c {
		panic(errors.New("unexpected character"))
	}
}

func (rr *recReader) readNum(delim byte) uint32 {
	s, err := rr.ReadString(delim)
	if err != nil {
		panic(err)
	}

	s = s[:len(s)-1] // Strip delim
	n, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		panic(err)
	}

	return uint32(n)
}

func (rr *recReader) copyn(w io.Writer, n uint32) {
	if _, err := io.CopyN(w, rr, int64(n)); err != nil {
		panic(err)
	}
}

func putNum(buf []byte, x uint32) {
	binary.LittleEndian.PutUint32(buf, x)
}

func writeNums(w io.Writer, x, y uint32, buf []byte) {
	putNum(buf, x)
	putNum(buf[4:], y)
	if _, err := w.Write(buf[:8]); err != nil {
		panic(err)
	}
}

type slot struct {
	h, pos uint32
}

func writeSlots(w io.Writer, slots []slot, buf []byte) (err error) {
	for _, np := range slots {
		putNum(buf, np.h)
		putNum(buf[4:], np.pos)
		if _, err = w.Write(buf[:8]); err != nil {
			return
		}
	}

	return nil
}
