package prompt

import (
	"fmt"
	"os"
)

type Progress struct {
	totalSteps int
	currentStep int
}

func NewProgress(totalSteps int) *Progress {
	return &Progress{
		totalSteps: totalSteps,
		currentStep: 0,
	}
}

func (p *Progress) Step(message string) {
	p.currentStep++
	fmt.Fprintf(os.Stderr, "[%d/%d] %s... ", p.currentStep, p.totalSteps, message)
}

func (p *Progress) Done() {
	fmt.Fprintln(os.Stderr, "done")
}

func (p *Progress) StepDone(message string) {
	p.Step(message)
	p.Done()
}

func (p *Progress) Fail(err error) {
	fmt.Fprintf(os.Stderr, "failed: %v\n", err)
}
