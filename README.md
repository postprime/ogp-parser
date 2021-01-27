### Open Graph Protocol Parser
A Golang library to fetch the OGP contents from URl

#### Parse the basic contents from meta tags

```go
type OgpPageInfo struct {
	Title       string `meta:"og:title"`
	Type        string `meta:"og:type"`
	Url         string `meta:"og:url"`
	SiteName    string `meta:"og:site_name"`
	Description string `meta:"og:description"`
	Locale      string `meta:"og:locale"`
	Images      []*OgpImage
	Content     string
}
```

```go
type OgpImage struct {
	Url    string `meta:"og:image,og:image:url"`
	Width  int    `meta:"og:image:width"`
	Height int    `meta:"og:image:height"`
	Type   string `meta:"og:image:type"`
}
```

### How to use
```go
url := "https://www.bloomberg.co.jp/news/articles/2021-01-19/QN4ILZT0AFBS01"
pageInfo, e := GetPageInfoFromUrl(url)
```