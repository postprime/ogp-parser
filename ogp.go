package ogpparser

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"strings"

	"strconv"

	"io/ioutil"

	"github.com/PuerkitoBio/goquery"
	"github.com/julianshen/go-readability"
)

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

func GetPageDataFromHtml(html []byte, data interface{}) error {
	buf := bytes.NewBuffer(html)
	doc, err := goquery.NewDocumentFromReader(buf)

	if err != nil {
		return err
	}

	return GetPageData(doc, data)
}

func GetPageData(doc *goquery.Document, data interface{}) error {
	doc = goquery.CloneDocument(doc)
	return getPageData(doc, data)
}

func GetPageInfoFromResponse(response *http.Response) (*PageInfo, error) {
	info := PageInfo{}
	html, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	err = GetPageDataFromHtml(html, &info)

	if err != nil {
		return nil, err
	}

	r, err := readability.NewDocument(string(html))
	if err != nil {
		return nil, err
	}

	info.Content = r.Text()

	return &info, nil
}

func GetPageInfoFromUrl(urlStr string) (*PageInfo, error) {
	resp, err := http.Get(urlStr)

	if err != nil {
		return nil, err
	}
	return GetPageInfoFromResponse(resp)
}

func getPageData(doc *goquery.Document, data interface{}) error {
	var rv reflect.Value
	var ok bool
	if rv, ok = data.(reflect.Value); !ok {
		rv = reflect.ValueOf(data)
	}

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("should not be non-ptr or nil")
	}

	rt := rv.Type()

	for i := 0; i < rv.Elem().NumField(); i++ {
		fv := rv.Elem().Field(i)
		field := rt.Elem().Field(i)

		switch fv.Type().Kind() {
		case reflect.Ptr:
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			e := getPageData(doc, fv)

			if e != nil {
				return e
			}
		case reflect.Struct:
			e := getPageData(doc, fv.Addr())

			if e != nil {
				return e
			}
		case reflect.Slice:
			if fv.IsNil() {
				fv.Set(reflect.MakeSlice(fv.Type(), 0, 0))
			}

			switch field.Type.Elem().Kind() {
			case reflect.Struct:
				last := reflect.New(field.Type.Elem())
				for {
					data := reflect.New(field.Type.Elem())
					e := getPageData(doc, data.Interface())

					if e != nil {
						return e
					}

					//Ugly solution (I can't remove nodes. Why?)
					if !reflect.DeepEqual(last.Elem().Interface(), data.Elem().Interface()) {
						fv.Set(reflect.Append(fv, data.Elem()))
						last.Elem().Set(data.Elem())

					} else {
						break
					}
				}
			case reflect.Ptr:
				last := reflect.New(field.Type.Elem().Elem())
				for {
					data := reflect.New(field.Type.Elem().Elem())
					e := getPageData(doc, data.Interface())

					if e != nil {
						return e
					}

					//Ugly solution (I can't remove nodes. Why?)
					if !reflect.DeepEqual(last.Elem().Interface(), data.Elem().Interface()) {
						fv.Set(reflect.Append(fv, data))
						last.Elem().Set(data.Elem())

					} else {
						break
					}
				}
			default:
				if tag, ok := field.Tag.Lookup("meta"); ok {
					tags := strings.Split(tag, ",")

					for _, t := range tags {
						var contents []reflect.Value

						processMeta := func(idx int, sel *goquery.Selection) {
							if c, existed := sel.Attr("content"); existed {
								if field.Type.Elem().Kind() == reflect.String {
									contents = append(contents, reflect.ValueOf(c))
								} else {
									i, e := strconv.Atoi(c)

									if e == nil {
										contents = append(contents, reflect.ValueOf(i))
									}
								}

								fv.Set(reflect.Append(fv, contents...))
							}
						}

						doc.Find(fmt.Sprintf("meta[property=\"%s\"]", t)).Each(processMeta)

						doc.Find(fmt.Sprintf("meta[name=\"%s\"]", t)).Each(processMeta)

						fv = reflect.Append(fv, contents...)
					}
				}
			}
		default:
			if tag, ok := field.Tag.Lookup("meta"); ok {

				tags := strings.Split(tag, ",")

				content := ""
				existed := false
				sel := (*goquery.Selection)(nil)
				for _, t := range tags {
					if sel = doc.Find(fmt.Sprintf("meta[property=\"%s\"]", t)).First(); sel.Size() > 0 {
						content, existed = sel.Attr("content")
					}

					if !existed {
						if sel = doc.Find(fmt.Sprintf("meta[name=\"%s\"]", t)).First(); sel.Size() > 0 {
							content, existed = sel.Attr("content")
						}
					}

					if existed {
						if fv.Type().Kind() == reflect.String {
							fv.Set(reflect.ValueOf(content))
						} else if fv.Type().Kind() == reflect.Int {
							if i, e := strconv.Atoi(content); e == nil {
								fv.Set(reflect.ValueOf(i))
							}
						}
						break
					}
				}
			}
		}
	}
	return nil
}
