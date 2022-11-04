package notify

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMsg(t *testing.T) {
	msg := fmt.Sprintf("\nBotID: %s\nEmergency Level: %d\nMessage: %s", "xxx123", 2, "測試測試")
	err := SendMsg(msg)
	assert.NoError(t, err)
}
