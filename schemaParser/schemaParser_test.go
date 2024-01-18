package schemaParser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	fileutil "graphql-go-schema-parser/util"
	"os"
	"regexp"
	"testing"
)

func TestParsing(t *testing.T) {

	root, _ := os.Getwd()
	println(root)
	sdl, e := fileutil.ReadFile("resource/schema/examplev1.sdl")
	require.NoError(t, e)
	assert.NotEqual(t, sdl, "")
	schemaType, err := Parsing(string(sdl))
	if err != nil {
		println(schemaType)
	}

}

func TestParsingFiled(t *testing.T) {
	str := "company(symbol: String!)"

	str = "company"
	//str = "aa"
	re := regexp.MustCompile(`^([^(]+)(\(([^)]+)\))?`)
	match := re.FindStringSubmatch(str)

	if len(match) > 1 {
		aa := match[1]
		bbb := match[2]
		fmt.Println("aa:", aa)
		fmt.Println("bbb:", bbb)
	} else {
		fmt.Println("No match found")
	}
}
