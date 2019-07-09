package kahla

import (
	"crypto/tls"
	"net/http"

	"net/http/cookiejar"
)

type Client struct {
	client       *http.Client
	baseUrl      string
	Auth         *AuthService
	Conversation *ConversationService
	Devices      *DevicesService
	Files        *FilesService
	Friendship   *FriendshipService
	Groups       *GroupsService
	Oss          *OssService
}

func NewClient(baseUrl string, ossUrl string) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	client.Jar, _ = cookiejar.New(nil)
	return &Client{
		client:       client,
		baseUrl:      baseUrl,
		Auth:         NewAuthService(client, baseUrl),
		Conversation: NewConversationService(client, baseUrl),
		Devices:      NewDevicesService(client, baseUrl),
		Files:        NewFilesService(client, baseUrl),
		Friendship:   NewFriendshipService(client, baseUrl),
		Groups:       NewGroupsService(client, baseUrl),
		Oss:          NewOssService(client, ossUrl),
	}
}
