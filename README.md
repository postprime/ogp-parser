### Open Graph Protocol Parser
A Golang library to fetch the OGP contents from URl

#### Parse the basic contents from meta tags

```go
type OgpPageInfo struct {
	Title       string `meta:"og:title"`
	Url         string `meta:"og:url"`
	SiteName    string `meta:"og:site_name"`
	Images      []*OgpImage
}
```

```go
type OgpImage struct {
	Url    string `meta:"og:image,og:image:url"`
	Width  int    `meta:"og:image:width"`
	Height int    `meta:"og:image:height"`
}
```

### How to use
```go
url := "https://www.bloomberg.co.jp/news/articles/2021-01-19/QN4ILZT0AFBS01"
pageInfo, e := GetPageInfoFromUrl(url)
```