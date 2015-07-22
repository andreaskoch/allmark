package fulltext

import (
	"fmt"
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
	"path/filepath"
	re "regexp"
	"testing"
)

func TestIndexer(t *testing.T) {
	fmt.Printf("TestIndexer\n")

	idx, err := NewIndexer()
	if err != nil {
		panic(err)
	}

	idx.AddDoc(IndexDoc{Id: []byte(`blah1`), StoreValue: []byte(`store this`), IndexValue: []byte(`test of the emergency broadcast system`)})
	idx.AddDoc(IndexDoc{Id: []byte(`blah2`), StoreValue: []byte(`store this stuff too, yeah store it`), IndexValue: []byte(`every good boy does fine`)})
	idx.AddDoc(IndexDoc{Id: []byte(`blah3`), StoreValue: []byte(`more storage here`), IndexValue: []byte(`a taco in the hand is worth two in the truck`)})

	var indexFs afero.Fs = &afero.MemMapFs{}

	f, err := indexFs.Create("idxout")
	if err != nil {
		panic(err)
	}
	err = idx.FinalizeAndWrite(f)
	if err != nil {
		panic(err)
	}
	f.Close()

	fmt.Printf("Wrote index file: %s\n", f.Name())

}

// A more extensive test - index the complete works of William Shakespeare
func NoTestTheBardIndexing(t *testing.T) {

	fmt.Println("TestTheBardIndexing")

	idx, err := NewIndexer()
	if err != nil {
		panic(err)
	}
	defer idx.Close()

	titlere := re.MustCompile("(?i)<title>([^<]+)</title>")

	n := 0

	filepath.Walk("testdata/shakespeare.mit.edu/", func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() /*&& n < 5*/ {
			n++
			fmt.Printf("indexing: %s\n", path)
			b, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			title := string(titlere.Find(b))
			body := HTMLStripTags(string(b))
			doc := IndexDoc{
				Id:         []byte(path),
				StoreValue: []byte(title),
				IndexValue: []byte(title + " " + title + " " + body),
			}
			idx.AddDoc(doc)
		}
		return nil
	})

	fmt.Println("Writing final index...")
	f, err := ioutil.TempFile("", "idxout")
	if err != nil {
		panic(err)
	}
	err = idx.FinalizeAndWrite(f)
	if err != nil {
		panic(err)
	}
	f.Close()

	fmt.Printf("Wrote index file: %s\n", f.Name())

}
