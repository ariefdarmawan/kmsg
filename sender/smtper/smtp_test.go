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
	err := sender.Send(&kmsg.Message{From: "internal.app@kanosolution.com",
		To:      "ariefda@hotmail.com",
		Cc:      []string{"adarmawan.2006@gmail.com", "arief@kanosolution.com"},
		Bcc:     []string{"arief@ciptaprimanugraha.com"},
		Title:   "Pengiriman yang dilakukan secara buta",
		Message: "Pesan ini dikirikan kepada anda sebagai bukti bahwa fungsi gomail sudah berjalan dengan lancar",
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("pengiriman berhasil")
}
