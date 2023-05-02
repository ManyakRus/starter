package email_imap

import (
	"testing"
)

func TestConnect(t *testing.T) {
	Connect()
	defer CloseConnection()
}

func TestSelectInbox(t *testing.T) {
	Connect()
	defer CloseConnection()

	SelectInbox()
}

func TestReplaceMessage(t *testing.T) {
	//Connect()
	//defer CloseConnection()
	//
	//SelectInbox()
	//ReplaceMessage(msg, FOLDER_NAME_INBOX)
}

func TestReadMessage(t *testing.T) {
	Connect()
	defer CloseConnection()

	SelectInbox()
	//SelectFolder("Архив")
	size, _ := Stat()
	if size <= 0 {
		t.Log("size: 0")
		return
	}

	msg, err := ReadMessage(size)
	if err != nil {
		t.Error("emal_imap_test.TestReadMessage() error: ", err)
	} else {
		t.Log("message: ", msg)
	}

}
