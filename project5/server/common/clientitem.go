package common

type ClientItem struct {
	VideoId string `json:"videoId"`
	VideoName string `json:"videoName"`
	VideoUrl string `json:"videoUrl"`
	VideoCaptionsUrl string `json:"videoCaptionsUrl"`
	VideoThumbnailUrl string `json:"videoThumbnailUrl"`
}
