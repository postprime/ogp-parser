package ogp

// Image content
type OgImage struct {
	Url    string `meta:"og:image,og:image:url"`
	Width  int    `meta:"og:image:width"`
	Height int    `meta:"og:image:height"`
	Type   string `meta:"og:image:type"`
}

// Video content
type OgVideo struct {
	Url       string `meta:"og:video,og:video:url"`
	SecureUrl string `meta:"og:video:secure_url"`
	Width     int    `meta:"og:video:width"`
	Height    int    `meta:"og:video:height"`
	Type      string `meta:"og:video:type"`
}

// Audio
type OgAudio struct {
	Url       string `meta:"og:audio,og:audio:url"`
	SecureUrl string `meta:"og:audio:secure_url"`
	Type      string `meta:"og:audio:type"`
}

// Page info
type PageInfo struct {
	Title       string `meta:"og:title"`
	Type        string `meta:"og:type"`
	Url         string `meta:"og:url"`
	Site        string `meta:"og:site"`
	SiteName    string `meta:"og:site_name"`
	Description string `meta:"og:description"`
	Locale      string `meta:"og:locale"`
	Images      []*OgImage
	Videos      []*OgVideo
	Audios      []*OgAudio
	Content     string
}
