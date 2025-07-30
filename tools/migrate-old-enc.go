package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	u "github.com/sunshine69/golang-tools/utils"
	"golang.org/x/term"
)

func main() {
	key := flag.String("key", "", "Current encryption key")
	file := flag.String("f", filepath.Join(u.Must(os.UserHomeDir()), ".gnote.db.key"), "Working mode. If provided it will be your database key file (if your current database file encrytped then you will find a same <databasefile name>.key. This file need to be migrated.). Default your home dir /.gnote.db.key")
	inputcipher := flag.String("data", "", "Input text to decrypt. If provided the cli will decode this texst only and not migrate key file")
	flag.Parse()
	if *key == "" {
		fmt.Fprintln(os.Stderr, "Enter encryption key to decode: ")
		*key = string(u.Must(term.ReadPassword(int(syscall.Stdin))))
	}
	if *inputcipher != "" {
		fmt.Fprintln(os.Stdout, u.Must(u.Decrypt_v0(*inputcipher, *key)))
	} else {
		o := u.Must(u.Decrypt_v0(string(u.Must(os.ReadFile(*file))), *key))
		u.Copy(*file, *file+".backup")
		os.WriteFile(*file, []byte(u.Must(u.Encrypt(o, *key))), 0640)
	}
}
