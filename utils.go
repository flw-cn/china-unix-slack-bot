package main

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/flw-cn/slack"
)

func slackDownloadFile(tmpdir string, api *slack.Client, fileID string) (string, func(), error) {
	file, _, _, err := api.GetFileInfo(fileID, 0, 0)
	if err != nil {
		return "", nil, err
	}

	url := file.URLPrivateDownload

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", nil, err
	}

	req.Header.Add("Authorization", "Bearer "+config.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}

	tmpdir, err = ioutil.TempDir(tmpdir, "slack-bot-downloaded-files-")
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
