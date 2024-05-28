package tools

import (
	"fmt"
	"testing"
)

func TestDecrypt(t *testing.T) {
	str, err := Decrypt("EMjG405pzOFPmn078/B/Gg==", []byte("ef9bd2441fd7ebdc6f91ab3a89b8f70b"), []byte("vo6hyrk71ueetvkj"))
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}
