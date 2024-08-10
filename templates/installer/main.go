package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ncruces/zenity"
)

var lang Lang = Default

func languagePage() error {
	res, err := zenity.List(
		"Select language:",
		AllLangs.Strings(),
		zenity.DisallowEmpty(),
		zenity.RadioList(),
		zenity.Title("Language"),
		zenity.DefaultItems(Default.String()), // not work :(
	)
	if err != nil {
		return err
	}
	for _, it := range AllLangs {
		if it.String() == res {
			lang = it
			return nil
		}
	}
	panic("unreachable")
}

func zstdMemoryPage() error {
	title := [...]string{
		English: "Settings",
		Russian: "Параметры",
	}
	message := [...]string{
		English: "Memory limit:",
		Russian: "Ограничение памяти:",
	}
	available := []string{
		"1G",
		"2G",
		"4G",
		"8G",
	}
	res, err := zenity.List(
		message[lang],
		available,
		zenity.DisallowEmpty(),
		zenity.RadioList(),
		zenity.Title(title[lang]),
		zenity.DefaultItems(available[len(available)-1]),
	)
	if err != nil {
		return err
	}
	zstdInputInfo.maxMem = 1
	for _, it := range available {
		if it == res {
			zstdInputInfo.maxMem *= 1024 * 1024 * 1024
			return nil
		}
		zstdInputInfo.maxMem *= 2
	}
	panic("unreachable")
}

func makeShortcutPage() error {
	if len(install.Shortcuts) == 0 {
		return nil
	}
	title := [...]string{
		English: "Settings",
		Russian: "Параметры",
	}
	message := [...]string{
		English: "Create shortcuts?",
		Russian: "Создать ярлыки?",
	}
	err := zenity.Question(
		message[lang],
		zenity.Title(title[lang]),
		zenity.OKLabel(Yes[lang]),
		zenity.ExtraButton(No[lang]),
	)
	if err == zenity.ErrExtraButton {
		install.Shortcuts = nil
	}
	return err
}

func installPathPage() error {
	title := [...]string{
		English: "Install path",
		Russian: "Путь установки",
	}
	file, err := zenity.SelectFile(
		zenity.Title(title[lang]),
		zenity.Directory(),
		zenity.Filename(filepath.FromSlash(install.TargetPath)),
	)
	if err == nil {
		install.TargetPath = filepath.Join(file, install.ProductName)
	}
	return err
}

func summaryPage() error {
	title := [...]string{
		English: "Summary",
		Russian: "Итог",
	}
	message := [...]string{
		English: "The program will be installed in %q\nContinue?",
		Russian: "Программа будет установлена в %q\nПродолжить?",
	}
	changePath := [...]string{
		English: "Change path",
		Russian: "Изменить путь",
	}
	options := make([]zenity.Option, 0, 4)
	options = append(options,
		zenity.Title(title[lang]),
		zenity.OKLabel(Yes[lang]),
		zenity.CancelLabel(No[lang]),
	)
	if install.TargetPathEditable {
		options = append(options, zenity.ExtraButton(changePath[lang]))
	}
	for {
		err := zenity.Question(
			fmt.Sprintf(message[lang], filepath.ToSlash(install.TargetPath)),
			options...,
		)
		if err == zenity.ErrExtraButton {
			if err := installPathPage(); err != nil {
				if err == zenity.ErrCanceled {
					continue
				}
				return err
			}
			continue
		}
		return err
	}
}

