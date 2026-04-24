package input

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/tailscale/hujson"
	toon "github.com/toon-format/toon-go"
)

var errEmptyInput = errors.New("input is empty")

func Parse(data []byte) (any, error) {
	trimmed := bytes.TrimSpace(stripUTF8BOM(data))
	if len(trimmed) == 0 {
		return nil, errEmptyInput
	}

	value, err := parseStream(trimmed)
	if err == nil {
		return value, nil
	}

	standardized, stdErr := hujson.Standardize(trimmed)
	if stdErr != nil {
		return nil, fmt.Errorf("parse input: %w", err)
	}

	value, err = parseSingleDocument(standardized)
	if err != nil {
		return nil, fmt.Errorf("parse jsonc input: %w", err)
	}

	return value, nil
}

func stripUTF8BOM(data []byte) []byte {
	return bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
}

func parseStream(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))

	values := make([]any, 0, 1)
	for {
		var raw json.RawMessage
		err := decoder.Decode(&raw)
		if err == nil {
			value, err := parseDocument(raw)
			if err != nil {
				return nil, err
			}
			values = append(values, value)
			continue
		}

		if errors.Is(err, io.EOF) {
			break
		}

		return nil, err
	}

	if len(values) == 0 {
		return nil, errEmptyInput
	}

	if len(values) == 1 {
		return values[0], nil
	}

	return values, nil
}

func parseSingleDocument(data []byte) (any, error) {
	value, err := parseStream(data)
	if err != nil {
		return nil, err
	}

	if _, ok := value.([]any); ok {
		return nil, errors.New("jsonc input must contain exactly one top-level value")
	}

	return value, nil
}

func parseDocument(raw []byte) (any, error) {
	value, err := hujson.Parse(raw)
	if err != nil {
		return nil, err
	}

	return convertValue(value)
}

func convertValue(value hujson.Value) (any, error) {
	switch trimmed := value.Value.(type) {
	case *hujson.Object:
		return convertObject(trimmed)
	case *hujson.Array:
		return convertArray(trimmed)
	case hujson.Literal:
		return convertLiteral(trimmed)
	default:
		return nil, fmt.Errorf("unsupported hujson value type: %T", value.Value)
	}
}

func convertObject(object *hujson.Object) (any, error) {
	fields := make([]toon.Field, 0, len(object.Members))

	for _, member := range object.Members {
		key, err := convertObjectKey(member.Name)
		if err != nil {
			return nil, err
		}

		value, err := convertValue(member.Value)
		if err != nil {
			return nil, err
		}

		fields = append(fields, toon.Field{Key: key, Value: value})
	}

	return toon.NewObject(fields...), nil
}

func convertArray(array *hujson.Array) (any, error) {
	items := make([]any, 0, len(array.Elements))

	for _, element := range array.Elements {
		value, err := convertValue(element)
		if err != nil {
			return nil, err
		}

		items = append(items, value)
	}

	return items, nil
}

func convertObjectKey(value hujson.Value) (string, error) {
	literal, ok := value.Value.(hujson.Literal)
	if !ok {
		return "", fmt.Errorf("unexpected hujson object key type: %T", value.Value)
	}

	if literal.Kind() != '"' {
		return "", fmt.Errorf("unexpected hujson object key kind: %q", literal.Kind())
	}

	return literal.String(), nil
}

func convertLiteral(literal hujson.Literal) (any, error) {
	switch literal.Kind() {
	case '{', '[':
		return nil, fmt.Errorf("unexpected composite literal kind: %q", literal.Kind())
	case 'n':
		return nil, nil
	case 'f', 't':
		return literal.Bool(), nil
	case '"':
		return literal.String(), nil
	case '0':
		return convertNumber(literal.String())
	default:
		return nil, fmt.Errorf("unsupported literal kind: %q", literal.Kind())
	}
}

func convertNumber(literal string) (any, error) {
	if strings.ContainsAny(literal, ".eE") {
		return json.Number(literal), nil
	}

	integer := new(big.Int)
	if _, ok := integer.SetString(literal, 10); !ok {
		return nil, fmt.Errorf("invalid integer literal: %q", literal)
	}

	return integer, nil
}
