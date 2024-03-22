package smtper_test

import (
	"testing"

	"github.com/ariefdarmawan/kmsg"
	"github.com/ariefdarmawan/kmsg/sender/smtper"
)

func TestSend(t *testing.T) {
	opts := smtper.Options{
		Server:   "smtp.outlook.com",
		Port:     587,
		UID:      "youruid@domain.com",
		Password: "yourextremeDifficultPassword",
		TLS:      true,
	}

	sender := smtper.NewSender(opts)
	err := sender.Send(&kmsg.Message{From: "",
		To:      "email1@hotmail.com",
		Cc:      []string{"email2@hotmail.com", "email1@gmail.com"},
		Bcc:     []string{"email2@gmail.com"},
		Title:   "Pengiriman yang dilakukan secara buta",
		Message: "Pesan ini dikirikan kepada anda sebagai bukti bahwa fungsi gomail sudah berjalan dengan lancar",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("pengiriman berhasil")
}
