package main

import (
	"path/filepath"
	"flag"
	"fmt"
	"os"
	"github.com/gotk3/gotk3/gtk"
	"github.com/sunshine69/gnote/forms"
)

func main() {
	gtk.Init(&os.Args)
	dbPath := flag.String("db","","Path to the database file")
	doMigrate := flag.Bool("mig",false,"Migrate")
	flag.Parse()
	if *doMigrate {
		forms.DoMigration()
		os.Exit(0)
	}

	workdir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
			fmt.Printf("%v\n",err)
			os.Exit(1)
	}

	if _, e := os.Stat(fmt.Sprintf("%s/glade", workdir) ); e == nil {
		os.Chdir(workdir)

	}

	homeDir, e := os.UserHomeDir()
	if e != nil {
		fmt.Printf("ERROR %v\n", e)
	}
	if *dbPath == "" {
		*dbPath =  fmt.Sprintf("%s%s%s", homeDir, string(os.PathSeparator), ".gnote.db")
		fmt.Printf("Use the database file %s\n", *dbPath)
	}
	os.Setenv("DBPATH", *dbPath)
	forms.SetupConfigDB()

	if _, e := forms.GetConfig("config_created"); e != nil {
		fmt.Println("Setup default config ....")
		forms.SetupDefaultConfig()
		forms.RestoreAssetsAll(workdir)
		forms.MessageBox("Initial setup db compelted. The program will exit now. You can start it again.")
		os.Exit(0)
	}

	builder, err := gtk.BuilderNewFromFile("glade/gnote.glade")
	if err != nil {
		panic(err)
	}
	gnoteApp := forms.GnoteApp {
		Builder: builder,
	}

	gnoteApp.InitApp()
	gtk.Main()
}