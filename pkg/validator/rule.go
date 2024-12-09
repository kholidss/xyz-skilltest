package validator

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
	"time"
)

var (
	rgxNumericString = regexp.MustCompile(`^\d+$`)
	rgxHumanName     = regexp.MustCompile(`^[a-zA-Z\'’.,\s]+$`)
)

func validateNumericString(v string) bool {
	if v == "" {
		return true
	}
	return rgxNumericString.MatchString(v)
}

func validateHumanName(v string) bool {
	if v == "" {
		return true
	}
	return rgxHumanName.MatchString(v)
}

func validateDOB(v string) bool {
	var (
		f  = []string{`2006-01-02`, `02-01-2006`}
		tm time.Time
	)

	if len(v) == 0 {
		return true
	}

	for i := 0; i < len(f); i++ {
		t, err := time.Parse(f[i], v)

		if err != nil && i == 0 {
			continue
		}

		if err != nil {
			return false
		}

		tm = t
		break
	}

	if tm.After(time.Now().AddDate(-15, 0, 0)) {
		return false
	}

	return true
}

func ValidateHumanName() validation.StringRule {
	return validation.NewStringRuleWithError(
		validateHumanName,
		validation.NewError("validation_name", "Invalid format. This field only allow these following characters: alphabet, single quote ('), space, comma(,), and period(.)."))
}

func ValidateNIK() validation.StringRule {
	return validation.NewStringRuleWithError(
		validateNumericString,
		validation.NewError("validation_nik", `Invalid format, e.g 11111111112XXXXX, 44128777124XXXXX, 64111331221XXXXX`))
}

func ValidateDOB() validation.StringRule {
	return validation.NewStringRuleWithError(
		validateDOB,
		validation.NewError("validation_dob", "Must be valid date of bird YYYY-MM-DD or DD-MM-YYYY and must be 17 years old"))
}
