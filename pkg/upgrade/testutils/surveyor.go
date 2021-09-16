package testutils

import (
	"fmt"
)

// AutoAnsweringSurveyor simulates a user always giving the configured answer
type AutoAnsweringSurveyor struct {
	currentAnswer int
	answers       []bool
}

// PromptUser returns the previously configured answer
//goland:noinspection GoUnusedParameter
func (s *AutoAnsweringSurveyor) PromptUser(message string) (bool, error) {
	if len(s.answers) == 0 {
		return true, nil
	}

	if s.currentAnswer >= len(s.answers) {
		return false, fmt.Errorf("no more answers configured. You need to initialize with none or all answers")
	}

	answer := s.answers[s.currentAnswer]
	s.currentAnswer++

	return answer, nil
}

// NewAutoAnsweringSurveyor returns a new AutoAnsweringSurveyor
func NewAutoAnsweringSurveyor(answers []bool) *AutoAnsweringSurveyor {
	return &AutoAnsweringSurveyor{
		answers: answers,
	}
}
