package util

import (
	"io/ioutil"
	"net/http"
	"os"
)

func DownloadFile(url string, auth string) (string, func(), error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Add("Authorization", auth)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}

	tmpdir, err := ioutil.TempDir("", "downloaded-files-")
	if err != nil {
		return "", nil, err
	}

	cleaner := func() {
		os.RemoveAll(tmpdir)
	}

	tmpFile := tmpdir + "/file"
	err = ioutil.WriteFile(tmpFile, content, 0666)
	if err != nil {
		os.RemoveAll(tmpdir)
		return "", nil, err
	}

	return tmpFile, cleaner, nil
}
