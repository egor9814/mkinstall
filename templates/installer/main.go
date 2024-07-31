package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

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

func welcomePage() error {
	title := [...]string{
		English: "Welcome!",
		Russian: "Добро пожаловать!",
	}
	message := [...]string{
		English: "Welcome to the installation wizard.",
		Russian: "Вас приветсвует мастер установки.",
	}
	return zenity.Info(message[lang], zenity.Title(title[lang]))
}

func installPathPage() error {
	title := [...]string{
		English: "Install path",
		Russian: "Путь установки",
	}
	file, err := zenity.SelectFile(
		zenity.Title(title[lang]),
		zenity.Directory(),
		zenity.Filename(install.Target.Path),
	)
	if err == nil {
		install.Target.Path = path.Join(file, install.Product.Name)
	}
	return err
}

func summaryPage() error {
	title := [...]string{
		English: "Summary",
		Russian: "Итог",
	}
	message := [...]string{
		English: "The program will be installed in%qContinue?",
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
	if install.Target.Editable {
		options = append(options, zenity.ExtraButton(changePath[lang]))
	}
	for {
		err := zenity.Question(
			fmt.Sprintf(message[lang], install.Target.Path),
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
	dialog, err := zenity.Progress(
		zenity.Title(title[lang]),
		zenity.TimeRemaining(),
	)
	if err != nil {
		return err
	}

	dialog.Text(preparing[lang])
	dialog.Value(0)

	input, err := NewInput()
	if err != nil {
		return err
	}
	for {
		it, err := input.Next()
		if err != nil {
			dialog.Close()
			showError(err)
			return err
		}
		if !it.IsValid() {
			break
		}
		dialog.Text(fmt.Sprintf(extracting[lang], it.Path))

		outFile := &VirtualFile{Path: path.Join(install.Target.Path, it.Path)}
		out, err := outFile.Create()
		if err != nil {
			dialog.Close()
			showError(err)
			return err
		}

		rc, err := it.Open()
		if err != nil {
			dialog.Close()
			out.Close()
			input.Close()
			return err
		}

		if _, err := io.Copy(out, rc); err != nil {
			dialog.Close()
			rc.Close()
			out.Close()
			input.Close()
			return err
		}

		out.Close()
		rc.Close()
		dialog.Value(100.0 * input.ProgressCurrent() / input.ProgressAll())
	}

	input.Close()
	dialog.Complete()
	dialog.Close()

	return nil
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
	return path.Base(os.Args[0])
}

func main() {
	for i, l := 1, len(os.Args); i < l; i++ {
		it := os.Args[i]
		switch it {
		case "version":
			parseVersion()
			fmt.Printf("%s v%d.%d.%d%s\n", exe(), Version.Major, Version.Minor, Version.Patch, Version.Suffix)
			return

		default:
			// skip
		}
	}
	handleError(languagePage(), false)

	handlePage(welcomePage)

	showError(install.load())

	if install.Target.Editable {
		// originalPath := install.Target.Path
		handlePage(installPathPage)
	} else {
		install.Target.Path = path.Join(install.Target.Path, install.Product.Name)
	}

	handlePage(summaryPage)

	handlePage(installPage)

	handlePage(finishPage)
}
