package fulltext

import (
	"bytes"
	"encoding/gob"
	"github.com/jbarham/go-cdb"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

// Interface for search.  Not thread-safe, but low overhead
// so having a separate one per thread should be workable.
type Searcher struct {
	file    *os.File
	docCdb  *cdb.Cdb
	wordCdb *cdb.Cdb
}

// Wraps a ReaderAt and adjusts (tweaks) it's offset by the specified amount
type tweakedReaderAt struct {
	readerAt io.ReaderAt
	tweak    int64
}

func (t *tweakedReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	n, err = t.readerAt.ReadAt(p, off+t.tweak)
	return
}

// A single item in a search result
type SearchResultItem struct {
	Id         []byte // id of this item (document)
	StoreValue []byte // the stored value of this document
	Score      int64  // the total score
}

// Implement sort.Interface
type SearchResultItems []SearchResultItem

func (s SearchResultItems) Len() int      { return len(s) }
func (s SearchResultItems) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SearchResultItems) Less(i, j int) bool {
	// if same score, then sort by raw bytes comparison of store value -
	// so we get consistently ordered results, even when score is same
	if s[i].Score == s[j].Score {
		return bytes.Compare(s[i].Id, s[j].Id) < 0
	}
	return s[i].Score < s[j].Score
}

// What happened during the search
type SearchResults struct {
	Items SearchResultItems
}

// Make a new searcher using the file at the specified path
// TODO: Make a variation that accepts a ReaderAt
func NewSearcher(fpath string) (*Searcher, error) {

	s := &Searcher{}

	f, err := os.Open(fpath)
	if err != nil {
		return s, err
	}
	s.file = f

	// write out the word data
	dec := gob.NewDecoder(f)
	lens := make([]int64, 2, 2)
	dec.Decode(&lens)

	s.docCdb = cdb.New(&tweakedReaderAt{f, HEADER_SIZE})
	s.wordCdb = cdb.New(&tweakedReaderAt{f, HEADER_SIZE + lens[0]})

	return s, nil
}

// Close and release resources
func (s *Searcher) Close() error {
	s.docCdb = nil
	s.wordCdb = nil
	return s.file.Close()
}

// Perform a search
func (s *Searcher) SimpleSearch(search string, maxn int) (SearchResults, error) {

	sr := SearchResults{}

	// break search into words_word
	searchWords := Wordize(search)

	itemMap := make(map[string]SearchResultItem)

	// read word data for each word that was provided
	for _, w := range searchWords {
		w = IndexizeWord(w)
		// find the docs for this word
		mapGob, err := s.wordCdb.Find([]byte(w))
		if err == io.EOF {
			continue
		}
		if err != nil {
			return sr, err
		}

		m := make(map[string]int)

		dec := gob.NewDecoder(mapGob)
		err = dec.Decode(&m)
		if err != nil {
			return sr, err
		}

		// for each doc, increase score
		for docId, cnt := range m {
			sri := itemMap[docId]
			if sri.Score < 1 {
				sri.Id = []byte(docId)
			}
			sri.Score += int64(cnt)
			itemMap[docId] = sri
		}

	}

	// convert to slice
	items := make(SearchResultItems, 0, maxn)
	for _, item := range itemMap {
		items = append(items, item)
	}

	// sort by score descending
	sort.Sort(sort.Reverse(items))

	// limit to maxn
	if len(items) > maxn {
		items = items[:maxn]
	}

	// pull document contents from doc cdb
	for i := range items {
		item := &items[i]
		v, err := s.docCdb.Find(item.Id)
		if err == io.EOF {
			panic("doc id " + string(item.Id) + " not found in index, this should never happen")
		}
		if err != nil {
			return sr, err
		}
		v1, err := ioutil.ReadAll(v)
		if err != nil {
			return sr, err
		}
		item.StoreValue = v1
	}

	sr.Items = items

	return sr, nil

}
