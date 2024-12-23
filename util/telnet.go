package util

import (
	"bytes"
	"crypto/cipher"
	"errors"
	"strings"

	"github.com/reiver/go-telnet"
)

// convert a command to bytes, and send to Telnet connection followed by '\r\n'
func WriteTelnet(conn *telnet.Conn, command string) {
	var commandBuffer []byte
	for _, char := range command {
		commandBuffer = append(commandBuffer, byte(char))
	}

	var crlfBuffer [2]byte = [2]byte{'\r', '\n'}
	crlf := crlfBuffer[:]

	//fmt.Println(commandBuffer)

	conn.Write(commandBuffer)
	conn.Write(crlf)
}

func ExecTelnet(conn *telnet.Conn, command string) string {
	WriteTelnet(conn, command)
	return readTelnet(conn, "/ #")
}

func PerformTelnetLogin(conn *telnet.Conn, user string, pass string) error {
	readTelnet(conn, "Login")
	WriteTelnet(conn, user)
	readTelnet(conn, "Password")
	WriteTelnet(conn, pass)
	telnetReply := readTelnet(conn, "/ #")
	if !strings.Contains(telnetReply, "/ #") {
		if strings.Contains(telnetReply, "Access denied") {
			return errors.New("access denied")
		}
		return errors.New("telnet not ready")
	}
	return nil
}

// Thin function reads from Telnet session. "expect" is a string I use as signal to stop reading
func readTelnet(conn *telnet.Conn, expect string) (out string) {
	var buffer [1]byte
	recvData := buffer[:]
	var n int
	var err error

	for {
		n, err = conn.Read(recvData)
		//fmt.Println("Bytes: ", n, "Data: ", recvData, string(recvData))
		if n <= 0 || err != nil || strings.Contains(out, expect) || strings.Contains(out, "Access denied") {
			break
		} else {
			out += string(recvData)
		}
	}
	return out
}

func Pad(dataToPad []byte, blockSize int) []byte {
	paddingLen := blockSize - len(dataToPad)%blockSize
	padding := bytes.Repeat([]byte{0}, paddingLen)
	return append(dataToPad, padding...)
}

func Unpad(paddedData []byte, blockSize int) []byte {
	trimIndex := len(paddedData) - blockSize
	trimmedData := paddedData[:trimIndex]

	for i := len(trimmedData) - 1; i >= trimIndex; i-- {
		if trimmedData[i] != 0 {
			break
		}
		trimIndex--
	}

	return paddedData[:trimIndex]
}

func EncryptAesEcb(data []byte, cy cipher.Block) []byte {
	encrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cy.Encrypt(encrypted[bs:be], data[bs:be])
	}

	return encrypted
}

func DecryptAesEcb(data []byte, cy cipher.Block) []byte {
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cy.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}
