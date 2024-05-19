package common

import (
    "strings"
    "errors"
)

type FileUrl struct {
    BucketKey BucketKey
    RootUrl string // "https://d271tjczb1hjof.cloudfront.net/" IMPORTANT: ends in a trailing slash
	               //  imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "goserverbucket", *obj.Key)
}

// Constructor
func NewFileUrl(cloudFrontUrl string) (*FileUrl, error) {
    dotNetSlashIndex := strings.Index(cloudFrontUrl, ".net/")

    if dotNetSlashIndex != -1 {
        // Extract the substring from the start to the ".net/" (inclusive)
        rootUrl := cloudFrontUrl[:dotNetSlashIndex+5]

        key := cloudFrontUrl[dotNetSlashIndex+5:]

        bucketKey := BucketKey{ Key: key }

        fileUrl := &FileUrl{ BucketKey: bucketKey, RootUrl: rootUrl }
       
        return fileUrl, nil
    } else {
        return nil, errors.New("invalid url .net not found")
    }
}

// Constructor
func NewFileUrlWithDefaultRootUrl(bucketKey BucketKey, rootUrl ...string) *FileUrl {
    defaultRootUrl := "https://d271tjczb1hjof.cloudfront.net/"

    fileUrl := &FileUrl{
        BucketKey: bucketKey,
        RootUrl:   defaultRootUrl,
    }

    if len(rootUrl) > 0 {
        fileUrl.RootUrl = rootUrl[0]
    }

    return fileUrl
}

func (f *FileUrl) GetUrl() string {
    url := f.RootUrl + f.BucketKey.Key

    return url
}

func (f *FileUrl) GetUrlCaptionsVtt() string {
    cloudFrontCaptionsUrl := f.RootUrl + f.BucketKey.GetKeyForCaptions() // substring + uuidString + "-captions.vtt"

    return cloudFrontCaptionsUrl
}

func (f *FileUrl) GetUrlThumbnail() string {
    cloudFrontThumbnailUrl := f.RootUrl + f.BucketKey.GetKeyForThumbnail() 

    return cloudFrontThumbnailUrl
}

func (f *FileUrl) GetUrlTranscriptionOutputJson() string {
    cloudFrontTranscriptionsUrl := f.RootUrl + f.BucketKey.GetKeyForTranscription() 

    return cloudFrontTranscriptionsUrl
}
