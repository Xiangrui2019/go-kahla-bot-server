package kahla

type Oss_Download_FromKeyRequest struct {
	HeadImgFileKey uint32  `json:"headImgFileKey,omitempty"`
	W              *uint32 `json:"w,omitempty"`
	H              *uint32 `json:"h,omitempty"`
}