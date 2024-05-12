package common

type FileUrl struct {
    BucketKey BucketKey
    RootUrl string // "https://d271tjczb1hjof.cloudfront.net/" ends in a trailing slash
	               //  imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", "goserverbucket", *obj.Key)
}

func (f *FileUrl) GetUrl() string {
    url := f.RootUrl + f.BucketKey.Key

    return url
}

