package envtpl

import (
	"cmp"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/flexstack/envtpl/pkg/nanoid"
	"github.com/flexstack/uuid"
	"github.com/joho/godotenv"
)

func Parse(file string) ([]EnvVar, error) {
	env, err := godotenv.Read(file)
	if err != nil {
		return nil, err
	}

	envVars := make([]EnvVar, len(env))
	i := 0
	for key, value := range env {
		val := ParseValue(value)
		envVars[i] = EnvVar{Key: key, Value: val}
		i++
	}

	// Sort the environment variables by key
	sort.Slice(envVars, func(i, j int) bool {
		return cmp.Less(envVars[i].Key, envVars[j].Key)
	})

	return envVars, nil
}

func ParseValue(value string) *EnvValue {
	if value == "" {
		return &EnvValue{
			Type:  Text,
			value: "",
		}
	}

	matches := valueRegex.FindStringSubmatch(value)
	if matches == nil {
		return &EnvValue{
			Type:  Plain,
			value: value,
		}
	}

	groups := make(map[string]string)
	for i, name := range valueRegex.SubexpNames() {
		if i != 0 && name != "" {
			groups[name] = matches[i]
		}
	}

	return &EnvValue{
		Type:  ValueType(groups["type"]),
		value: groups["value"],
	}
}

var valueRegex = regexp.MustCompile(`<(?P<type>text|password|enum|uuid|alpha|hex|base64|ascii85|int)(:(?P<value>.*))?>`)

type EnvVar struct {
	Key   string
	Value *EnvValue
}

type EnvValue struct {
	Type  ValueType
	value string
}

func (v *EnvValue) Generate() (interface{}, error) {
	if v.Type == Plain {
		return v.value, nil
	}

	switch v.Type {
	case Text, Password:
		return v.value, nil

	case Enum:
		possibleValues := strings.Split(v.value, ",")
		if len(possibleValues) == 0 {
			return nil, ErrInvalidArg
		}
		for i, value := range possibleValues {
			possibleValues[i] = strings.TrimSpace(value)
		}
		return possibleValues, nil

	case UUID:
		return uuid.Must(uuid.NewV4()).String(), nil

	case Int:
		min := int64(0)
		max := int64(100)
		var err error
		if v.value != "" {
			if strings.Contains(v.value, "-") {
				intRange := strings.Split(v.value, "-")
				if len(intRange) != 2 {
					return nil, ErrInvalidArg
				}
				min, err = strconv.ParseInt(intRange[0], 10, 64)
				if err != nil {
					return nil, ErrInvalidArg
				}
				max, err = strconv.ParseInt(intRange[1], 10, 64)
				if err != nil {
					return nil, ErrInvalidArg
				}
			}
		} else {
			max, err = strconv.ParseInt(v.value, 10, 64)
			if err != nil {
				return nil, ErrInvalidArg
			}
		}

		return fmt.Sprint(rand.Intn(int(max-min)) + int(min)), nil

	case Alpha, Hex, Base64, Ascii85:
		length := 16

		if v.value != "" {
			l, err := strconv.ParseInt(v.value, 10, 64)
			if err != nil {
				return nil, ErrInvalidArg
			}
			length = int(l)
		}

		alphabet := nanoid.AlphabetDefault
		switch v.Type {
		case Hex:
			alphabet = nanoid.AlphabetHex
		case Base64:
			alphabet = nanoid.AlphabetBase64
		case Ascii85:
			alphabet = nanoid.AlphabetAscii85
		}

		return nanoid.New(length, alphabet), nil
	}

	return nil, ErrInvalidArg
}

type ValueType string

const (
	Plain    ValueType = "plain"
	Text     ValueType = "text"
	Password ValueType = "password"
	Enum     ValueType = "enum"
	UUID     ValueType = "uuid"
	Alpha    ValueType = "alpha"
	Hex      ValueType = "hex"
	Base64   ValueType = "base64"
	Ascii85  ValueType = "ascii85"
	Int      ValueType = "int"
)

var (
	ErrInvalidArg = errors.New("invalid argument")
)
