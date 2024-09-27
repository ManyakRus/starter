package email_imap

import (
	//"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ManyakRus/starter/contextmain"
	"github.com/ManyakRus/starter/email"
	"github.com/ManyakRus/starter/log"
	"github.com/ManyakRus/starter/micro"
	"github.com/ManyakRus/starter/stopapp"
	"github.com/emersion/go-imap"
	imapModule "github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	mail "github.com/emersion/go-message/mail"
	"github.com/joho/godotenv"
	simplemail "github.com/xhit/go-simple-mail/v2"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

var Conn *imapModule.Client
var MailInbox *imap.MailboxStatus // папка inbox

// Settings хранит все нужные переменные окружения
var Settings SettingsINI

var FOLDER_NAME_INBOX = `INBOX`
var ErrEmptyInbox = fmt.Errorf("empty inbox")

// SettingsINI - структура для хранения всех нужных переменных окружения
type SettingsINI struct {
	EMAIL_IMAP_SERVER  string
	EMAIL_IMAP_PORT    string
	EMAIL_LOGIN        string
	EMAIL_PASSWORD     string
	EMAIL_SEND_TO_TEST string
	//EMAIL_SUBJECT        string
	EMAIL_AUTHENTICATION string
	EMAIL_ENCRYPTION     string
}

type Attachment struct {
	Filename string
	Data     []byte
}

// Connect - подключение клиента Email
func Connect() {
	err := Connect_err()
	LogInfo_Connected(err)

}

// LogInfo_Connected - выводит сообщение в Лог, или паника при ошибке
func LogInfo_Connected(err error) {
	if err != nil {
		log.Panicln("Connect() error: ", err)
	} else {
		log.Info("Email connected: ", Settings.EMAIL_LOGIN)
	}

}

// Connect_err - Однократно Устанавливает соединение по требованию
func Connect_err() error {
	var err error

	if Settings.EMAIL_LOGIN == "" {
		LoadEnv()
	}

	strFrom := Settings.EMAIL_LOGIN
	strPass := Settings.EMAIL_PASSWORD
	strHost := Settings.EMAIL_IMAP_SERVER

	//log.Debugf("Connecting to server %s", strHost)

	//подключение к серверу
	fn := func() error {
		Conn, err = imapModule.Dial(strHost + ":" + Settings.EMAIL_IMAP_PORT)
		if err != nil {
			err1 := fmt.Errorf("Send() error: %w", err)
			return err1
		}
		return err
	}

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()
	err = micro.GoGo(ctx, fn)

	//Conn, err := imapModule.Dial(strHost + ":" + strconv.Itoa(Settings.EMAIL_SMTP_PORT))
	if err != nil {
		return err
	}

	// Login
	fn = func() error {
		err = Conn.Login(strFrom, strPass)
		if err != nil {
			err1 := fmt.Errorf("Send() error: %w", err)
			return err1
		}
		return err
	}

	ctx, cancel2 := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel2()
	err = micro.GoGo(ctx, fn)

	//err = Conn.Login(strFrom, strPass)
	if err != nil {
		return err
	}
	log.Infof("Logged in %s", strFrom)

	MailInbox = SelectInbox()

	return nil
}

// SelectInbox - возвращает емайл папку Inbox
func SelectInbox() *imap.MailboxStatus {
	var MailInbox *imap.MailboxStatus
	MailInbox = SelectFolder(FOLDER_NAME_INBOX)

	return MailInbox
}

// SelectInbox - возвращает емайл папку Inbox
func SelectFolder(FolderName string) *imap.MailboxStatus {
	var MailInbox *imap.MailboxStatus

	if Conn == nil {
		log.Errorf("mailconn.SelectFolder() error conn = nil")
		return MailInbox
	}

	// Select INBOX
	var err error

	fn := func() error {
		MailInbox, err = Conn.Select(FolderName, false)
		if err != nil {
			err1 := fmt.Errorf("Send() error: %w", err)
			return err1
		}
		return err
	}

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()
	err = micro.GoGo(ctx, fn)

	//MailInbox, err = Conn.Select(FolderName, false)
	if err != nil {
		log.Error("mailconn.SelectFolder() Select(", FolderName, ") error: %s", err)
		panic(err)
	}
	if MailInbox == nil {
		log.Error("mailconn.SelectFolder() Select(", FolderName, ") MailInbox=nil, error: %s", err)
		panic(err)
	}

	return MailInbox
}

// ReplaceMessage -- перемещает сообщение в другую папку
func ReplaceMessage(msg *imap.Message, FolderName string) error {
	if msg == nil {
		return fmt.Errorf("ReplaceMessage(): msg==nil")
	}

	SeqSet := FindSeqSet(msg.Uid)

	var err error
	//ctxMain := contextmain.GetContext()
	//Ctx1, CancelFunc1 := context.WithTimeout(ctxMain, time.Second*30)
	//defer CancelFunc1()

	err = Conn.UidMove(SeqSet, FolderName)

	//select {
	//case <-Ctx1.Done():
	//	Text1 := "mailconn.MsgNext() Fetch() error: TimeOut"
	//	err = errors.New(Text1)
	//	return err
	//case err = <-Chan1:
	//}
	////err := mc.conn.UidMove(SeqSet, FolderName)

	return err
}

// FindSeqSet -- находит SeqSet по номеру сообщения
func FindSeqSet(Id uint32) *imap.SeqSet {

	SeqSet := new(imap.SeqSet)
	SeqSet.AddNum(Id)

	return SeqSet
}

// FindFolderName - возвращает имя папки imap
func FindFolderName(MainFolderName, SubFolderName string) string {
	Otvet := MainFolderName

	if SubFolderName != "" {
		Otvet = MainFolderName + `/` + SubFolderName
	}

	return Otvet
}

// ForwardMessage -- перенаправляет емайл
func ForwardMessage(msg *imap.Message, email_send_to string) error {
	//BodyText := msg.Body
	//BodyText = "Обращение поступило от: " + msg.Envelope.From.MailboxName + "\n\r<BR>" +
	//	"----------------------------------------------------------\n\r<BR><BR>" +
	//	BodyText

	var err error
	//ctxMain := contextmain.GetContext()
	//Ctx1, CancelFunc1 := context.WithTimeout(ctxMain, time.Second*30)
	//defer CancelFunc1()

	Body, Attachments, err := ReadBody(msg)

	MassAttachments := make([]simplemail.File, 0)
	for _, v := range Attachments {
		Attachment1 := simplemail.File{}
		Attachment1.Name = v.Filename
		Attachment1.Data = v.Data
		MassAttachments = append(MassAttachments, Attachment1)
	}

	err = email.SendEmail(email_send_to, Body, msg.Envelope.Subject, MassAttachments)

	//select {
	//case <-Ctx1.Done():
	//	Text1 := "mailconn.MsgNext() Fetch() error: TimeOut"
	//	err = errors.New(Text1)
	//	return err
	//case err = <-Chan1:
	//}

	return err
}

// WaitStop - ожидает отмену глобального контекста
func WaitStop() {

	select {
	case <-contextmain.GetContext().Done():
		log.Warn("Context app is canceled.")
	}

	//
	stopapp.WaitTotalMessagesSendingNow("email")

	//
	CloseConnection()
	stopapp.GetWaitGroup_Main().Done()
}

// Start - необходимые процедуры для подключения к серверу email imap
func Start() {
	var err error

	ctx := contextmain.GetContext()
	WaitGroup := stopapp.GetWaitGroup_Main()
	err = Start_ctx(&ctx, WaitGroup)
	LogInfo_Connected(err)

}

// Start_ctx - необходимые процедуры для подключения к серверу email imap
// Свой контекст и WaitGroup нужны для остановки работы сервиса Graceful shutdown
// Для тех кто пользуется этим репозиторием для старта и останова сервиса можно просто Start()
func Start_ctx(ctx *context.Context, WaitGroup *sync.WaitGroup) error {
	var err error

	//запомним к себе контекст и WaitGroup
	contextmain.Ctx = ctx
	stopapp.SetWaitGroup_Main(WaitGroup)

	//
	LoadEnv()
	err = Connect_err()
	if err != nil {
		return err
	}

	stopapp.GetWaitGroup_Main().Add(1)
	go WaitStop()

	return err
}

// LoadEnv - загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv() {

	dir := micro.ProgramDir()
	filename := dir + ".env"
	LoadEnv_FromFile(filename)
}

