package main

import (
	"os"
	"fmt"
	"flag"
	s "github.com/xiak/remote-command/pkg/services"
)

var version = "0.0.1"
var commit = "Xiak"

var fn = []string {
	"Data domain file system GC",
	"Data domain restart ddboost",
	"Data domain delete all storage unit",
	"Data domain delete oldest storage unit",
	"Data movement to cloud",
	"Avamar GC",
}

func fns (fn *[]string) *string {
	var temp string
	for key, value := range *fn {
		temp = fmt.Sprintf("%s%d: %s\n\t", temp, key, value)
	}
	return &temp
}

func main() {
	v := flag.Bool("v", false, "prints current version and exits")
	host := flag.String("host", "", "Remote server host name or ip address")
	user := flag.String("user", "", "Remote server user name")
	password := flag.String("password", "", "Remote server password")
	x := flag.Int("x", -1, *fns(&fn))
	flag.Parse()
	if *v {
		fmt.Printf("Version %s (commit: %s)\n", version, commit)
		os.Exit(0)
	}
	if *host == "" || *user == "" || *password == "" {
		fmt.Printf("Please input available host (-host), user (-user), password (-password)\n")
		flag.Usage()
		os.Exit(0)
	}
	if *x < 0 || *x > len(fn) {
		fmt.Printf("Please input available number to execute\n")
		flag.Usage()
		os.Exit(0)
	}

	cmd := s.NewCommand(*host, *user, *password)
	dd := s.NewDataDomain(cmd)
	if *x == 0 {
		err := dd.FilesysGc()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *x == 1 {
		err := dd.DDBoostRestart()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *x == 2 {
		err := dd.DeleteAllStorageUnit()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *x == 3 {
		err := dd.DeleteOldestStorageUnit()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	if *x == 4 {
		err := dd.DataMoveToCloud()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	avamar := s.NewAvamar(cmd)
	if *x == 5 {
		err := avamar.Gc()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	fmt.Println("No function was executed\n")
	os.Exit(1)
}