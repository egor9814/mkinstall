package main

import "strings"

type Lang uint
type LangArray []Lang

const (
	English Lang = iota
	Russian

	Default = English
)

func (l Lang) String() string {
	switch l {
	case English:
		return "English"

	case Russian:
		return "Русский"

	default:
		panic("unsupported language")
	}
}

func (la LangArray) Strings() []string {
	res := make([]string, len(la))
	for i, it := range la {
		res[i] = it.String()
	}
	return res
}

func (la LangArray) String() string {
	return strings.Join(la.Strings(), ", ")
}

var AllLangs = LangArray{
	English,
	Russian,
}

var Yes = [...]string{
	English: "Yes",
	Russian: "Да",
}

var No = [...]string{
	English: "No",
	Russian: "Нет",
}
