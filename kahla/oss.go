package kahla

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type OssService struct {
	client  *http.Client
	baseUrl string
}

func NewOssService(client *http.Client, baseUrl string) *OssService {
	return &OssService{client: client, baseUrl: baseUrl}
}

func (s *OssService) HeadImgFile(headImgFileKey uint32, w int, h int) ([]byte, *http.Response, error) {
	v := url.Values{}
	v.Set("w", strconv.Itoa(w))
	v.Set("h", strconv.Itoa(h))
	resp, err := s.client.Get("https://oss.aiursoft.com/Download/FromKey/%d" + strconv.FormatUint(uint64(headImgFileKey), 10) + "?" + v.Encode())
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, resp, err
	}
	return bytes, resp, nil
}
