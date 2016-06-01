package logging

import (
	"bytes"
	"container/list"
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// EmailWriter group logs and send as email
type EmailWriter struct {
	// Email server
	Server string

	// Email sender and account name
	Sender string

	// Email account password
	Password string

	// Email Receivers
	Receivers []string

	// Email subject
	Subject string

	// the logs arrived after first log in "Delay" duration will be grouped with the first log
	Delay time.Duration

	ch chan []byte
}

// Initialize setup the writer
func (writer *EmailWriter) Initialize() {
	// 启动发送守候
	writer.ch = make(chan []byte)
	go func() {
		messages := list.New()
		for {
			messages.Init()
			messages.PushBack(<-writer.ch)
			ticker := time.NewTicker(writer.Delay)
			for delay := true; delay; {
				select {
				case message := <-writer.ch:
					messages.PushBack(message)
				case <-ticker.C:
					ticker.Stop()
					delay = false
				}
			}
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "From: %s\r\n", writer.Sender)
			fmt.Fprintf(&buf, "To: %s\r\n", strings.Join(writer.Receivers, ","))
			fmt.Fprintf(&buf, "Subject: %s\r\n", "=?utf-8?B?"+base64.StdEncoding.EncodeToString([]byte(writer.Subject))+"?=")
			buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
			for elem := messages.Front(); elem != nil; elem = elem.Next() {
				buf.Write(elem.Value.([]byte))
			}
			auth := smtp.PlainAuth(writer.Sender, writer.Sender, writer.Password,
				strings.Split(writer.Server, ":")[0])
			if err := smtp.SendMail(writer.Server, auth, writer.Sender,
				writer.Receivers, buf.Bytes()); err != nil {
				fmt.Fprintf(os.Stderr, "send log email fail: subject=%q, error=%q\n",
					writer.Subject, err.Error())
			} else {
				fmt.Fprintf(os.Stderr, "send log email success: subject=%q\n", writer.Subject)
			}
		}
	}()
}

// Write receives logs
func (writer *EmailWriter) Write(bytes []byte) (int, error) {
	writer.ch <- bytes
	return len(bytes), nil
}