func installPage() error {
	input, err := install.InputType.Open()
	if err != nil {
		showError(err)
		return err
	}

	title := [...]string{
		English: "Installing",
		Russian: "Установка",
	}
	preparing := [...]string{
		English: "Preparing...",
		Russian: "Подготовка...",
	}
	extracting := [...]string{
		English: "Extracting %q...",
		Russian: "Распаковка %q...",
	}

	dialogMtx := sync.Mutex{}
	makeDialog := func() (zenity.ProgressDialog, error) {
		return zenity.Progress(
			zenity.Title(title[lang]),
			zenity.TimeRemaining(),
			zenity.MaxValue(int(input.Progress().All())),
		)
	}
	dialog, err := makeDialog()
	if err != nil {
		showError(err)
		return err
	}

	dialog.Text(preparing[lang])
	dialog.Value(0)

	errChan := make(chan error, 1)
	mtx := sync.Mutex{}
	pause := false

	go func() {
		for it := range input.Progress().Chan() {
			dialogMtx.Lock()
			dialog.Value(int(it))
			dialogMtx.Unlock()
		}
	}()

	go func() {
		for {
			mtx.Lock()
			p := pause
			mtx.Unlock()
			if p {
				time.Sleep(time.Second)
				continue
			}
			it, err := input.Next()
			if err != nil {
				input.Close()
				dialog.Close()
				errChan <- err
				return
			}
			if !it.IsValid() {
				break
			}
			dialog.Text(fmt.Sprintf(extracting[lang], it.Path))

			outFile := &VirtualFile{Path: filepath.Join(install.TargetPath, filepath.FromSlash(it.Path))}
			out, err := outFile.Create()
			if err != nil {
				input.Close()
				dialog.Close()
				errChan <- err
				return
			}

			rc, err := it.Open()
			if err != nil {
				out.Close()
				input.Close()
				dialog.Close()
				errChan <- err
				return
			}

			if _, err := io.Copy(out, rc); err != nil {
				rc.Close()
				out.Close()
				input.Close()
				dialog.Close()
				errChan <- err
				return
			}

			rc.Close()
			out.Close()
		}
		errChan <- nil
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				showError(err)
				return err
			}
		case <-dialog.Done():
			mtx.Lock()
			pause = true
			mtx.Unlock()
			if !handleError(zenity.ErrCanceled, true) {
				dialogMtx.Lock()
				dialog, err = makeDialog()
				dialogMtx.Unlock()
				if err != nil {
					return err
				}
				mtx.Lock()
				pause = false
				mtx.Unlock()
			}
			continue
		case <-context.Background().Done():
		}
		break
	}

	input.Close()
	dialog.Close()

	return nil
}

func makeShortcuts() error {
	l := len(install.Shortcuts)
	if l == 0 {
		return nil
	}
	errList := make([]error, 0, l)
	for _, it := range install.Shortcuts {
		if err := it.Make(); err != nil {
			errList = append(errList, err)
		}
	}
	if l := len(errList); l == 0 {
		return nil
	} else if l == 1 {
		return errList[0]
	} else {
		s := make([]string, l)
		for i, it := range errList {
			s[i] = it.Error()
		}
		return errors.New(strings.Join(s, "\n"))
	}
}

func finishPage() error {
	title := [...]string{
		English: "Finish",
		Russian: "Завершение",
	}
	message := [...]string{
		English: "Installation successful finsished",
		Russian: "Установка успешно завершена",
	}
	zenity.Info(message[lang], zenity.Title(title[lang]))
	return nil
}

func handleError(err error, sureCancel bool) bool {
	if err == zenity.ErrCanceled {
		if sureCancel {
			title := [...]string{
				English: "Install cancelation",
				Russian: "Отмена установки",
			}
			message := [...]string{
				English: "Are you sure?",
				Russian: "Вы уверены?",
			}
			if err := zenity.Question(
				message[lang],
				zenity.Title(title[lang]),
				zenity.OKLabel(Yes[lang]),
				zenity.CancelLabel(No[lang]),
				zenity.DefaultCancel(),
			); err != nil {
				if err == zenity.ErrCanceled {
					return false
				}
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func handlePage(f func() error) {
	for {
		if handleError(f(), true) {
			break
		}
	}
}

func showError(err error) {
	if err == nil {
		return
	}
	zenity.Error(err.Error(), zenity.Title("Fatal error"))
	os.Exit(1)
}

func exe() string {
	return filepath.Base(os.Args[0])
}

func main() {
	for i, l := 1, len(os.Args); i < l; i++ {
		it := os.Args[i]
		switch it {
		case "version":
			fmt.Printf("%s v%d.%d.%d%s\n", exe(), Version.Major, Version.Minor, Version.Patch, Version.Suffix)
			return

		case "help":
			fmt.Printf("Usage: %s [COMMAND]\n", exe())
			fmt.Println("Commands:")
			fmt.Println(" help           - Print this help")
			fmt.Println(" version        - Print installer version")
			return

		default:
			// skip
		}
	}
	handleError(languagePage(), false)

	showError(install.init())

	if install.TargetPathEditable {
		handlePage(installPathPage)
	} else {
		install.TargetPath = filepath.Join(install.TargetPath, install.ProductName)
	}

	handlePage(zstdMemoryPage)

	handlePage(makeShortcutPage)

	handlePage(summaryPage)

	handlePage(installPage)

	showError(makeShortcuts())

	handlePage(finishPage)
}
