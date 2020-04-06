package gcorecloud

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/instance/v1/types"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate          *validator.Validate
	Trans             ut.Translator
	interfaceTypes    = types.InterfaceType("").StringList()
	floatingIPSources = types.FloatingIPSource("").StringList()
)

func init() { // nolint
	Validate = validator.New()
	translator := en.New()
	uni := ut.New(translator, translator)
	var found bool
	Trans, found = uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	if err := translations.RegisterDefaultTranslations(Validate, Trans); err != nil {
		log.Fatal(err)
	}

	err := Validate.RegisterTranslation("required", Trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	FailOnErrorF(err, "Cannot register translation for tag: %s", "required")

	err = Validate.RegisterTranslation("url", Trans, func(ut ut.Translator) error {
		return ut.Add("url", "{0} must be a valid url", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url", fe.Field())
		return t
	})
	FailOnErrorF(err, "Cannot register translation for tag: %s", "url")

	err = Validate.RegisterTranslation("rfe", Trans, func(ut ut.Translator) error {
		return ut.Add("rfe", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("rfe", fe.Field())
		return t
	})
	FailOnErrorF(err, "Cannot register translation for tag: %s", "rfe")

	err = Validate.RegisterTranslation("sfe", Trans, func(ut ut.Translator) error {
		return ut.Add("sfe", "{0} is not a required field. Should be empty", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("sfe", fe.Field())
		return t
	})
	FailOnErrorF(err, "Cannot register translation for tag: %s", "sfe")

	err = Validate.RegisterTranslation("floating-ip-source", Trans, func(ut ut.Translator) error {
		return ut.Add(
			"floating-ip-source",
			fmt.Sprintf("{0} should be one of: %v", floatingIPSources),
			true,
		)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("floating-ip-source", fe.Field())
		return t
	})
	FailOnErrorF(err, "Cannot register translation for tag: %s", "floating-api-source")

	err = Validate.RegisterTranslation("interface-type", Trans, func(ut ut.Translator) error {
		return ut.Add(
			"interface-type",
			fmt.Sprintf("{0} should be one of: %v", interfaceTypes),
			true,
		)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("interface-type", fe.Field())
		return t
	})
	FailOnErrorF(err, "Cannot register translation for tag: %s", "interface-type")

	err = Validate.RegisterValidation("floating-ip-source", func(fl validator.FieldLevel) bool {
		return ContainsString(floatingIPSources, fl.Field().String())
	})
	FailOnErrorF(err, "Cannot register custom validator tag: %v", "floating-api-source")

	err = Validate.RegisterValidation("interface-type", func(fl validator.FieldLevel) bool {
		return ContainsString(interfaceTypes, fl.Field().String())
	})
	FailOnErrorF(err, "Cannot register custom validator tag: %v", "interface-type")

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("rfe"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("sfe"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("regex"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("allowed_without"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err = Validate.RegisterValidation(`rfe`, rfe)
	FailOnErrorF(err, "Cannot register custom validation tag: %s", "rfe")
	err = Validate.RegisterValidation(`sfe`, sfe)
	FailOnErrorF(err, "Cannot register custom validation tag: %s", "sfe")
	err = Validate.RegisterValidation(`allowed_without`, allowedWithout)
	FailOnErrorF(err, "Cannot register custom validation tag: %s", "allowed_without")
	err = Validate.RegisterValidation(`regex`, regex)
	FailOnErrorF(err, "Cannot register custom validation tag: %s", "regex")

}

func allowedWithout(fl validator.FieldLevel) bool {
	param := strings.TrimSpace(fl.Param())
	requiredField := requireCheckFieldKind(fl, param, true)
	return requiredField && hasValue(fl)
}

func rfe(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramValue := param[1]

	if paramField == `` {
		return true
	}

	paramValueParts := strings.Split(paramValue, ";")

	// param field reflect.Value.
	var paramFieldValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
	} else {
		paramFieldValue = fl.Parent().FieldByName(paramField)
	}

	found := false

	for _, value := range paramValueParts {
		if isEq(paramFieldValue, value) {
			found = true
		}
	}

	if !found {
		return true
	}

	return hasValue(fl)
}

func sfe(fl validator.FieldLevel) bool {
	param := strings.Split(fl.Param(), `:`)
	paramField := param[0]
	paramValue := param[1]

	if paramField == `` {
		return true
	}

	paramValueParts := strings.Split(paramValue, ";")

	// param field reflect.Value.
	var paramFieldValue reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
	} else {
		paramFieldValue = fl.Parent().FieldByName(paramField)
	}

	found := false

	for _, value := range paramValueParts {
		if isEq(paramFieldValue, value) {
			found = true
		}
	}

	return !found || !hasValue(fl)
}

// requireCheckField is a func for check field kind
func requireCheckFieldKind(fl validator.FieldLevel, param string, defaultNotFoundValue bool) bool {
	field := fl.Field()
	kind := field.Kind()
	var nullable, found bool
	if len(param) > 0 {
		field, kind, nullable, found = fl.GetStructFieldOKAdvanced2(fl.Parent(), param)
		if !found {
			return defaultNotFoundValue
		}
	}
	switch kind {
	case reflect.Invalid:
		return defaultNotFoundValue
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return field.IsNil()
	default:
		if nullable && field.Interface() != nil {
			return false
		}
		return field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface()
	}
}

func hasValue(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		if field.Kind() == reflect.Ptr && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

// regex is the validation function for validating if the current field's value matches to regex.
func regex(fl validator.FieldLevel) bool {

	field := fl.Field()
	rgx := fl.Param()

	validRegex := regexp.MustCompile(rgx)

	if field.Kind() == reflect.String {
		value := field.String()
		return validRegex.MatchString(value)
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func isEq(field reflect.Value, param string) bool {

	switch field.Kind() {

	case reflect.String:
		return field.String() == param

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() == p

	case reflect.Bool:
		p := asBool(param)

		return field.Bool() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asBool(param string) bool {

	i, err := strconv.ParseBool(param)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func asFloat(param string) float64 {

	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func asInt(param string) int64 {

	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func TranslateValidationError(err error) error {
	if err == nil {
		return nil
	}
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	errs := make([]string, len(validationErrors))
	for i, e := range validationErrors {
		errs[i] = fmt.Sprintf("%s: %s", e.Field(), e.Translate(Trans))
	}
	errorText := strings.Join(errs, ", ")
	return fmt.Errorf(errorText)
}
