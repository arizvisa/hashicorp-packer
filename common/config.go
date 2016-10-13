package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"net/url"
)

// ScrubConfig is a helper that returns a string representation of
// any struct with the given values stripped out.
func ScrubConfig(target interface{}, values ...string) string {
	conf := fmt.Sprintf("Config: %+v", target)
	for _, value := range values {
		if value == "" {
			continue
		}
		conf = strings.Replace(conf, value, "<Filtered>", -1)
	}
	return conf
}

// ChooseString returns the first non-empty value.
func ChooseString(vals ...string) string {
	for _, el := range vals {
		if el != "" {
			return el
		}
	}

	return ""
}

// DownloadableURL processes a URL that may also be a file path and returns
// a completely valid URL. For example, the original URL might be "local/file.iso"
// which isn't a valid URL. DownloadableURL will return "file:///local/file.iso"
func DownloadableURL(original string) (string, error) {

	// Verify that the scheme is something we support in our common downloader.
	supported := []string{"file", "http", "https", "ftp", "smb"}
	found := false
	for _, s := range supported {
		if strings.HasPrefix(strings.ToLower(original), s + "://") {
			found = true
			break
		}
	}

	// If it's properly prefixed with something we support, then we don't need
	//	to make it a uri.
	if found {
		original = filepath.ToSlash(original)

		// make sure that it can be parsed though..
		uri,err := url.Parse(original)
		if err != nil { return "", err }

		uri.Scheme = strings.ToLower(uri.Scheme)

		return uri.String(), nil
	}

	// If the file exists, then make it an absolute path
	_,err := os.Stat(original)
	if err == nil {
		original, err = filepath.Abs(filepath.FromSlash(original))
		if err != nil { return "", err }

		original, err = filepath.EvalSymlinks(original)
		if err != nil { return "", err }

		original = filepath.Clean(original)
		original  = filepath.ToSlash(original)
	}

	// Since it wasn't properly prefixed, let's make it into a well-formed
	//	file:// uri.

	return "file://" + original, nil
}
