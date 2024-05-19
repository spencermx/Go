package common

type ClientItem struct {
	CloudFrontUrl string `json:"url"`
	FileName string `json:"alt"`
	VideoCaptionsUrl string `json:"videocaptionsurl"`
	VideoThumbnailUrl string `json:"videothumbnailurl"`
}
