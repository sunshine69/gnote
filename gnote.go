package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/gotk3/gotk3/gtk"
	"github.com/sunshine69/gnote/forms"
	u "github.com/sunshine69/golang-tools/utils"
	"golang.org/x/term"
)

var (
	cliCommand = map[string]func() error{
		"migrate-key-file": MigrateKeyFile,
	}
)
var (
	dbPath = flag.String("db", "", "Path to the database file")

	cliflag            = flag.NewFlagSet("cli", flag.ExitOnError)
	doMigrate          = cliflag.Bool("mig", false, "Migrate")
	oldDB              = cliflag.String("old-db", "", "Path to the old database file. If it is encrypted pass the key like filename?_pragma_key=x'<YOUR_KEY>'")
	createWinBundle    = cliflag.Bool("create-win-bundle", false, "Create a windows bundle script")
	mingw64Prefix      = cliflag.String("mingw64-root", "c:/tools/msys64/mingw64", "Mingw64 root dir. Under this we have the /bin dir which has all gtk dll files")
	cliCmd             = cliflag.String("command", "", "Cli commands. List commands: "+u.JsonDump(u.MapKeysToSlice(cliCommand), ""))
	keyfilePass        = cliflag.String("keyfile-pass", "", "Current keyfile password. Will prompt if it is empty")
	inputcipher        = cliflag.String("data", "", "Input text to decrypt/encrypt. If provided the cli will action on this text only and not migrate key file")
	inputcipherVersion = cliflag.Int("data-version", 1, "The version that the data is encrypted with. If 0 it is the old depricated version. Other than that it could be a version number, like 1 or 2 matching with encryption config on utils")
	action             = cliflag.String("action", "decrypt", "What to do, encrypt or decrypt")
)

func MigrateKeyFile() error {
	keyFile := os.Getenv("GNOTE_KEY_FILE")
	if keyFile == "" {
		fmt.Fprintln(os.Stderr, "[INFO] env var GNOTE_KEY_FILE is empty. Will use the default .gnote.key in your user home dir")
		keyFile = filepath.Join(u.Must(os.UserHomeDir()), ".gnote.db.key")
	}
	if *keyfilePass == "" {
		fmt.Fprintln(os.Stderr, "Enter encryption key to decode/encode: ")
		*keyfilePass = string(u.Must(term.ReadPassword(int(syscall.Stdin))))
	}

	if *inputcipher != "" {
		o := ""
		switch *inputcipherVersion {
		case 0:
			fmt.Fprintf(os.Stderr, "[WARN] use old and weak Encrypt_v0 func\n")
			switch *action {
			case "encrypt":
				o = u.Must(u.Encrypt_v0(*inputcipher, *keyfilePass))
			case "decrypt":
				o = u.Must(u.Decrypt_v0(*inputcipher, *keyfilePass))
			}
			fmt.Fprintln(os.Stdout, o)
		default:
			encCfg := u.Must(u.NewEncConfigForVersion(byte(*inputcipherVersion)))
			switch *action {
			case "encrypt":
				o = u.Must(u.Encrypt(*inputcipher, *keyfilePass, encCfg))
			case "decrypt":
				o = u.Must(u.Decrypt(*inputcipher, *keyfilePass, encCfg))
			}
			fmt.Fprintln(os.Stdout, o)
		}
	} else {
		o := u.Must(u.Decrypt_v0(string(u.Must(os.ReadFile(keyFile))), *keyfilePass))
		u.CheckErr(u.Copy(keyFile, keyFile+".backup"), "Backup current file")
		u.CheckErr(os.WriteFile(keyFile, []byte(u.Must(u.Encrypt(o, *keyfilePass, nil))), 0640), "Write new key file")
	}
	return nil
}

