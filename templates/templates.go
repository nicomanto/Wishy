package templates

import (
	"sync"
	"text/template"
)

type HtmlTemplateType string

const (
	WishListHtmlTemplateType HtmlTemplateType = "wishlist"
)

var once sync.Once
var HtmlTpls map[HtmlTemplateType]*template.Template

func InitHtmlTpls() map[HtmlTemplateType]*template.Template {
	var err error
	once.Do(func() {
		HtmlTpls, err = loadHtmlTpls()
		if err != nil {
			panic(err)
		}
	})
	return HtmlTpls
}

func loadHtmlTpls() (map[HtmlTemplateType]*template.Template, error) {
	tpls := map[HtmlTemplateType]*template.Template{}
	// wishlist
	if t, e := template.New("wishList.tpl").ParseFiles("./templates/wishList.tpl"); e != nil {
		return nil, e
	} else {
		tpls[WishListHtmlTemplateType] = t
	}
	return tpls, nil
}
