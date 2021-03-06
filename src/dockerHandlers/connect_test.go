/*
author Vyacheslav Kryuchenko
*/
package dockerHandlers

import (
	"testing"
)

func TestCreateClientConnection(t *testing.T) {
	const (
		dockerhost = "tcp://localhost:2376"
		apiversion = "1.32"
	)
	testClient, err := CreateClientConnection(dockerhost, apiversion)
	if err != nil {
		t.Error(err)
		return
	}
	defer testClient.Close()
}
