package cmd

import (
	"fmt"
	"testing"
)

func TestUpdateTemplate(t *testing.T) {
	err := updateTemplate()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	return
}
