package templates

import (
	"html/template"
	"sync"
)

type HtmlTemplateType string

const (
	WishListHtmlTemplateType   HtmlTemplateType = "wishlist"
	ErrorPageHtmlTemplateType  HtmlTemplateType = "errorpage"
	ManagementHtmlTemplateType HtmlTemplateType = "management"
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
	if t, e := template.New("wishList.tpl").ParseFiles("./templates/html/wishList.tpl"); e != nil {
		return nil, e
	} else {
		tpls[WishListHtmlTemplateType] = t
	}
	// error page
	if t, e := template.New("errorPage.tpl").ParseFiles("./templates/html/errorPage.tpl"); e != nil {
		return nil, e
	} else {
		tpls[ErrorPageHtmlTemplateType] = t
	}
	if t, e := template.New("management.tpl").ParseFiles("./templates/html/management.tpl"); e != nil {
		return nil, e
	} else {
		tpls[ManagementHtmlTemplateType] = t
	}
	return tpls, nil
}
