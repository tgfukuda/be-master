package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgfukuda/be-master/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	assert.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "Test email"
	content := `
	<h2>Hello World</h2>
	<p>This is a test email from <a href="https://github.com/tgfukuda/be-master">My BE Master</a></p>
	`
	to := []string{config.EmailSenderAddress}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	assert.NoError(t, err)
}
