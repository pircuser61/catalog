package log

import "fmt"

func Msg(text string) {
	fmt.Println(text)
}

func Msgf(format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	fmt.Println(str)
}

func ErrorMsg(text string) {
	fmt.Println(text)
}

func Warnf(format string, a ...any) {
	str := fmt.Sprintf(format, a...)
	fmt.Println(str)
}
