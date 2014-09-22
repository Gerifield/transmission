package transmission

import (
	"encoding/base64"
	"encoding/json"
	"github.com/longnguyen11288/go-transmission/client"
	"io/ioutil"
)

var ()

type TransmissionClient struct {
	apiclient client.ApiClient
}

type command struct {
	Method    string    `json:"method,omitempty"`
	Arguments arguments `json:"arguments,omitempty"`
	Result    string    `json:"result,omitempty"`
}

type arguments struct {
	Fields       []string     `json:"fields,omitempty"`
	Torrents     []Torrent    `json:"torrents,omitempty"`
	Ids          []int        `json:"ids,omitempty"`
	DeleteData   bool         `json:"delete-local-data,omitempty"`
	DownloadDir  string       `json:"download-dir,omitempty"`
	MetaInfo     string       `json:"metainfo,omitempty"`
	TorrentAdded TorrentAdded `json:"torrent-added"`
}

type Torrent struct {
	Id            int     `json:"id"`
	Name          string  `json:"name"`
	Status        int     `json:"status"`
	LeftUntilDone int     `json:"leftUntilDone"`
	Eta           int     `json:"eta"`
	UploadRatio   float64 `json:"uploadRatio"`
	RateDownload  int     `json:"rateDownload"`
	RateUpload    int     `json:"rateUpload"`
}

type TorrentAdded struct {
	HashString string `json:"hashString"`
	Id         int    `json:"id"`
	Name       string `json:"name"`
}

func New(url string, username string, password string) TransmissionClient {
	apiclient := client.NewClient(url, username, password)
	tc := TransmissionClient{apiclient: apiclient}
	return tc
}

func (ac *TransmissionClient) GetTorrents() ([]Torrent, error) {
	var getCommand command
	var outputCommand command
	getCommand.Method = "torrent-get"
	getCommand.Arguments.Fields = []string{"id", "name",
		"status", "leftUntilDone", "eta", "uploadRatio",
		"rateDownload", "rateUpload"}
	body, err := json.Marshal(getCommand)
	if err != nil {
		return []Torrent{}, err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return []Torrent{}, err
	}
	err = json.Unmarshal(output, &outputCommand)
	if err != nil {
		return []Torrent{}, err
	}
	return outputCommand.Arguments.Torrents, nil
}

func (ac *TransmissionClient) RemoveTorrent(id int, removeFile bool) (string, error) {
	var removeCommand command
	var outputCommand command

	removeCommand.Method = "torrent-remove"
	removeCommand.Arguments.Ids = []int{id}
	removeCommand.Arguments.DeleteData = removeFile
	body, err := json.Marshal(removeCommand)
	if err != nil {
		return "", err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(output, &outputCommand)
	if err != nil {
		return "", err
	}
	return outputCommand.Result, nil
}

func (ac *TransmissionClient) AddTorrent(file string, downloadDir string) (TorrentAdded, error) {
	var addCommand command
	var outputCommand command

	encodedFile, err := encodeFile(file)
	if err != nil {
	}

	addCommand.Method = "torrent-add"
	addCommand.Arguments.MetaInfo = encodedFile
	addCommand.Arguments.DownloadDir = downloadDir

	body, err := json.Marshal(addCommand)
	if err != nil {
		return TorrentAdded{}, err
	}
	output, err := ac.apiclient.Post(string(body))
	if err != nil {
		return TorrentAdded{}, err
	}
	err = json.Unmarshal(output, &outputCommand)
	if err != nil {
		return TorrentAdded{}, err
	}
	return outputCommand.Arguments.TorrentAdded, nil
}

func encodeFile(file string) (string, error) {
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileData), nil
}
