package kahla

import (
	"fmt"
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

func (s *OssService) HeadImgFile(in *Oss_Download_FromKeyRequest) ([]byte, *http.Response, error) {
	v := url.Values{}
	if in.W != nil {
		v.Add("w", strconv.FormatUint(uint64(*in.W), 10))
	}
	if in.H != nil {
		v.Add("h", strconv.FormatUint(uint64(*in.H), 10))
	}
	resp, err := s.client.Get(s.baseUrl + fmt.Sprintf("/Download/FromKey/%d", in.HeadImgFileKey) + "?" + v.Encode())
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
