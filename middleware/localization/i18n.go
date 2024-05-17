package localization

import (
	"github.com/gopi-frame/contract/web"
	"github.com/gopi-frame/web/request"
)

type Locale struct {
	LanguageGetter func(request *request.Request) string
}

func (locale *Locale) Handle(request *request.Request, next func(*request.Request) web.Responser) web.Responser {
	var language string
	if locale.LanguageGetter == nil {
		language = request.Header("Accept-Language", "en").String()
	} else {
		language = locale.LanguageGetter(request)
	}
	request.SetLocale(language)
	return next(request)
}
