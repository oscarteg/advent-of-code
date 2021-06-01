package challenge

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/stateful"
)

var (
	passportLexer = lexer.Must(stateful.NewSimple([]stateful.Rule{
		{"DateTime", `\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\d(\.\d+)?(-\d\d:\d\d)?`, nil},
		{"Date", `\d\d\d\d-\d\d-\d\d`, nil},
		{"Time", `\d\d:\d\d:\d\d(\.\d+)?`, nil},
		{"Ident", `[a-zA-Z_][a-zA-Z_0-9]*`, nil},
		{"String", `"[^"]*"`, nil},
		{"Number", `[-+]?[.0-9]+\b`, nil},
		{"Punct", `\[|]|[-!()+/*=,]`, nil},
		{"comment", `#[^\n]+`, nil},
		{"whitespace", `\s+`, nil},
	}))
	tomlParser = participle.MustBuild(&Passport{},
		participle.Lexer(
			passportLexer,
		),
		participle.Unquote("String"),
	)

	cli struct {
		File string `help:"TOML file to parse." arg:""`
	}
	mandatoryFields = [7]string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}
)


type Passport struct {
	fields     map[string]string
	Properties []*Property `@@*`
}

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
}

type Property struct {
	Key   string `@Ident "="`
	Value *Value `@@`
}

type Value struct {
	String *string  `  @String`
	Number *float64 `| @Float`
}

func (p *Passport) NewPassport(fields map[string]string) *Passport {
	// Create a new passport
	return &Passport{fields: fields}
}

func (p *Passport) IsValid() bool {
	for _, field := range mandatoryFields {
		if _, ok := p.fields[field]; ok == false {
			return false
		}
		
	}
	return true

}

//func main() {
//	ctx := kong.Parse(&cli)
//	toml := &TOML{}
//	r, err := os.Open(cli.File)
//	ctx.FatalIfErrorf(err)
//	defer r.Close()
//	err = tomlParser.Parse(cli.File, r, toml)
//	ctx.FatalIfErrorf(err)
//	repr.Println(toml)
//}
