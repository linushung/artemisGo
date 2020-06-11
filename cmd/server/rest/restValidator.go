package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/go-playground/validator/v10"

	log "github.com/sirupsen/logrus"
)

var (
	once     sync.Once
	validate *validator.Validate
)

const (
	ParseErrMsg = "***** [VALIDATION][FAIL] ***** Cannot parse request:: %#v"
	ValidateErrMsg = "***** [VALIDATION][FAIL] ***** Invalid request body:: %v"
)

func init() {
	once.Do(func() {
		// Use a single instance of Validate, it caches struct info
		validate = validator.New()
	})
}
/* Ref:
1. https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
2. https://wang.yuxuan.org/blog/item/2017/02/some-go-memory-notes
*/
func validateReq(body io.ReadCloser, request interface{}) error {

	/* Ref: https://ahmet.im/blog/golang-json-decoder-pitfalls/ */
	if err := json.NewDecoder(body).Decode(request); err != nil {
		log.Errorf(ParseErrMsg, err)
		return fmt.Errorf(ValidateErrMsg, err)
	}

	if errs := validate.Struct(request); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			log.Errorf(ParseErrMsg, err)
			return fmt.Errorf(ValidateErrMsg, err)
		}
	}

	return nil
}
