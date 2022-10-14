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

	workdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if _, e := os.Stat(fmt.Sprintf("%s/glade", workdir)); e != nil {
		forms.RestoreAssetsAll(workdir)
	}
	os.Chdir(workdir)

	homeDir, e := os.UserHomeDir()
	if e != nil {
		fmt.Printf("ERROR %v\n", e)
	}
	if *dbPath == "" {
		*dbPath = fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), ".gnote.db")
		fmt.Printf("Use the database file %s\n", *dbPath)
	}

	key := forms.InputDialog("title", "Enter decode 32 bytes key (64 char long)", "label", "Enter decode key. Type auto to auto generate. To disable leave it empty just hit enter", "password-mask", '*')
	var fullDBPath string = ""
	switch key {
	case "auto":
		key, _ = u.RandomHex(32)
		fmt.Printf("[INFO] HERE IS YOUR KEY. WRITE IT DOWN SAVE TO SOMWHERE. IF GET LOST ALL YOUR FUTURE DATA WILL BE GONE\n%s\n", key)
		fullDBPath = fmt.Sprintf("%s?_pragma_key=x'%s'", *dbPath, key)
	case "":
		fmt.Println("[INFO] Encryption is disabled")
		fullDBPath = *dbPath
	default:
		if len(key) == 64 {
			fullDBPath = fmt.Sprintf("%s?_pragma_key=x'%s'", *dbPath, key)
		} else {
			fmt.Printf("[WARN] key length is not 64 char long, so use non hex key\n")
			fullDBPath = fmt.Sprintf("%s?_pragma_key='%s'", *dbPath, key)
		}
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
