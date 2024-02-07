package cbor

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_ParseCBORToStruct_Success(t *testing.T) {
	t.Parallel()

	hexCBOR := `0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff000000`
	bytesCBOR, err := hexutil.Decode(hexCBOR)
	assert.NoError(t, err)

	parsed := struct {
		Url  string   `cbor:"url"`
		Path []string `cbor:"path"`
	}{}
	err = ParseDietCBORToStruct(bytesCBOR, &parsed)

	require.NoError(t, err)
	require.Equal(t, "https://etherprice.com/api", parsed.Url)
	require.Equal(t, []string{"recent", "usd"}, parsed.Path)
}

func Test_ParseCBORToStruct_WrongFieldType(t *testing.T) {
	t.Parallel()

	hexCBOR := `0xbf6375726c781a68747470733a2f2f657468657270726963652e636f6d2f61706964706174689f66726563656e7463757364ffff000000`
	bytesCBOR, err := hexutil.Decode(hexCBOR)
	assert.NoError(t, err)

	parsed := struct {
		Url  string `cbor:"url"`
		Path []int  `cbor:"path"` // exect int but get string
	}{}
	err = ParseDietCBORToStruct(bytesCBOR, &parsed)

	require.Error(t, err)
}

func Test_ParseCBORToStruct_BinaryStringOfWrongType(t *testing.T) {
	t.Parallel()

	// {"key":"value"} but with last byte replaced with invalid unicode (0x88)
	hexCBOR := `0x636B65796576616C7588`
	bytesCBOR, err := hexutil.Decode(hexCBOR)
	assert.NoError(t, err)

	parsed := struct {
		Key string `cbor:"key"`
	}{}
	err = ParseDietCBORToStruct(bytesCBOR, &parsed)
	require.Error(t, err)
}

func Test_autoAddMapDelimiters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []byte
		want []byte
	}{
		{
			"map(0)",
			hexutil.MustDecode("0xA0"),
			hexutil.MustDecode("0xA0"),
		},
		{
			`map(1) {"key":"value"}`,
			hexutil.MustDecode("0xA1636B65796576616C7565"),
			hexutil.MustDecode("0xA1636B65796576616C7565"),
		},
		{
			"array(0)",
			hexutil.MustDecode("0x80"),
			hexutil.MustDecode("0x80"),
		},
		{
			`map(*) {"key":"value"}`,
			hexutil.MustDecode("0xbf636B65796576616C7565ff"),
			hexutil.MustDecode("0xbf636B65796576616C7565ff"),
		},
		{
			`map(*) {"key":"value"} missing open delimiter`,
			hexutil.MustDecode("0x636B65796576616C7565ff"),
			hexutil.MustDecode("0xbf636B65796576616C7565ffff"),
		},
		{
			`map(*) {"key":"value"} missing closing delimiter`,
			hexutil.MustDecode("0xbf636B65796576616C7565"),
			hexutil.MustDecode("0xbf636B65796576616C7565"),
		},
		{
			`map(*) {"key":"value"} missing both delimiters`,
			hexutil.MustDecode("0x636B65796576616C7565"),
			hexutil.MustDecode("0xbf636B65796576616C7565ff"),
		},
		{
			"empty input adds delimiters",
			[]byte{},
			[]byte{0xbf, 0xff},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, autoAddMapDelimiters(test.in))
		})
	}
}

func jsonMustUnmarshal(t *testing.T, in string) interface{} {
	var j interface{}
	err := json.Unmarshal([]byte(in), &j)
	require.NoError(t, err)
	return j
}

func TestCoerceInterfaceMapToStringMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{"empty map", map[interface{}]interface{}{}, map[string]interface{}{}},
		{"simple map", map[interface{}]interface{}{"key": "value"}, map[string]interface{}{"key": "value"}},
		{"int map", map[int]interface{}{1: "value"}, map[int]interface{}{1: "value"}},
		{
			"nested string map map",
			map[string]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}},
			map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}},
		},
		{
			"nested map map",
			map[interface{}]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}},
			map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}},
		},
		{
			"nested map array",
			map[interface{}]interface{}{"key": []interface{}{1, "value"}},
			map[string]interface{}{"key": []interface{}{1, "value"}},
		},
		{"empty array", []interface{}{}, []interface{}{}},
		{"simple array", []interface{}{1, "value"}, []interface{}{1, "value"}},
		{
			"nested array map",
			[]interface{}{map[interface{}]interface{}{"key": map[interface{}]interface{}{"nk": "nv"}}},
			[]interface{}{map[string]interface{}{"key": map[string]interface{}{"nk": "nv"}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoded, err := CoerceInterfaceMapToStringMap(test.input)
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(test.want, decoded))
		})
	}
}

func TestCoerceInterfaceMapToStringMap_BadInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input interface{}
	}{
		{"error map", map[interface{}]interface{}{1: "value"}},
		{"error array", []interface{}{map[interface{}]interface{}{1: "value"}}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := CoerceInterfaceMapToStringMap(test.input)
			assert.Error(t, err)
		})
	}
}

func TestJSON_CBOR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
	}{
		{"empty object", jsonMustUnmarshal(t, `{}`)},
		{"array", jsonMustUnmarshal(t, `[1,2,3,4]`)},
		{
			"basic object",
			jsonMustUnmarshal(t, `{"path":["recent","usd"],"url":"https://etherprice.com/api"}`),
		},
		{
			"complex object",
			jsonMustUnmarshal(t, `{"a":{"1":[{"b":"free"},{"c":"more"},{"d":["less", {"nesting":{"4":"life"}}]}]}}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encoded := mustMarshal(t, test.in)

			var decoded interface{}
			err := cbor.Unmarshal(encoded, &decoded)
			require.NoError(t, err)

			decoded, err = CoerceInterfaceMapToStringMap(decoded)
			require.NoError(t, err)
			assert.True(t, reflect.DeepEqual(test.in, decoded))
		})
	}
}

// mustMarshal returns a bytes array of the JSON map or array encoded to CBOR.
func mustMarshal(t *testing.T, j interface{}) []byte {
	switch v := j.(type) {
	case map[string]interface{}, []interface{}, nil:
		b, err := cbor.Marshal(v)
		if err != nil {
			t.Fatalf("failed to marshal CBOR: %v", err)
		}
		return b
	default:
		t.Fatalf("unable to coerce JSON to CBOR for type %T", v)
		return nil
	}
}
