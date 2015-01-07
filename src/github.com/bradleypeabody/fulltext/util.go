package fulltext

import (
	"bytes"
	"html/template"
	re "regexp"
	"strings"
)

var wordizeRe *re.Regexp

func init() {
	wordizeRe = re.MustCompile("[\\s,.;:!?[\\]()'\"]+")
}

type WordSplitter func(string) []string

// Split a string up into words
func Wordize(t string) []string {
	return wordizeRe.Split(t, -1)
}

type WordCleaner func(string) string

// Make word appropriate for indexing
func IndexizeWord(w string) string {
	return strings.TrimSpace(strings.ToLower(w))
}

type StopWordChecker func(string) bool

// This function copied from here: https://github.com/kennygrant/sanitize/blob/master/sanitize.go
// License is: https://github.com/kennygrant/sanitize/blob/master/License-BSD.txt
// Strip html tags, replace common entities, and escape <>&;'" in the result.
// Note the returned text may contain entities as it is escaped by HTMLEscapeString,
// and most entities are not translated.
func HTMLStripTags(s string) (output string) {

	output = ""

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {

		// First remove line breaks etc as these have no meaning outside html tags (except pre)
		// this means pre sections will lose formatting... but will result in less uninentional paras.
		s = strings.Replace(s, "\n", "", -1)

		// Then replace line breaks with newlines, to preserve that formatting
		s = strings.Replace(s, "</p>", "\n", -1)
		s = strings.Replace(s, "<br>", "\n", -1)
		s = strings.Replace(s, "</br>", "\n", -1)

		// Walk through the string removing all tags
		b := bytes.NewBufferString("")
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}

	// In case we have missed any tags above, escape the text - removes <, >, &, ' and ".
	output = template.HTMLEscapeString(output)

	// Remove a few common harmless entities, to arrive at something more like plain text
	// This relies on having removed *all* tags above
	output = strings.Replace(output, "&nbsp;", " ", -1)
	output = strings.Replace(output, "&quot;", "\"", -1)
	output = strings.Replace(output, "&apos;", "'", -1)
	output = strings.Replace(output, "&#34;", "\"", -1)
	output = strings.Replace(output, "&#39;", "'", -1)
	// NB spaces here are significant - we only allow & not part of entity
	output = strings.Replace(output, "&amp; ", "& ", -1)
	output = strings.Replace(output, "&amp;amp; ", "& ", -1)

	return output
}

var titlere *re.Regexp
var descre *re.Regexp

func init() {
	titlere = re.MustCompile("(?i)<title>([^<]+)</title>")
	descre = re.MustCompile("(?i)<meta\\s+name\\s*=\\s*[\"']description[\"']\\s+content\\s*=\\s*[\"'](.*?)[\"']\\s*/?>")
}

// Helper to extract an HTML title from the title tag
func HTMLExtractTitle(html string) string {
	tret := titlere.FindSubmatch([]byte(html))
	title := ""
	if len(tret) > 1 {
		title = strings.TrimSpace(string(tret[1]))
	}
	return title
}

// Helper to extract an HTML description from the meta[name=description] tag
func HTMLExtractDescription(html string) string {
	tret := descre.FindSubmatch([]byte(html))
	desc := ""
	if len(tret) > 1 {
		desc = strings.TrimSpace(string(tret[1]))
	}
	return desc
}
