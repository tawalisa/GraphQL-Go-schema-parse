package schemaParser

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	fileutil "graphql-go-schema-parser/util"
	"log"
	"os"
	"reflect"
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

func TestRef(t *testing.T) {
	//tt := reflect.StructOf([]reflect.StructField{
	//	{
	//		Name: "Height",
	//		Type: reflect.TypeOf(0.0),
	//	},
	//	{
	//		Name: "Name",
	//		Type: reflect.TypeOf(""),
	//	},
	//})
	//ref
	//fmt.Println(tt.NumField())
	//fmt.Println(tt.FieldByName("Name"))
	//fmt.Println(tt.Field(1))

	emptyStruct := reflect.StructOf([]reflect.StructField{
		{Name: "Field1", Type: reflect.TypeOf("")},
		{Name: "Field2", Type: reflect.TypeOf(0)},
	})

	fmt.Println(emptyStruct)
	instance := reflect.New(emptyStruct).Elem()
	// 设置结构体的字段值
	value1 := "Hello"
	value2 := 42

	instance.Field(0).SetString("11")
	instance.Field(0).Set(reflect.ValueOf(value1))
	instance.FieldByName("Field2").Set(reflect.ValueOf(value2))

	// 打印结构体的字段值
	fmt.Println(instance.Field(0).Interface()) // 输出: Hello
	//fmt.Println(reflect.ValueOf(instance).Field(1).Interface()) // 输出: 42
}

type testStruct struct {
}

func Test1(t *testing.T) {
	var ptr *testStruct // 声明一个int类型的指针变量
	ptrType := reflect.TypeOf(ptr).Elem()
	fmt.Println("Pointer type:", ptrType)
	//ptr := graphql.NewObject(graphql.ObjectConfig{
	//	Name: "Company"})

	// 使用reflect包的TypeOf函数获取指针变量的类型

}

func TestGraphql(t *testing.T) {
	root, _ := os.Getwd()
	println(root)
	sdl, e := fileutil.ReadFile("resource/schema/examplev1.sdl")
	require.NoError(t, e)
	assert.NotEqual(t, sdl, "")
	QueryType, _ := Parsing(string(sdl))
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
	})

	query := `
		{
  StockPrice(symbol: "AAOI", start_date: "2023-05-10", end_date: "2023-06-10") {
    AdjClose
    Close
    Date
    High
    Low
    Open
    Volume
  }
  company(symbol: "AAOI") {
    Name
    LongName
    Address1
    Address2
    City
    Zip
    Country
    Phone
    Fax
    Website
    CompanyOfficers {
      MaxAge
      Name
      Age
      YearBorn
      ExercisedValue
      UnexercisedValue
      FiscalYear
    }
  }
}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON) // {"data":{"hello":"world"}}
}
