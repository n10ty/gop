package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// path to keychain file
var keychainPath = filepath.Join(os.Getenv("HOME"), ".gop")

// default number of totp digits
var digits = 6

// character to separate values in a keychain record
var keychainRowDelimiter = " "

type Keychain struct {
	keys map[string]string
}

func (k *Keychain) add(name, key string) error {
	if _, exists := k.keys[name]; exists {
		return fmt.Errorf("key '%s' alreay exists", name)
	}

	k.keys[name] = key

	return nil
}

func (k *Keychain) get(name string) ([]byte, error) {
	key, exists := k.keys[name]
	if !exists {
		return nil, fmt.Errorf("key '%s' does not exist", name)
	}

	return base32.StdEncoding.DecodeString(strings.ToUpper(key))
}

var rootCommand = &cobra.Command{
	Use: "gop",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	DisableFlagsInUseLine: true,
}

var addCommand = &cobra.Command{
	Use:                   "add <name> <secret key>",
	Short:                 "Add new key",
	Args:                  cobra.MinimumNArgs(2),
	Aliases:               []string{"a"},
	RunE:                  addKey,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
}

var generateCommand = &cobra.Command{
	Use:                   "generate <name>",
	Short:                 "Generate password and copy to clipboard",
	Args:                  cobra.MinimumNArgs(1),
	Aliases:               []string{"g"},
	RunE:                  generateTotp,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
}

var listCommand = &cobra.Command{
	Use:                   "list",
	Short:                 "List all keys",
	Aliases:               []string{"l"},
	RunE:                  list,
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
}

func main() {
	rootCommand.AddCommand(addCommand)
	rootCommand.AddCommand(generateCommand)
	rootCommand.AddCommand(listCommand)
	rootCommand.Execute()
}

func addKey(_ *cobra.Command, args []string) error {
	keychain, err := readKeychain(keychainPath)
	if err != nil {
		return err
	}

	name := args[0]
	key := args[1]

	if err := keychain.add(name, key); err != nil {
		return err
	}
	if err = writeKeychain(keychain, keychainPath); err != nil {
		return err
	}

	return nil
}

func generateTotp(_ *cobra.Command, args []string) error {
	keychain, err := readKeychain(keychainPath)
	if err != nil {
		return err
	}

	name := args[0]

	key, err := keychain.get(name)
	if err != nil {
		return err
	}

	p := totp(key, time.Now(), digits)
	code := fmt.Sprintf("%0*d\n", digits, p)
	err = clipboard.WriteAll(code)
	if err != nil {
		return err
	}
	fmt.Print(code)

	return nil
}

func list(_ *cobra.Command, _ []string) error {
	keychain, err := readKeychain(keychainPath)
	if err != nil {
		return err
	}

	for name := range keychain.keys {
		fmt.Println(name)
	}

	return nil
}

func readKeychain(keychainPath string) (*Keychain, error) {
	keychain := Keychain{
		keys: map[string]string{},
	}
	f, err := os.OpenFile(keychainPath, os.O_CREATE|os.O_RDWR, 0600)
	defer f.Close()
	if err != nil {
		return nil, fmt.Errorf("error during open keychain: %s", err)
	}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		t := scanner.Text()
		entry := strings.Split(t, keychainRowDelimiter)
		if len(entry) < 2 {
			return nil, fmt.Errorf("corrupted keychain file")
		}

		keychain.keys[entry[0]] = entry[1]
	}

	return &keychain, nil
}

func writeKeychain(keychain *Keychain, keychainPath string) error {
	f, err := os.OpenFile(keychainPath, os.O_CREATE|os.O_RDWR, 0600)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("error during open keychain: %s", err)
	}

	for name, key := range keychain.keys {
		row := []string{
			strings.TrimSpace(name),
			strings.TrimSpace(key),
		}
		_, err := fmt.Fprintf(f, "%s\n", strings.Join(row, keychainRowDelimiter))
		if err != nil {
			return err
		}
	}

	return nil
}

func hotp(key []byte, counter uint64, digits int) int {
	h := hmac.New(sha1.New, key)
	binary.Write(h, binary.BigEndian, counter)
	hmacResult := h.Sum(nil)
	offset := hmacResult[19] & 0x0F
	code := uint32(hmacResult[offset])&0x7f<<24 |
		uint32(hmacResult[offset+1])&0xff<<16 |
		uint32(hmacResult[offset+2])&0xff<<8 |
		uint32(hmacResult[offset+3])&0xff
	d := uint32(1)
	for i := 0; i < digits && i < 8; i++ {
		d *= 10
	}

	return int(code % d)
}

func totp(key []byte, t time.Time, digits int) int {
	return hotp(key, uint64(t.UnixNano())/30e9, digits)
}
