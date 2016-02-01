package feedparser

import "testing"

func TestKeywordSet(t *testing.T) {
	k := make(keywordSet)

	k.add("foo")
	k.add("bar")
	k.add("foo") // should be unique
	k.add("czah")

	s := k.slice()
	if len(s) != 3 {
		t.Fatal("Should be unique 3 items")
	}

}
