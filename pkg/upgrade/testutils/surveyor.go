package testutils

// AutoAnsweringSurveyor simulates a user always giving the configured answer
type AutoAnsweringSurveyor struct {
	answer bool
}

// AskUserIfReady returns the previously configured answer
func (s AutoAnsweringSurveyor) AskUserIfReady() (bool, error) {
	return s.answer, nil
}

// NewAutoAnsweringSurveyor returns a new AutoAnsweringSurveyor
func NewAutoAnsweringSurveyor(answer bool) AutoAnsweringSurveyor {
	return AutoAnsweringSurveyor{
		answer,
	}
}
