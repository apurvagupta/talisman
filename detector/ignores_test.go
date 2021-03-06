package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thoughtworks/talisman/git_repo"
)

func TestShouldIgnoreEmptyLinesInTheFile(t *testing.T) {
	for _, s := range []string{"", " ", "  "} {
		assert.True(t, NewIgnores(s).AcceptsAll(), "Expected '%s' to result in no ignore patterns.", s)
	}
}

func TestShouldIgnoreLinesThatBeginWithThePound(t *testing.T) {
	for _, s := range []string{"#", "#monkey", "# this monkey likes bananas  "} {
		assert.True(t, NewIgnores(s).AcceptsAll(), "Expected commented line '%s' to result in no ignore patterns", s)
	}
}

func TestRawPatterns(t *testing.T) {
	assertAccepts("foo", "bar", t)
	assertAccepts("foo", "foobar", t)
	assertAccepts("foo", "foo/bar", t)
	assertAccepts("foo", "foo/bar/baz", t)

	assertDenies("foo", "foo", t)
	assertDenies("foo", "bar/foo", t)
	assertDenies("foo", "bar/baz/foo", t)
}

func TestSingleStarPatterns(t *testing.T) {
	assertAccepts("foo*", "bar", t)
	assertAccepts("foo*", "foo/bar", t)
	assertAccepts("foo*", "foo/bar/baz", t)

	assertDenies("foo*", "foo", t)
	assertDenies("foo*", "foobar", t)

	assertAccepts("*.pem", "foo.txt", t)
	assertAccepts("*.pem", "pem.go", t)
	assertDenies("*.pem", "secret.pem", t)
	assertDenies("*.pem", "foo/bar/secret.pem", t)
}

func TestDirectoryPatterns(t *testing.T) {
	assertAccepts("foo/", "bar", t)
	assertAccepts("foo/", "foo", t)
	assertDenies("foo/", "foo/bar", t)
	assertDenies("foo/", "foo/bar.txt", t)
	assertDenies("foo/", "foo/bar/baz.txt", t)
}

func assertDenies(pattern, path string, t *testing.T) {
	assert.True(t, NewIgnores(pattern).Deny(testAddition(path)), "%s is expected to deny a file named %s.", pattern, path)
}

func assertAccepts(pattern, path string, t *testing.T) {
	assert.True(t, NewIgnores(pattern).Accept(testAddition(path)), "%s is expected to accept a file named %s.", pattern, path)
}

func testAddition(path string) git_repo.Addition {
	return git_repo.NewAddition(path, make([]byte, 0))
}
