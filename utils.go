package main

import (
	"bytes"
	"crypto/cipher"
	"errors"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/reiver/go-telnet"
)

func normalizeString(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

func removeLastNChars(s string, lengthNChars int) string {
	return s[:len(s)-lengthNChars]
}

func parseDuration(timeString string) int64 {
	durationSplit := strings.Fields(timeString)
	daysConv, _ := strconv.ParseInt(durationSplit[0], 10, 64)
	hoursConv, _ := strconv.ParseInt(durationSplit[2], 10, 64)
	minsConv, _ := strconv.ParseInt(durationSplit[4], 10, 64)
	secsConv, _ := strconv.ParseInt(durationSplit[6], 10, 64)
	return daysConv*86400 + hoursConv*3600 + minsConv*60 + secsConv
}

func randInt(min int, max int) int {
	return rand.Intn(max-min) + min
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

// convert a command to bytes, and send to Telnet connection followed by '\r\n'
func writeTelnet(conn *telnet.Conn, command string) {
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

func execTelnet(conn *telnet.Conn, command string) string {
	writeTelnet(conn, command)
	return readTelnet(conn, "/ #")
}

func performTelnetLogin(conn *telnet.Conn, user string, pass string) error {
	readTelnet(conn, "Login")
	writeTelnet(conn, user)
	readTelnet(conn, "Password")
	writeTelnet(conn, pass)
	telnetReply := readTelnet(conn, "/ #")
	if !strings.Contains(telnetReply, "/ #") {
		if strings.Contains(telnetReply, "Access denied") {
			return errors.New("access denied")
		}
		return errors.New("telnet not ready")
	}
	return nil
}

func pad(dataToPad []byte, blockSize int) []byte {
	paddingLen := blockSize - len(dataToPad)%blockSize
	padding := bytes.Repeat([]byte{0}, paddingLen)
	return append(dataToPad, padding...)
}

func unpad(paddedData []byte, blockSize int) []byte {
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

func encryptAesEcb(data []byte, cy cipher.Block) []byte {
	encrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cy.Encrypt(encrypted[bs:be], data[bs:be])
	}

	return encrypted
}

func decryptAesEcb(data []byte, cy cipher.Block) []byte {
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cy.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func parseHtmlPage(elemList *[]string, xpath string) error {
	htmlNode := htmlquery.FindOne(cachedPage.GetPage(), xpath)
	if htmlNode != nil {
		*elemList = append(*elemList, normalizeString(htmlquery.InnerText(htmlNode)))
		return nil
	}
	return errors.New("unable to find xpath: " + xpath)
}
