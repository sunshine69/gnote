package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sunshine69/gnote/forms"
	u "github.com/sunshine69/golang-tools/utils"
)

func main() {
	gtk.Init(&os.Args)
	dbPath := flag.String("db", "", "Path to the database file")
	doMigrate := flag.Bool("mig", false, "Migrate")
	oldDB := flag.String("old-db", "", "Path to the old database file. If it is encrypted pass the key like filename?_pragma_key=x'<YOUR_KEY>'")
	createWinBundle := flag.Bool("create-win-bundle", false, "Create a windows bundle script")

	flag.Parse()

	if *createWinBundle {
		forms.CreateWinBundle()
		os.Exit(1)
	}

	binaryDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	u.CheckErr(err, "binaryDir")

	if _, e := os.Stat(fmt.Sprintf("%s/glade", binaryDir)); e != nil {
		forms.RestoreAssetsAll(binaryDir)
	}
	// For loading the glade resources etc.. DBPATH should be absolute path
	os.Chdir(binaryDir)

	u.CheckErr(err, "Getwd")
	homeDir, e := os.UserHomeDir()
	u.CheckErr(e, "UserHomeDir")

	var keyFile string = ""

	if *dbPath == "" {
		*dbPath = fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), ".gnote.db")
		fmt.Println("Use the database file in user home dir")
		keyFile = fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), ".gnote.db.key")
	} else {
		keyFile = *dbPath+".key"
	}
	var key, passphrase string
	var initialSetup bool = false
	if exist, _ := u.FileExists(keyFile); !exist {
		initialSetup = true
	}
	passphrase = forms.InputDialog("title", "Enter Passphrase", "label", "Enter passphrase to decode key. hit enter if you know your DB is not encrypted", "password-mask", '*')
	if passphrase != "" {
		if initialSetup {
			key, _ = u.RandomHex(32)
			encryptedKey := u.Encrypt(key, passphrase)
			err = os.WriteFile(keyFile, []byte(encryptedKey), 0600)
			u.CheckErr(err, "Write encrypted key file")
		} else {
			keyEncodedByte, err := os.ReadFile(keyFile)
			u.CheckErr(err, "keyEncodedByte")
			key, err = u.Decrypt(string(keyEncodedByte), passphrase)
			u.CheckErr(err, "Decode Key")
		}
	} else {
		key = ""
	}

	var fullDBPath string = ""
	switch key {
	case "":
		fullDBPath = *dbPath
	default:
		fullDBPath = fmt.Sprintf("%s?_pragma_key=x'%s'", *dbPath, key)
	}

	if *doMigrate {
		forms.DoMigrationV1(*oldDB, fullDBPath)
		os.Exit(0)
	}

	os.Setenv("DBPATH", fullDBPath)
	forms.SetupConfigDB()

	if _, e := forms.GetConfig("config_created"); e != nil {
		fmt.Println("Setup default config ....")
		forms.SetupDefaultConfig()
		forms.MessageBox("Initial setup db completed. The program will exit now. You can start it again.")
		os.Exit(0)
	}

	builder, err := gtk.BuilderNewFromFile("glade/gnote.glade")
	if err != nil {
		panic(err)
	}
	gnoteApp := forms.GnoteApp{
		Builder: builder,
	}

	gnoteApp.InitApp()
	gtk.Main()
}
