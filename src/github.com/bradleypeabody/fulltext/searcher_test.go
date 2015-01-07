package fulltext

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	re "regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Extract a single file from a zip and return it's contents
func zipExtract(zfpath string, fpath string) ([]byte, error) {

	zr, err := zip.OpenReader(zfpath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()

	fpath = strings.Trim(filepath.Clean(filepath.ToSlash(fpath)), "/")

	for _, f := range zr.File {

		fn := strings.Trim(filepath.Clean(filepath.ToSlash(f.Name)), "/")

		// keep going until we find it
		if fn != fpath {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		rc.Close()

		return b, nil

	}

	return nil, io.EOF

}

// Index and search the complete works of William Shakespeare
func TestTheBardSearch(t *testing.T) {

	fmt.Println("TestTheBardIndexing")

	idx, err := NewIndexer("")
	if err != nil {
		panic(err)
	}
	defer idx.Close()

	// use English stop words
	idx.StopWordCheck = EnglishStopWordChecker

	titlere := re.MustCompile("(?i)<title>([^<]+)</title>")

	zr, err := zip.OpenReader("testdata/shakespeare.mit.edu.zip")
	if err != nil {
		panic(err)
	}
	defer zr.Close()

	for _, f := range zr.File {
		fmt.Printf("indexing: %s\n", f.Name)

		rc, err := f.Open()
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(rc)
		if err != nil {
			panic(err)
		}

		// extract title tag
		tret := titlere.FindSubmatch(b)
		title := ""
		if len(tret) > 1 {
			title = strings.TrimSpace(string(tret[1]))
		}

		// strip html from entire doc and get text
		body := HTMLStripTags(string(b))

		// make a doc out of it
		doc := IndexDoc{
			Id:         []byte(f.Name),
			StoreValue: []byte(title),
			IndexValue: []byte(title + " " + title + " " + body),
		}
		idx.AddDoc(doc)

		rc.Close()
	}

	fmt.Println("Writing final index...")
	f, err := ioutil.TempFile("", "idxout")
	if err != nil {
		panic(err)
	}
	err = idx.FinalizeAndWrite(f)
	if err != nil {
		panic(err)
	}

	fmt.Println("Debug data: \n")
	idx.DumpStatus(os.Stdout)

	// panic("DONE")

	f.Close()

	fmt.Printf("Wrote index file: %s\n", f.Name())

	/////////////////////////////////

	start := time.Now()

	s, err := NewSearcher(f.Name())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Opening searcher took: %s\n", time.Since(start).String())

	start = time.Now()

	sr, err := s.SimpleSearch("king", 20)
	if err != nil {
		panic(err)
	}

	if len(sr.Items) == 0 {
		t.Fatalf("Search for 'king' returned 0 results, but should have gotten something")
	}

	fmt.Printf("Searching took: %s\n", time.Since(start).String())

	fmt.Printf("Total Results for 'king': %d\n", len(sr.Items))
	for k, v := range sr.Items {
		fmt.Printf("----------- #:%d\n", k)
		fmt.Printf("Id: %s\n", v.Id)
		fmt.Printf("Score: %d\n", v.Score)
		fmt.Printf("StoreValue: %s\n", v.StoreValue)
	}

	fmt.Printf("Raw dump: %+v\n", sr)

	// look for a stop word and make sure it's not there

	sr, err = s.SimpleSearch("the", 20)
	if err != nil {
		panic(err)
	}
	if len(sr.Items) != 0 {
		t.Fatalf("Search for 'the' returned %d results when it should have been 0 because it's a stop word", len(sr.Items))
	}
	fmt.Printf("Check for stop word passed\n")

	///////////////////////////////////////////////////

	fmt.Printf("Starting Shakespeare's very own search interface at :1414 ...")

	ln, err := net.Listen("tcp", ":1414")
	if err != nil {
		panic(err)
	}

	timeoutStr := os.Getenv("SEARCHER_WEB_TIMEOUT_SECONDS")

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeout = 10
	}

	zfpath := "testdata/shakespeare.mit.edu.zip"

	// wait for specified time
	go func() { time.Sleep(time.Duration(timeout) * time.Second); ln.Close() }()

	// main request handler
	err = http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// home page redirect
		if r.URL.Path == "/" || r.URL.Path == "/Shakespeare" {
			http.Redirect(w, r, "/shakespeare.mit.edu/index.html", 302)
			return
		}

		// handle search result page
		if r.URL.Path == "/searchresults.html" {

			w.Header().Set("Content-type", "text/html")

			q := r.FormValue("q")

			// do search
			sr, err := s.SimpleSearch(q, 20)
			if err != nil {
				panic(err)
			}

			// render results page
			sres, err := ioutil.ReadFile("testdata/searchresults.html")
			if err != nil {
				panic(err)
			}
			t := template.Must(template.New("main").Parse(string(sres)))
			var buf bytes.Buffer
			t.Execute(&buf, &map[string]interface{}{
				"q":  q,
				"sr": sr,
			})
			sresbytes := buf.Bytes()

			w.Write(sresbytes)

			return
		}

		// by default look through zip file
		b, err := zipExtract(zfpath, r.URL.Path)
		if err != nil {
			http.Error(w, "File not found", 404)
		}
		if strings.HasSuffix(r.URL.Path, ".css") {
			w.Header().Set("Content-type", "text/css")
		}
		if strings.HasSuffix(r.URL.Path, ".gif") {
			w.Header().Set("Content-type", "image/gif")
		}
		if strings.HasSuffix(r.URL.Path, ".jpg") {
			w.Header().Set("Content-type", "image/jpeg")
		}

		// for html files we inject a search box
		if strings.HasSuffix(r.URL.Path, ".html") {
			w.Header().Set("Content-type", "text/html")

			// render search form
			sf, err := ioutil.ReadFile("testdata/searchform.html")
			if err != nil {
				panic(err)
			}
			t := template.Must(template.New("main").Parse(string(sf)))
			var buf bytes.Buffer
			t.Execute(&buf, r.FormValue("q"))
			sfbytes := buf.Bytes()

			// inject into page

			pagebytes := re.MustCompile("(<body[^>]*>)").ReplaceAllLiteral(b, []byte("<body bgcolor=\"#ffffff\" text=\"#000000\">"+string(sfbytes)))
			w.Write(pagebytes)
			return

		}

		w.Write(b)

	}))

	if err != nil {
		fmt.Printf("err from listen: %s\n", err)
	}

	s.Close()

}
