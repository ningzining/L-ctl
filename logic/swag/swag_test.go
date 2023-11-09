package swag

import (
	"log"
	"testing"
)

func TestSwag_Upload(t *testing.T) {
	if err := NewSwag("./swagger.json", "3567347").Upload(); err != nil {
		log.Println("error", err)
		return
	}
}