// LoadEnv_FromFile загружает переменные окружения в структуру из файла или из переменных окружения
func LoadEnv_FromFile(filename string) {
	//var err error
	//err := godotenv.Load(Filename_Settings)
	//if err != nil {
	//	log.Fatal("Error loading " + Filename_Settings + " file, error: " + err.Error())
	//}

	err := godotenv.Load(filename)
	if err != nil {
		log.Debug("Error parse .env file error: " + err.Error())
	} else {
		log.Info("load .env from file: ", filename)
	}

	Settings = SettingsINI{}
	Settings.EMAIL_IMAP_SERVER = os.Getenv("EMAIL_IMAP_SERVER")
	Settings.EMAIL_IMAP_PORT = os.Getenv("EMAIL_IMAP_PORT")
	Settings.EMAIL_LOGIN = os.Getenv("EMAIL_LOGIN")
	Settings.EMAIL_PASSWORD = os.Getenv("EMAIL_PASSWORD")
	Settings.EMAIL_SEND_TO_TEST = os.Getenv("EMAIL_SEND_TO_TEST")
	//Settings.EMAIL_SUBJECT = os.Getenv("EMAIL_SUBJECT")
	Settings.EMAIL_AUTHENTICATION = os.Getenv("EMAIL_AUTHENTICATION")
	Settings.EMAIL_ENCRYPTION = os.Getenv("EMAIL_ENCRYPTION")

	if Settings.EMAIL_IMAP_SERVER == "" {
		log.Panicln("Need fill EMAIL_SMTP_SERVER ! in file ", filename)
	}

	if Settings.EMAIL_IMAP_PORT == "" {
		log.Panicln("Need fill EMAIL_SMTP_PORT ! in file ", filename)
	}

	if Settings.EMAIL_LOGIN == "" {
		log.Panicln("Need fill EMAIL_LOGIN ! in file ", filename)
	}

	if Settings.EMAIL_PASSWORD == "" {
		log.Panicln("Need fill EMAIL_PASSWORD ! in file ", filename)
	}

	if Settings.EMAIL_SEND_TO_TEST == "" && micro.IsTestApp() == true {
		log.Info("Need fill EMAIL_SEND_TO_TEST ! in file ", filename)
	}

	//if Settings.EMAIL_SUBJECT == "" {
	//	log.Panicln("Need fill EMAIL_SUBJECT ! in file ", filename)
	//}

	if Settings.EMAIL_AUTHENTICATION == "" {
		log.Warn("Need fill EMAIL_AUTHENTICATION ! in file ", filename)
	}

	if Settings.EMAIL_ENCRYPTION == "" {
		log.Warn("Need fill EMAIL_ENCRYPTION ! in file ", filename)
	}

}

