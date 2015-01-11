package fulltext

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/jbarham/go-cdb"
	"io"
	"io/ioutil"
	"os"
	"syscall"
)

// Size of header block to prepend - make it 4k to align disk reads
const HEADER_SIZE = 4096

// Produces a set of cdb files from a series of AddDoc() calls
type Indexer struct {
	docTxtFile    *os.File
	wordTxtFile   *os.File
	docCdbFile    *os.File
	wordCdbFile   *os.File
	wordMap       map[string]map[string]int // map of [word][docId]count
	WordSplit     WordSplitter
	WordClean     WordCleaner
	StopWordCheck StopWordChecker
}

// Contents of a single document to be indexed
type IndexDoc struct {
	Id         []byte // the id, this is usually the path to the document
	IndexValue []byte // index this data
	StoreValue []byte // store this data
}

// Creates a new indexer, using the given temp dir while building
// the index.
func NewIndexer(tempDir string) (*Indexer, error) {
	idx := &Indexer{}
	var err error
	idx.docTxtFile, err = ioutil.TempFile(tempDir, "doctmp")
	if err != nil {
		return nil, err
	}
	idx.wordTxtFile, err = ioutil.TempFile(tempDir, "wordtmp")
	if err != nil {
		return nil, err
	}
	idx.docCdbFile, err = ioutil.TempFile(tempDir, "doccdb")
	if err != nil {
		return nil, err
	}
	idx.wordCdbFile, err = ioutil.TempFile(tempDir, "wordcdb")
	if err != nil {
		return nil, err
	}
	idx.wordMap = make(map[string]map[string]int)
	idx.WordSplit = Wordize
	idx.WordClean = IndexizeWord
	return idx, nil
}

// Add a document to the index - writes to temporary files and stores some data in memory while building the index.
func (idx *Indexer) AddDoc(idoc IndexDoc) error {
	// add to docs
	docId := string(idoc.Id)
	writeTextLine(idx.docTxtFile, []byte(docId), idoc.StoreValue)
	words := append(idx.WordSplit(string(idoc.IndexValue)), idx.WordSplit(string(idoc.StoreValue))...)
	for _, word := range words {
		word = idx.WordClean(word)

		// skip if stop word
		if idx.StopWordCheck != nil {
			if idx.StopWordCheck(word) {
				continue
			}
		}

		// ensure nested map exists
		if idx.wordMap[word] == nil {
			idx.wordMap[word] = make(map[string]int)
		}
		// increment count by one for this combination
		c := idx.wordMap[word][docId] + 1
		idx.wordMap[word][docId] = c
	}
	return nil
}

// Builds a final single index file, which consists of some simple header info,
// followed by the cdb binary files that comprise the full index.
func (idx *Indexer) FinalizeAndWrite(w io.Writer) error {

	var buf bytes.Buffer

	// write out the word data
	for word, m := range idx.wordMap {
		enc := gob.NewEncoder(&buf)
		enc.Encode(m)
		writeTextLine(idx.wordTxtFile, []byte(word), buf.Bytes())
		buf.Reset()
	}

	var err error

	idx.docTxtFile.Write([]byte("\n"))
	idx.wordTxtFile.Write([]byte("\n"))

	_, err = idx.docTxtFile.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = idx.wordTxtFile.Seek(0, 0)
	if err != nil {
		return err
	}

	// make cdb files
	err = cdb.Make(idx.docCdbFile, idx.docTxtFile)
	if err != nil {
		return err
	}
	err = cdb.Make(idx.wordCdbFile, idx.wordTxtFile)
	if err != nil {
		return err
	}

	// make sure the contents are all settled
	idx.docCdbFile.Sync()
	idx.wordCdbFile.Sync()
	_, err = idx.docCdbFile.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = idx.wordCdbFile.Seek(0, 0)
	if err != nil {
		return err
	}

	docStat, err := idx.docCdbFile.Stat()
	if err != nil {
		return err
	}
	wordStat, err := idx.wordCdbFile.Stat()
	if err != nil {
		return err
	}

	// now package it all up
	buf.Reset()
	enc := gob.NewEncoder(&buf)
	bhead := []int{int(docStat.Size()), int(wordStat.Size())}
	enc.Encode(bhead)

	// extend buffer to be HEADER_SIZE len
	bpadsize := HEADER_SIZE - buf.Len()
	buf.Write(make([]byte, bpadsize, bpadsize))
	b := buf.Bytes()

	_, err = w.Write(b)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, idx.docCdbFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, idx.wordCdbFile)
	if err != nil {
		return err
	}

	return nil
}

// Dump some human readable status information
func (idx *Indexer) DumpStatus(w io.Writer) {
	fmt.Fprintf(w, "files used:\n\t%s\n\t%s\n\t%s\n\t%s\n", idx.docTxtFile.Name(), idx.wordTxtFile.Name(), idx.docCdbFile.Name(), idx.wordCdbFile.Name())
	// fmt.Fprintf(w, "wordMap: %+v\n", idx.wordMap)
}

// close and remove all resources
func (idx *Indexer) Close() {
	syscall.Unlink(idx.docTxtFile.Name())
	idx.docTxtFile.Close()
	syscall.Unlink(idx.wordTxtFile.Name())
	idx.wordTxtFile.Close()
	syscall.Unlink(idx.docCdbFile.Name())
	idx.docCdbFile.Close()
	syscall.Unlink(idx.wordCdbFile.Name())
	idx.wordCdbFile.Close()
	idx.wordMap = nil
}

// Write a single line of data in cdb's text format
func writeTextLine(w io.Writer, key []byte, data []byte) (err error) {
	_, err = fmt.Fprintf(w, "+%d,%d:%s->%s\n", len(key), len(data), key, data)
	return
}
