package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

// Sweeper .
type Sweeper struct {
	modsDir string
	modIDs  []string
}

// NewSweeper .
func NewSweeper() *Sweeper {
	s := &Sweeper{}
	s.init()

	return s
}

func (s *Sweeper) init() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 581320`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	l, _, err := k.GetStringValue("InstallLocation")
	if err != nil {
		log.Fatal(err)
	}

	s.modsDir = filepath.Join(l, "Insurgency", "Mods", "modio")

	_, err = os.Stat(s.modsDir)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Sweeper) execute() {
	fmt.Printf("Mods directory: %s\n", s.modsDir)

	e, err := ioutil.ReadDir(s.modsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, _e := range e {
		if _e.IsDir() {
			s.modIDs = append(s.modIDs, _e.Name())
			continue
		}
	}

	var delCnt int

	for _, modID := range s.modIDs {
		e, err := ioutil.ReadDir(filepath.Join(s.modsDir, modID))
		if err != nil {
			log.Fatal(err)
		}

		var _modFiles []int
		for _, _e := range e {
			if _e.IsDir() {
				fID, _ := strconv.Atoi(_e.Name())
				_modFiles = append(_modFiles, fID)
				continue
			}
		}

		sort.Ints(_modFiles)

		if len(_modFiles) == 1 {
			fmt.Printf("[%s] Skip\n", modID)
			continue
		}

		for _, _fID := range _modFiles[:len(_modFiles)-1] {
			if err := os.RemoveAll(filepath.Join(s.modsDir, modID, strconv.Itoa(_fID))); err != nil {
				fmt.Printf("[%s] %d FAILED. Skip.\n", modID, _fID)
				continue
			}

			fmt.Printf("[%s] %d DELETED.\n", modID, _fID)
			delCnt = delCnt + 1
		}
	}

	fmt.Printf("%d Mod(s) found.\n", len(s.modIDs))
	if delCnt > 0 {
		fmt.Printf("Delete: %d File(s) SUCCEEDED.\n", delCnt)
	} else {
		fmt.Println("There's no need to clean it up.")
	}
	fmt.Println("Done.")
	os.Exit(0)
}

func main() {
	fmt.Println("+++++++++++++++++++++++++++++++++++++++")
	fmt.Println("+ Mod Sweeper for Insurgency: Sandstorm ")
	fmt.Println("+ Author: Franky <franky@jpnhq.net>")
	fmt.Println("+ Copyright: 2020 JPNHQ")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++")
	NewSweeper().execute()
}
