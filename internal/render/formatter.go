package render

import (
	"fmt"

	"github.com/spetix/bar-out-adapters/pkg/barout/models"
)

type dataFormatterImpl struct {
	models.Formatter
	format string
}

func NewDataFormatterImpl() models.Formatter {
	return &dataFormatterImpl{}
}

func (f *dataFormatterImpl) Render(s ...string) string {
	if f.format != "" {
		if len(s) < 2 {
			return ""
		}
		return fmt.Sprintf(f.format, s[0], s[1])
	}
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

func (f *dataFormatterImpl) Validate(format string) bool {
	if format == "" {
		return false
	}
	placeholderCount := 0
	for i := 0; i+1 < len(format); i++ {
		if format[i] == '%' && format[i+1] == 's' {
			placeholderCount++
			i++
		}
	}
	return placeholderCount == 2
}