func main() {
	if len(os.Args) == 1 {
		*dbPath = ""
	} else {
		switch os.Args[1] {
		case "cli":
			cliflag.Parse(os.Args[2:])
		default:
			flag.Parse()
		}
	}

	if *createWinBundle {
		forms.CreateWinBundle(*mingw64Prefix)
		os.Exit(0)
	}
	if *cliCmd != "" { // Process cli utilities
		if runFunc, ok := cliCommand[*cliCmd]; ok {
			u.CheckErr(runFunc(), "Run "+*cliCmd)
		}
		os.Exit(0)
	}

	gtk.Init(&os.Args)
	binaryDir := u.Must(filepath.Abs(filepath.Dir(os.Args[0])))

	if _, e := os.Stat(fmt.Sprintf("%s/glade", binaryDir)); e != nil {
		forms.RestoreAssetsAll(binaryDir)
	}
	// For loading the glade resources etc.. DBPATH should be absolute path
	u.Must("", os.Chdir(binaryDir))

	builder := u.Must(gtk.BuilderNewFromFile("glade/gnote.glade"))

	gnoteApp := forms.GnoteApp{
		Builder: builder,
	}
	// DB will only available after DoStartup to prompt passphrase etc. InitApp only init graphic resources
	gnoteApp.InitApp()
	DoStartup()
	gnoteApp.SetDefaultWindowSize() // This one need the db to get the saved sizes
	gnoteApp.MainWindow.ShowAll()
	gtk.Main()
}

func DoStartup() {
	var keyFile string = ""
	homeDir := u.Must(os.UserHomeDir())
	if *dbPath == "" {
		*dbPath = fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), ".gnote.db")
		fmt.Println("Use the database file in user home dir")
		keyFile = fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), ".gnote.db.key")
	} else {
		keyFile = *dbPath + ".key"
	}
	var key, passphrase string
	var initialSetup bool = false
	if exist, _ := u.FileExists(keyFile); !exist {
		initialSetup = true
	}
	passphrase = forms.InputDialog("title", "Enter Passphrase", "label", "Enter passphrase to decode key. hit enter if you know your DB is not encrypted", "password-mask", '*')

	if passphrase != "" {
		if initialSetup {
			// It needs to be hex encoded 32 byte key; see the dsn below. See https://github.com/mutecomm/go-sqlcipher
			key, _ = u.RandomHex(32)
			encryptedKey, _ := u.Encrypt(key, passphrase, nil)
			u.Must("[ERROR] os.WriteFile", os.WriteFile(keyFile, []byte(encryptedKey), 0600))
		} else {
			keyEncodedByte := u.Must(os.ReadFile(keyFile))
			key = u.Must(u.Decrypt(string(keyEncodedByte), passphrase, nil))
		}
	} else {
		key = ""
	}

	var fullDBPath string = ""
	switch key {
	case "":
		fullDBPath = *dbPath
	default:
		// To re-key To change the key on an existing encrypted database, it must first be unlocked with the current encryption key. Once the database is readable and writeable, PRAGMA rekey can be used to re-encrypt every page in the database with a new key.
		// Example see https://www.zetetic.net/sqlcipher/sqlcipher-api/index.html#rekey
		// sqlite> PRAGMA key = x'old hex';
		// sqlite> PRAGMA rekey = x'new hex'; In go use the func Query
		fullDBPath = fmt.Sprintf("%s?_pragma_key=x'%s'", *dbPath, key)
	}
	// fmt.Println(fullDBPath)
	os.Setenv("DBPATH", fullDBPath)

	if *doMigrate {
		forms.DoMigrationV1(*oldDB, fullDBPath)
		os.Exit(0)
	}
	// forms.SetupConfigDB()

	config_created, err := forms.GetConfig("config_created")
	u.CheckErrNonFatal(err, "GetConfig config_created")
	if config_created == "" {
		fmt.Println("Setup default config ....")
		forms.SetupConfigDB()
		forms.SetupDefaultConfig()
		forms.MessageBox("Initial setup db completed.")
		// os.Exit(0)
	}
	forms.DateLayout, _ = forms.GetConfig("date_layout")
	forms.WebNoteUser, _ = forms.GetConfig("webnote_user")
}
