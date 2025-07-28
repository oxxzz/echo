package hooks

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/labstack/echo/v4"
)

// SonicJSONSerializer implements JSON encoding using encoding/json.
type SonicJSONSerializer struct{}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
func (d SonicJSONSerializer) Serialize(c echo.Context, i any, indent string) error {
	enc := sonic.ConfigDefault.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d SonicJSONSerializer) Deserialize(c echo.Context, i any) error {
	err := sonic.ConfigDefault.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*decoder.MismatchTypeError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, field=%v, offset=%v", ute.Type, ute.Src, ute.Pos)).SetInternal(err)
	} else if se, ok := err.(*decoder.SyntaxError); ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Pos, se.Error())).SetInternal(err)
	}
	return err
}