// CloseConnection - остановка работы клиента Email
func CloseConnection() {
	err := CloseConnection_err()
	if err != nil {
		log.Panic("Email imap CloseConnection() error: ", err)
	} else {
		log.Info("Email imap connection closed")
	}

}

// CloseConnection_err -- закрывает соединение с почтовым сервером
func CloseConnection_err() error {
	var err error

	fn := func() error {
		err := Conn.Logout()
		if err != nil {
			err1 := fmt.Errorf("Send() error: %w", err)
			return err1
		}
		return err
	}

	ctxMain := contextmain.GetContext()
	ctx, cancel := context.WithTimeout(ctxMain, 60*time.Second)
	defer cancel()
	err = micro.GoGo(ctx, fn)

	//err := Conn.Logout()
	if err != nil {
		log.Printf("MailConn.Logout(): err=\t%v", err)
	}

	Conn = nil

	return err
}

// Stat - возвращает количество сообщений и ИД первого непрочитанного
func Stat() (count, id int) {
	// Select INBOX
	MailInbox := SelectInbox()
	count = int(MailInbox.Messages)
	id = int(MailInbox.UnseenSeqNum)

	return count, id
}

// Безопасно читает тело сообщения
func ReadBody(Msg *imap.Message) (BodyText string, Attachments []Attachment, err error) {

	Attachments = make([]Attachment, 0)

	chErr := make(chan string, 2)
	if Msg == nil {
		//chErr <- "MsgEmail.readBody().fnRead(): сообщение не присвоена в структуре"
		sError := "MsgEmail.readBody(): Msg == nil"
		err = fmt.Errorf(sError)
		chErr <- sError
		return
	}
	var section imap.BodySectionName
	l := Msg.GetBody(&section)
	if l == nil {
		sError := "MsgEmail.readBody(): GetBody() = nil"
		err = fmt.Errorf(sError)
		chErr <- sError
		return
		//log.Fatal("Server didn't returned message body")
	}

	//sBody := l.(*bytes.Buffer).String()

	// Create a new mail reader
	mr, err := mail.CreateReader(l)
	if err != nil {
		sError := fmt.Sprintf("MsgEmail.readBody(): CreateReader() error: %v", err)
		err = fmt.Errorf(sError)
		chErr <- sError
		return
		//log.Fatal(err)
	}

	//buf := make([]byte, 0)
	//n, err := mr.Read(buf)

	// Process each message's part
	//sBody := ""
	BodyText = ""
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if message.IsUnknownCharset(err) {
			continue
		}
		if err != nil {
			sError := fmt.Sprintf("MsgEmail.readBody(): NextPart() error: %v", err)
			err = fmt.Errorf(sError)
			chErr <- sError
			if strings.Contains(err.Error(), "EOF") {
				break
			}
		}

		if p == nil {
			continue
		}

		switch p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, err := ioutil.ReadAll(p.Body)
			if err != nil {
				sError := fmt.Sprintf("MsgEmail.readBody(): ReadAll() error: %v", err)
				err = fmt.Errorf(sError)
				chErr <- sError
			}
			//sBody = sBody + string(b)
			BodyText = string(b)
			break

		case *mail.AttachmentHeader:
			h := p.Header.(*mail.AttachmentHeader)
			// This is an attachment
			filename, _ := h.Filename()
			//log.Println("Got attachment: %v", filename)
			var Data []byte
			Data, err := ioutil.ReadAll(p.Body)
			if err != nil {
				sError := fmt.Sprintf("MsgEmail.readBody(): ReadAll() error: %v", err)
				err = fmt.Errorf(sError)
				chErr <- sError
			}

			if filename != "" {
				Attachment1 := Attachment{}
				Attachment1.Filename = filename
				Attachment1.Data = Data
				Attachments = append(Attachments, Attachment1)
			}

		}
	}

	//sf.StrBody_ = sBody

	//body_, err := io.ReadAll(sf.Msg.GetBody(&section))
	//if err != nil {
	//	chErr <- fmt.Sprintf("MsgEmail.readBody().fnRead(): при чтении тела письма, err=\t%v", err)
	//}
	//sf.StrBody_ = html.UnescapeString(string(body_))
	//if err = sf.getEmail(); err != nil {
	//	chErr <- fmt.Sprintf("MsgEmail.readBody().fnRead(): при получении адресата письма, err=\t%v", err)
	//}

	return
}

