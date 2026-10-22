package opt

import (
	"os"
	"strings"
)

var (
	Options = make(map[string]Value)
	allArgs []string
	Args []string
)

type Value interface {
	NextIsOption() bool
	SetValue(string)
}

type StringFlag struct {
	NextOption bool
	Value *string
}

type BoolFlag struct {
	NextOption bool
	Value *bool
}

func (stringValue StringFlag) SetValue(arg string) {
	*stringValue.Value = arg
}

func (boolValue BoolFlag) SetValue(arg string) {
	*boolValue.Value = true
}

func (stringValue StringFlag) NextIsOption() bool {
	return stringValue.NextOption
}

func (boolValue BoolFlag) NextIsOption() bool {
	return boolValue.NextOption
}

func StringVar(v *string, shortName string, longName string, defaultValue string) {
	ShortStringVar(v, shortName, defaultValue)
	LongStringVar(v, longName, defaultValue)
}

func ShortStringVar(v *string, name string, defaultValue string) {
	*v = defaultValue
	Options["-"+name] = StringFlag{
		NextOption: true,
		Value: v,
	}
}

func LongStringVar(v *string, name string, defaultValue string) {
	*v = defaultValue
	Options["--"+name] = StringFlag{
		NextOption: true,
		Value: v,
	}
}

func BoolSwitch(v *bool, shortName string, longName string) {
	ShortBoolSwitch(v, shortName)
	LongBoolSwitch(v, longName)
}

func LongBoolSwitch(v *bool, name string) {
	Options["--"+name] = BoolFlag{
		Value: v,
	}
}

func ShortBoolSwitch(v *bool, name string) {
	Options["-"+name] = BoolFlag{
		Value: v,
	}
}

func Parse() {
	initArgs()
	scan()
}

func scan() {
	var nextValue Value
	if len(os.Args) < 2 {
		return
	}

	for _, arg := range allArgs {
		switch true {
		case nextValue != nil && nextValue.NextIsOption():
			nextValue.SetValue(arg)
			nextValue = nil

		case strings.HasPrefix(arg, "-"):
			var ok bool
			nextValue, ok = Options[arg]
			if ok && !nextValue.NextIsOption() {
				nextValue.SetValue("")
			}

		default:
			Args = append(Args, arg)
		}
	}
}

func initArgs() {
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-") && strings.ContainsRune(arg, '=') {
			fields := strings.SplitN(arg, "=", 2)
			value, ok := Options[fields[0]]
			switch true {
			case !ok:
				continue

			case value.NextIsOption():
				allArgs = append(allArgs, fields...)

			default:
				allArgs = append(allArgs, fields[0])
			}
		}
		allArgs = append(allArgs, arg)
	}
}
