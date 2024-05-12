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
    uuidString := f.BucketKey.GetGuid().String()

    // Find the index of the last occurrence of "/"
    lastSlashIndex := strings.LastIndex(f.BucketKey.Key, "/") // users/spencer/<uuidString>-<filename>.<filetype> 

    if lastSlashIndex != -1 {
        // Extract the substring from the start to the last "/"
        substring := f.BucketKey.Key[:lastSlashIndex+1]
       
        cloudFrontCaptionsUrl := f.RootUrl + substring + uuidString + "-captions.vtt"

        return cloudFrontCaptionsUrl
    } else {
        cloudFrontCaptionsUrl := f.RootUrl + uuidString + "-captions.vtt" 

        return cloudFrontCaptionsUrl
    } 
}

func (f *FileUrl) GetUrlTranscriptionOutputJson() string {
    uuidString := f.BucketKey.GetGuid().String()

    // Find the index of the last occurrence of "/"
    lastSlashIndex := strings.LastIndex(f.BucketKey.Key, "/") // users/spencer/<uuidString>-<filename>.<filetype> 

    if lastSlashIndex != -1 {
        // Extract the substring from the start to the last "/"
        substring := f.BucketKey.Key[:lastSlashIndex+1]
       
        cloudFrontTranscriptionUrl := f.RootUrl + substring + uuidString + "-transcription-output.json"

        return cloudFrontTranscriptionUrl
    } else {
        cloudFrontTranscriptionUrl := f.RootUrl + uuidString + "-transcription-output.json"

        return cloudFrontTranscriptionUrl
    } 
}