// Безопасно получает тему сообщения
func ReadHeader(msg *imap.Message) mail.Header {
	var Header mail.Header

	// Get the whole message body
	var section imap.BodySectionName
	section = imap.BodySectionName{}
	//items := []imap.FetchItem{section.FetchItem()}

	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	Header = mr.Header

	return Header
}

// ReadMessage -- возвращает пиьмо с сервера, с номером по порядку=id
func ReadMessage(id int) (*imap.Message, error) {
	var err error
	Otvet := &imap.Message{}

	messages := make(chan *imap.Message, 1)
	nextNum := uint32(id)
	seqset := new(imap.SeqSet)
	seqset.AddNum(nextNum)
	var section imap.BodySectionName
	items := []imap.FetchItem{imap.FetchBody, imap.FetchUid, imap.FetchEnvelope, imap.FetchBodyStructure, section.FetchItem()}
	err = Conn.Fetch(seqset, items, messages)
	if err != nil {
		return Otvet, err
	}

	TimeOutSeconds := 60
	duration := time.Duration(TimeOutSeconds) * time.Second
	Ctx1, CancelFunc1 := context.WithTimeout(contextmain.GetContext(), duration)
	defer CancelFunc1()
	select {
	case <-Ctx1.Done():
		Text1 := fmt.Sprint("mailconn.MsgNext() Fetch() error: TimeOut ", TimeOutSeconds, " seconds")
		err = errors.New(Text1)
		return Otvet, err
	case Otvet = <-messages:
	}

	if Otvet == nil {
		text1 := fmt.Sprint("mailconn.MsgNext() No email with id: ", nextNum)
		err := errors.New(text1)
		//chRes_ <- text1
		return Otvet, err
	}

	return Otvet, err
}
