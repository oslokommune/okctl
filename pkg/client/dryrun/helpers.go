package dryrun

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

const toBeGenerated = "<generated>"

func formatCreate(resource string) string {
	return fmt.Sprintf("%s Create %s \n",
		aurora.Green("+"),
		resource,
	)
}

func formatDelete(resource string) string {
	return fmt.Sprintf("%s Delete %s \n",
		aurora.Red("-"),
		resource,
	)
}
