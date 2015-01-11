Overview
========

This is a simple, pure-Go, full text indexing and search library.

I made it for use on small to medium websites, although there is nothing web-specific about it's API or operation.

Cdb (http://github.com/jbarham/go-cdb) is used to perform the indexing and lookups.

Status
------

This project is experimental.  Breaking changes very well may occur.

Notes on Building
--------

fulltext requires CDB:

	go get github.com/jbarham/go-cdb

Usage
------

First, you must create an index.  Like this:

	import "github.com/bradleypeabody/fulltext"

	// create new index with temp dir (usually "" is fine)
	idx, err := fulltext.NewIndexer(""); if err != nil { panic(err) }
	defer idx.Close()

	// provide stop words if desired
	idx.StopWordCheck = fulltext.EnglishStopWordChecker

	// for each document you want to add, you do something like this:
	doc := fulltext.IndexDoc{
		Id: []byte(uuid), // unique identifier (the path to a webpage works...)
		StoreValue: []byte(title), // bytes you want to be able to retrieve from search results
		IndexValue: []byte(data), // bytes you want to be split into words and indexed
	}
	idx.AddDoc(doc) // add it

	// when done, write out to final index
	err = idx.FinalizeAndWrite(f); if err != nil { panic(err) }

Once you have an index file, you can search it like this:

	s, err := fulltext.NewSearcher("/path/to/index/file"); if err != nil { panic(err) }
	defer s.Close()
	sr, err := s.SimpleSearch("Horatio", 20); if err != nil { panic(err) }
	for k, v := range sr.Items {
		fmt.Printf("----------- #:%d\n", k)
		fmt.Printf("Id: %s\n", v.Id)
		fmt.Printf("Score: %d\n", v.Score)
		fmt.Printf("StoreValue: %s\n", v.StoreValue)
	}

It's rather simplistic.  But it's fast and it works.

TODOs
-----

* ~~Will likely need some sort of "stop word" functionality.~~

* ~~Wordize(), IndexizeWord()~~ and the scoring aggregation logic should be extracted to callback functions with the existing functionality as default.

* The search logic is currently very naive.  Ideally this project would have something as sophisticated as <a href="http://lucene.apache.org/core/4_10_0/queryparser/org/apache/lucene/queryparser/classic/package-summary.html" target="_blank">Lucene's query parser</a>.  But in reality what I'll likely do is a simple survey of which common features are actually used on any on-site search engines I can get my hands on.  Quoting ("black cat"), and logical operators (Jim OR James) would likely be at the top of the list and implementing that sort of thing would be higher priority than trying to duplicate Lucene.

* If there is some decent b-tree disk storage that is portable then it would be worth looking at using that instead of CDB and implementing LIKE-style matching.  As it is, CDB is quite efficient, but it is a hash index.


Implementation Notes
--------------------

I originally tried doing this on top of Sqlite.  It was dreadfully slow.  Cdb is orders of magnitude faster.

Two main disadvantages from going the Cdb route are that the index cannot be edited once it is built (you have to recreate it in full), and since it's hash-based it will not support any sort of fuzzy matching unless those variations are included in the index (which they are not, in the current implementation.)   For my purposes these two disadvantages are overshadowed by the fact that it's blinding fast, easy to use, portable (pure-Go), and its interface allowed me to build the indexes I needed into a single file.

In the test suite is included a copy of the complete works of William Shakespeare (thanks to Jeremy Hylton's http://shakespeare.mit.edu/) and this library is used to create a simple search engine on top of that corpus.  By default it only runs for 10 seconds, but you can run it for longer by doing something like:

	SEARCHER_WEB_TIMEOUT_SECONDS=120 go test fulltext -v

Works on Windows.

Future Work
-----------

It might be feasible to supplant this project with something using suffix arrays ( http://golang.org/pkg/index/suffixarray/ ).  The main down side would be the requirement of a lot more storage space (and memory to load and search it).  Retooling the index/suffixarray package so it can work against the disk is an idea, but is not necessarily simple.  The upside of an approach like that would be full regex support for searches with decent performance - which would rock.  The index could potentially be sharded by the first character or two of the search - but that's still not as good as something with sensible caching where the whole set can be kept on disk and the "hot" parts cached in memory, etc.
