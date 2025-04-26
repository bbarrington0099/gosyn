package main

import (
	"strings"
	"fmt"
	"log"
	"os"
)

const (
    Reset  = "\033[0m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
	Cyan   = "\033[36m"
    
    BoldRed    = "\033[1;31m"
    BoldGreen  = "\033[1;32m"
    BoldYellow = "\033[1;33m"
	BoldPurple = "\033[1;35m"
	BoldCyan   = "\033[1;36m"

	BoldUnderline = "\033[1;4m"

	Italic = "\033[3m"
	BoldItalic = "\033[1;3m"
)

type command struct {
	action string
	args []string
}

type section struct {
	name string
	subsections []subsection
}

type subsection struct {
	name string
	content string
}

func parseCommand() (command, error) {
	var err error = nil
	var cmd command
	if len(os.Args) < 2 {
		err = fmt.Errorf("%sERROR%s parseCommand(): no command provided", BoldRed, Reset)
		return cmd, err
	}
	if len(os.Args) < 3 {
		cmd = command{
			action: os.Args[1],
			args: []string{""},
		}
	} else {
		cmd = command{
			action: os.Args[1],
			args: os.Args[2:],
		}
	}
	return cmd, err
}

func executeCommand() (string, error) {
	var err error = nil
	cmd, parseErr := parseCommand()
	if parseErr != nil {
		return "", parseErr	
	}
	sections := initializeSections()

	switch strings.ToLower(cmd.action) {
	case "h":
		fallthrough
	case "help":
		return listActions(), nil
	case "lsec":
		fallthrough
	case "listSections":
		return listSections(sections), nil
	case "lsub":
		fallthrough
	case "listSubsections":
		if cmd.args[0] == "" {
			err = fmt.Errorf("%sERROR%s executeCommand(): no section name provided for listSubsections <sectionName>", BoldRed, Reset)
			return "", err
		}
		return listSubsections(sections, cmd.args[0])
	default:
		if len(cmd.args) < 1 {
			err = fmt.Errorf("%sERROR%s executeCommand(): no subsection name provided for tax <sectionName> <subsectionName> in args \"%v\"", BoldRed, Reset, os.Args[1:])
			return "", err
		}
		return tax(sections, cmd.action, cmd.args[0])
	}
}

func listActions() string {
	return fmt.Sprintf(("%sAvailable commands%s:\n" +
		" - %s(help | h)%s: List all available commands\n" +
		" - %s(listSections | lsec)%s: List all sections\n" +
		" - %s(listSubsections | lsub) <sectionName>%s: List all subsections in a section\n" +
		"    - %s<sectionName>%s is the name of the section to list subsections for\n" +
		" - %s<sectionName> <subsectionName>%s: Get syntax information for a subsection\n" +
		"    - %s<sectionName>%s is the name of the section\n" +
		"    - %s<subsectionName>%s is the name of the subsection\n"),
		BoldUnderline, Reset, // Available commands
		BoldYellow, Reset, // help
		BoldCyan, Reset, // listSections
		BoldCyan, Reset, // listSubsections
		Italic, Reset, // > sectionName
		BoldGreen, Reset, // tax
		Italic, Reset, // > sectionName
		Italic, Reset, // > subsectionName
	)
}

func listSections(sections []section) (string) {
	output := fmt.Sprintf("%sSections%s:\n", BoldItalic, Reset)
	for _, sec := range sections {
		output += fmt.Sprintf(" - %s%s%s\n", BoldUnderline, sec.name, Reset)
		for _, sub := range sec.subsections {
			output += fmt.Sprintf("   - %s%s%s\n", Yellow, sub.name, Reset)
		}
	}
	return output
}

func listSubsections(sections []section, sectionName string) (string, error) {
	var err error = nil
	output := fmt.Sprintf("%sSubsections%s in %s%s%s:\n", 
	BoldYellow, Reset, // Subsections
	BoldGreen, sectionName, Reset) // sectionName
	for _, sec := range sections {
		if sec.name == sectionName {
			for _, sub := range sec.subsections {
				output += fmt.Sprintf("   - %s\n", sub.name)
			}
			return output, err
		}
	}
	err = fmt.Errorf("%sERROR%s listSubsections(): section \"%s\" not found", BoldRed, Reset, sectionName)
	return "", err
}

func tax(sections []section, sectionName string, subsectionName string) (string, error) {
	var err error = nil
	for _, sec := range sections {
		if strings.EqualFold(sec.name, sectionName) {
			for _, sub := range sec.subsections {
				if strings.EqualFold(sub.name, subsectionName) {
					return fmt.Sprintf("%sSyntax information%s for %s%s%s in %s%s%s: %s\n", 
					BoldPurple, Reset, // Syntax information
					Yellow, subsectionName, Reset, // subsectionName
					Green, sectionName, Reset, // sectionName
					sub.content), err
				}
			}
			err = fmt.Errorf("%sERROR%s tax(): subsection \"%s\" not found in section \"%s\"", BoldRed, Reset, subsectionName, sectionName)
			return "", err
		}
	}
	err = fmt.Errorf("%sERROR%s tax(): section \"%s\" not found", BoldRed, Reset, sectionName)
	return "", err
}

func main() {
	output, commandError := executeCommand()
	if commandError != nil {
		log.Fatal(commandError)
	}
	fmt.Println(output)
}

func initializeSections() ([]section) {
	sections := []section{
		{
			name: "Variables",
			subsections: []subsection{
				{name: "Declaration", content: fmt.Sprintf(
					("%sVariable Declaration%s:\n\n" +
				    "\t%svar%s %s<variableName> <type>%s = %s<value>%s\n" +
					"\t%vconst%s %s<variableName> <type>%s = %s<value>%s\n\n" +
					"\t%s<varaibleName>%s := %s<value>%s\n" +
					"\t\t - The %s:=%s (Walrus) operator is a shorthand for declaring and initializing a variable using type inference\n"),
					BoldItalic, Reset, // Variable Declaration
					Cyan, Reset, // var
					Yellow, Reset, // <variableName> <type>
					Green, Reset, // <value>
					Cyan, Reset, // const
					Yellow, Reset, // <variableName> <type>
					Green, Reset, // <value>
					Cyan, Reset, // <variableName>
					Green, Reset, // <value>
					BoldPurple, Reset, // :=
				)},
				{name: "Types", content: fmt.Sprintf(
					("%sBasic Types%s:\n\n" +
					"\t%sbool%s\n" + 
					"\t%sstring%s\n" +
					"\t%sint%s, %sint8%s, %sint16%s, %sint32%s, %sint64%s\n" +
					"\t%suint%s, %suint8%s, %suint16%s, %suint32%s, %suint64%s\n" +
					"\t%sfloat32%s, %sfloat64%s\n" +
					"\t%scomplex64%s, %scomplex128%s\n" +
					"\t%sbyte%s (alias for uint8)\n" +
					"\t%srune%s (alias for int32, represents Unicode code point)\n"),
					BoldItalic, Reset, // Basic Types
					Yellow, Reset, // bool
					Yellow, Reset, // string
					Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset, // int types
					Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset, // uint types
					Yellow, Reset, Yellow, Reset, // float types
					Yellow, Reset, Yellow, Reset, // complex types
					Yellow, Reset, // byte
					Yellow, Reset, // rune
				)},
			},
		},
		{
			name: "Conditionals",
			subsections: []subsection{
				{name: "If", content: fmt.Sprintf(
					("%sIf Statement%s:\n\n" +
					"\t%sif%s %s<condition>%s {\n" +
					"\t\t// code\n" +
					"\t}\n\n" +
					"\t// With initialization statement\n" +
					"\t%sif%s %s<initialization>%s; %s<condition>%s {\n" +
					"\t\t// code\n" +
					"\t}\n"),
					BoldItalic, Reset, // If Statement
					Cyan, Reset, // if
					BoldPurple, Reset, // <condition>
					Cyan, Reset, // if
					Green, Reset, // <initialization>
					BoldPurple, Reset, // <condition>
				)},
				{name: "IfElse", content: fmt.Sprintf(
					("%sIf-Else Statement%s:\n\n" +
					"\t%sif%s %s<condition>%s {\n" +
					"\t\t// code if condition is true\n" +
					"\t} %selse%s {\n" +
					"\t\t// code if condition is false\n" +
					"\t}\n"),
					BoldItalic, Reset, // If-Else Statement
					Cyan, Reset, // if
					BoldPurple, Reset, // <condition>
					Cyan, Reset, // else
				)},
				{name: "ElseIf", content: fmt.Sprintf(
					("%sElse-If Statement%s:\n\n" +
					"\t%sif%s %s<condition1>%s {\n" +
					"\t\t// code if condition1 is true\n" +
					"\t} %selse if%s %s<condition2>%s {\n" +
					"\t\t// code if condition2 is true\n" +
					"\t} %selse%s {\n" +
					"\t\t// code if all conditions are false\n" +
					"\t}\n"),
					BoldItalic, Reset, // Else-If Statement
					Cyan, Reset, // if
					BoldPurple, Reset, // <condition1>
					Cyan, Reset, // else if
					BoldPurple, Reset, // <condition2>
					Cyan, Reset, // else
				)},
				{name: "Switch", content: fmt.Sprintf(
					("%sSwitch Statement%s:\n\n" +
					"\t%sswitch%s %s<expression>%s {\n" +
					"\t%scase%s %s<value1>%s:\n" +
					"\t\t// code if expression == value1\n" +
					"\t%scase%s %s<value2>%s, %s<value3>%s:\n" +
					"\t\t// code if expression == value2 or expression == value3\n" +
					"\t%sdefault%s:\n" +
					"\t\t// code if no case matches\n" +
					"\t}\n\n" + 
					"\t// With initialization\n" +
					"\t%sswitch%s %s<initialization>%s; %s<expression>%s {\n" +
					"\t\t// cases\n" +
					"\t}\n"),
					BoldItalic, Reset, // Switch Statement
					Cyan, Reset, // switch
					Yellow, Reset, // <expression>
					Cyan, Reset, // case
					Green, Reset, // <value1>
					Cyan, Reset, // case
					Green, Reset, Green, Reset, // <value2>, <value3>
					Cyan, Reset, // default
					Cyan, Reset, // switch
					Green, Reset, // <initialization>
					Yellow, Reset, // <expression>
				)},
				{name: "TypeSwitch", content: fmt.Sprintf(
					("%sType Switch%s:\n\n" +
					"\t%sswitch%s %s<variable>%s := %s<expression>%s.(type) {\n" +
					"\t%scase%s %sint%s:\n" +
					"\t\t// code if expression is an int\n" +
					"\t%scase%s %sstring%s:\n" +
					"\t\t// code if expression is a string\n" +
					"\t%scase%s %snil%s:\n" +
					"\t\t// code if expression is nil\n" +
					"\t%sdefault%s:\n" +
					"\t\t// code for any other type\n" +
					"\t}\n"),
					BoldItalic, Reset, // Type Switch
					Cyan, Reset, // switch
					Yellow, Reset, // <variable>
					Yellow, Reset, // <expression>
					Cyan, Reset, // case
					Yellow, Reset, // int
					Cyan, Reset, // case
					Yellow, Reset, // string
					Cyan, Reset, // case
					Yellow, Reset, // nil
					Cyan, Reset, // default
				)},
			},
		},
		{
			name: "Loops",
			subsections: []subsection{
				{name: "For", content: fmt.Sprintf(
					("%sFor Loop (Standard)%s:\n\n" +
					"\t%sfor%s %s<initialization>%s; %s<condition>%s; %s<post>%s {\n" +
					"\t\t// code\n" +
					"\t}\n\n" +
					"\t// Examples:\n" +
					"\t%sfor%s %si%s := 0; i < 10; i++ {\n" +
					"\t\t// code\n" +
					"\t}\n"),
					BoldItalic, Reset, // For Loop
					Cyan, Reset, // for
					Green, Reset, // <initialization>
					BoldPurple, Reset, // <condition>
					Green, Reset, // <post>
					Cyan, Reset, // for
					Yellow, Reset, // i
				)},
				{name: "WhileStyle", content: fmt.Sprintf(
					("%sFor Loop (While Style)%s:\n\n" +
					"\t%sfor%s %s<condition>%s {\n" +
					"\t\t// code\n" +
					"\t}\n\n" +
					"\t// Example:\n" +
					"\t%sfor%s %si < 10%s {\n" +
					"\t\t// code\n" +
					"\t\ti++\n" +
					"\t}\n"),
					BoldItalic, Reset, // While Style
					Cyan, Reset, // for
					BoldPurple, Reset, // <condition>
					Cyan, Reset, // for
					BoldPurple, Reset, // i < 10
				)},
				{name: "Infinite", content: fmt.Sprintf(
					("%sInfinite Loop%s:\n\n" +
					"\t%sfor%s {\n" +
					"\t\t// code runs indefinitely\n" +
					"\t\t%sif%s %s<condition>%s {\n" +
					"\t\t\t%sbreak%s\n" +
					"\t\t}\n" +
					"\t}\n"),
					BoldItalic, Reset, // Infinite Loop
					Cyan, Reset, // for
					Cyan, Reset, // if
					BoldPurple, Reset, // <condition>
					BoldYellow, Reset, // break
				)},
				{name: "Range", content: fmt.Sprintf(
					("%sFor Range Loop%s:\n\n" +
					"\t%sfor%s %s<index>%s, %s<value>%s := %srange%s %s<collection>%s {\n" +
					"\t\t// code\n" +
					"\t}\n\n" +
					"\t// Ignore index with _\n" +
					"\t%sfor%s %s_%s, %s<value>%s := %srange%s %s<collection>%s {\n" +
					"\t\t// code\n" +
					"\t}\n\n" +
					"\t// Ignore value\n" +
					"\t%sfor%s %s<index>%s, %s_%s := %srange%s %s<collection>%s {\n" +
					"\t\t// code\n" +
					"\t}\n\n" +
					"\t// Only index\n" +
					"\t%sfor%s %s<index>%s := %srange%s %s<collection>%s {\n" +
					"\t\t// code\n" +
					"\t}\n"),
					BoldItalic, Reset, // For Range Loop
					Cyan, Reset, // for
					Yellow, Reset, // <index>
					Yellow, Reset, // <value>
					Cyan, Reset, // range
					Green, Reset, // <collection>
					Cyan, Reset, // for
					Yellow, Reset, // _
					Yellow, Reset, // <value>
					Cyan, Reset, // range
					Green, Reset, // <collection>
					Cyan, Reset, // for
					Yellow, Reset, // <index>
					Yellow, Reset, // _
					Cyan, Reset, // range
					Green, Reset, // <collection>
					Cyan, Reset, // for
					Yellow, Reset, // <index>
					Cyan, Reset, // range
					Green, Reset, // <collection>
				)},
				{name: "ControlFlow", content: fmt.Sprintf(
					("%sLoop Control Flow%s:\n\n" +
					"\t%sbreak%s - Exit the loop immediately\n" +
					"\t%scontinue%s - Skip the current iteration and move to the next one\n" +
					"\t%sbreak <label>%s - Break out of the labeled loop (for nested loops)\n" +
					"\t%scontinue <label>%s - Continue to next iteration of labeled loop\n\n" +
					"\t// Example with label:\n" +
					"\t%sOuterLoop%s:\n" +
					"\t%sfor%s i := 0; i < 5; i++ {\n" +
					"\t\t%sfor%s j := 0; j < 5; j++ {\n" +
					"\t\t\t%sif%s j == 3 {\n" +
					"\t\t\t\t%sbreak OuterLoop%s\n" +
					"\t\t\t}\n" +
					"\t\t}\n" +
					"\t}\n"),
					BoldItalic, Reset, // Loop Control Flow
					BoldYellow, Reset, // break
					BoldYellow, Reset, // continue
					BoldYellow, Reset, // break <label>
					BoldYellow, Reset, // continue <label>
					BoldCyan, Reset, // OuterLoop:
					Cyan, Reset, // for
					Cyan, Reset, // for
					Cyan, Reset, // if
					BoldYellow, Reset, // break OuterLoop
				)},
			},
		},
		{
			name: "Functions",
			subsections: []subsection{
				{name: "Declaration", content: fmt.Sprintf(
					"%sFunction Declaration%s:\n\n"+
					"\tfunc %sadd%s(%sa%s, %sb%s int) int {\n"+
					"\t\treturn %sa%s + %sb%s\n"+
					"\t}\n\n"+
					"\tfunc %sswap%s(%sa%s, %sb%s string) (%sstring%s, %sstring%s) {\n"+
					"\t\treturn %sb%s, %sa%s\n"+
					"\t}",
					BoldItalic, Reset,
					Cyan, Reset, Yellow, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset,
					Cyan, Reset, Yellow, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset,
				)},
				{name: "Variadic", content: fmt.Sprintf(
					"%sVariadic Functions%s:\n\n"+
					"\tfunc %ssum%s(%snums%s ...int) int {\n"+
					"\t\ttotal := 0\n"+
					"\t\tfor _, %sv%s := range %snums%s {\n"+
					"\t\t\ttotal += %sv%s\n"+
					"\t\t}\n"+
					"\t\treturn total\n"+
					"\t}",
					BoldItalic, Reset,
					Cyan, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset,
					Yellow, Reset,
				)},
				{name: "Closures", content: fmt.Sprintf(
					"%sClosures%s:\n\n"+
					"\tfunc %sintSeq%s() func() int {\n"+
					"\t\t%si%s := 0\n"+
					"\t\treturn func() int {\n"+
					"\t\t\t%si%s++\n"+
					"\t\t\treturn %si%s\n"+
					"\t\t}\n"+
					"\t}",
					BoldItalic, Reset,
					Cyan, Reset,
					Yellow, Reset,
					Yellow, Reset,
					Yellow, Reset,
				)},
			},
		},
		{
			name: "DataStructures",
			subsections: []subsection{
				{name: "Slices", content: fmt.Sprintf(
					("%sSlices%s:\n\n" +
					"\t// %sDeclaration%s\n" +
					"\t%svar%s %s<name>%s []%s<type>%s\n" +
					"\t%s<name>%s := []%s<type>%s{%s<values>%s}\n" +
					"\t%s<name>%s := %smake%s([]%s<type>%s, %s<length>%s, %s<capacity>%s)\n\n" +
					"\t// %sAccessing Elements%s\n" +
					"\t%selement%s := %s<slice>%s[%s<index>%s]\n\n" +
					"\t// %sSlicing%s\n" +
					"\t%s<newSlice>%s := %s<slice>%s[%s<start>%s:%s<end>%s]\n" +
					"\t%s<newSlice>%s := %s<slice>%s[%s<start>%s:%s<end>%s:%s<capacity>%s]\n\n" +
					"\t// %sOperations%s\n" +
					"\t%s<slice>%s = %sappend%s(%s<slice>%s, %s<element1>%s, %s<element2>%s)\n" +
					"\t%s<length>%s := %slen%s(%s<slice>%s)\n" +
					"\t%s<capacity>%s := %scap%s(%s<slice>%s)\n" +
					"\t%scopy%s(%s<dst>%s, %s<src>%s)\n\n" +
					"\t// %sIteration%s\n" +
					"\t%sfor%s %s<index>%s, %s<value>%s := %srange%s %s<slice>%s {\n" +
					"\t\t// code\n" +
					"\t}\n"),
					BoldItalic, Reset, // Slices
					BoldUnderline, Reset, // Declaration
					Cyan, Reset, // var
					Yellow, Reset, // <name>
					Yellow, Reset, // <type>
					Yellow, Reset, // <name>
					Yellow, Reset, // <type>
					Green, Reset, // <values>
					Yellow, Reset, // <name>
					Cyan, Reset, // make
					Yellow, Reset, // <type>
					Green, Reset, // <length>
					Green, Reset, // <capacity>
					BoldUnderline, Reset, // Accessing Elements
					Yellow, Reset, // element
					Yellow, Reset, // <slice>
					Green, Reset, // <index>
					BoldUnderline, Reset, // Slicing
					Yellow, Reset, // <newSlice>
					Yellow, Reset, // <slice>
					Green, Reset, // <start>
					Green, Reset, // <end>
					Yellow, Reset, // <newSlice>
					Yellow, Reset, // <slice>
					Green, Reset, // <start>
					Green, Reset, // <end>
					Green, Reset, // <capacity>
					BoldUnderline, Reset, // Operations
					Yellow, Reset, // <slice>
					Cyan, Reset, // append
					Yellow, Reset, // <slice>
					Green, Reset, // <element1>
					Green, Reset, // <element2>
					Yellow, Reset, // <length>
					Cyan, Reset, // len
					Yellow, Reset, // <slice>
					Yellow, Reset, // <capacity>
					Cyan, Reset, // cap
					Yellow, Reset, // <slice>
					Cyan, Reset, // copy
					Yellow, Reset, // <dst>
					Yellow, Reset, // <src>
					BoldUnderline, Reset, // Iteration
					Cyan, Reset, // for
					Yellow, Reset, // <index>
					Yellow, Reset, // <value>
					Cyan, Reset, // range
					Yellow, Reset, // <slice>
				)},
				{name: "Maps", content: fmt.Sprintf(
					("%sMaps%s:\n\n" +
					"\t// %sDeclaration%s\n" +
					"\t%svar%s %s<name>%s map[%s<keyType>%s]%s<valueType>%s\n" +
					"\t%s<name>%s := map[%s<keyType>%s]%s<valueType>%s{%s<key1>%s: %s<value1>%s, %s<key2>%s: %s<value2>%s}\n" +
					"\t%s<name>%s := %smake%s(map[%s<keyType>%s]%s<valueType>%s, %s<capacity>%s)\n\n" +
					"\t// %sAccessing Elements%s\n" +
					"\t%s<value>%s := %s<map>%s[%s<key>%s]\n" +
					"\t%s<value>%s, %s<exists>%s := %s<map>%s[%s<key>%s] // Check if key exists\n\n" +
					"\t// %sModifying%s\n" +
					"\t%s<map>%s[%s<key>%s] = %s<value>%s // Add or update\n" +
					"\t%sdelete%s(%s<map>%s, %s<key>%s) // Remove\n\n" +
					"\t// %sOperations%s\n" +
					"\t%s<length>%s := %slen%s(%s<map>%s)\n\n" +
					"\t// %sIteration%s\n" +
					"\t%sfor%s %s<key>%s, %s<value>%s := %srange%s %s<map>%s {\n" +
					"\t\t// code\n" +
					"\t}\n"),
					BoldItalic, Reset, // Maps
					BoldUnderline, Reset, // Declaration
					Cyan, Reset, // var
					Yellow, Reset, // <name>
					Yellow, Reset, // <keyType>
					Yellow, Reset, // <valueType>
					Yellow, Reset, // <name>
					Yellow, Reset, // <keyType>
					Yellow, Reset, // <valueType>
					Green, Reset, // <key1>
					Green, Reset, // <value1>
					Green, Reset, // <key2>
					Green, Reset, // <value2>
					Yellow, Reset, // <name>
					Cyan, Reset, // make
					Yellow, Reset, // <keyType>
					Yellow, Reset, // <valueType>
					Green, Reset, // <capacity>
					BoldUnderline, Reset, // Accessing Elements
					Yellow, Reset, // <value>
					Yellow, Reset, // <map>
					Green, Reset, // <key>
					Yellow, Reset, // <value>
					Yellow, Reset, // <exists>
					Yellow, Reset, // <map>
					Green, Reset, // <key>
					BoldUnderline, Reset, // Modifying
					Yellow, Reset, // <map>
					Green, Reset, // <key>
					Green, Reset, // <value>
					Cyan, Reset, // delete
					Yellow, Reset, // <map>
					Green, Reset, // <key>
					BoldUnderline, Reset, // Operations
					Yellow, Reset, // <length>
					Cyan, Reset, // len
					Yellow, Reset, // <map>
					BoldUnderline, Reset, // Iteration
					Cyan, Reset, // for
					Yellow, Reset, // <key>
					Yellow, Reset, // <value>
					Cyan, Reset, // range
					Yellow, Reset, // <map>
				)},
				{name: "Structs", content: fmt.Sprintf(
					("%sStructs%s:\n\n" +
					"\t// %sDefinition%s\n" +
					"\t%stype%s %s<Name>%s struct {\n" +
					"\t\t%s<field1>%s %s<type1>%s\n" +
					"\t\t%s<field2>%s %s<type2>%s\n" +
					"\t\t%s<field3>%s %s<type3>%s `%s<tag>%s`\n" +
					"\t}\n\n" +
					"\t// %sCreation%s\n" +
					"\t%s<var1>%s := %s<Name>%s{%s<field1>%s: %s<value1>%s, %s<field2>%s: %s<value2>%s}\n" +
					"\t%s<var2>%s := %s<Name>%s{%s<value1>%s, %s<value2>%s} // Positional initialization\n" +
					"\t%s<var3>%s := new(%s<Name>%s) // Zero-initialized\n\n" +
					"\t// %sAccessing Fields%s\n" +
					"\t%s<value>%s := %s<struct>%s.%s<field>%s\n" +
					"\t%s<struct>%s.%s<field>%s = %s<newValue>%s\n\n" +
					"\t// %sPointer to Struct%s\n" +
					"\t%s<ptr>%s := &%s<struct>%s\n" +
					"\t%s<value>%s := %s<ptr>%s.%s<field>%s // Automatic dereferencing\n\n" +
					"\t// %sMethods of Structs%s\n" +
					"\t%sfunc%s (%s<receiver>%s %s<Name>%s) %s<methodName>%s(%s<param>%s %s<type>%s) %s<returnType>%s {\n" +
					"\t\t// code\n" +
					"\t}\n" +
					"\t%s<receiver>%s.%s<methodName>%s(%s<param>%s)\n"),
					BoldItalic, Reset, // Structs
					BoldUnderline, Reset, // Definition
					Cyan, Reset, // type
					Yellow, Reset, // <Name>
					Yellow, Reset, // <field1>
					Yellow, Reset, // <type1>
					Yellow, Reset, // <field2>
					Yellow, Reset, // <type2>
					Yellow, Reset, // <field3>
					Yellow, Reset, // <type3>
					Italic, Reset, // <tag>
					BoldUnderline, Reset, // Creation
					Yellow, Reset, // <var1>
					Yellow, Reset, // <Name>
					Yellow, Reset, // <field1>
					Green, Reset, // <value1>
					Yellow, Reset, // <field2>
					Green, Reset, // <value2>
					Yellow, Reset, // <var2>
					Yellow, Reset, // <Name>
					Green, Reset, // <value1>
					Green, Reset, // <value2>
					Yellow, Reset, // <var3>
					Yellow, Reset, // <Name>
					BoldUnderline, Reset, // Accessing Fields
					Yellow, Reset, // <value>
					Yellow, Reset, // <struct>
					Yellow, Reset, // <field>
					Yellow, Reset, // <struct>
					Yellow, Reset, // <field>
					Green, Reset, // <newValue>
					BoldUnderline, Reset, // Pointer to Struct
					Yellow, Reset, // <ptr>
					Yellow, Reset, // <struct>
					Yellow, Reset, // <value>
					Yellow, Reset, // <ptr>
					Yellow, Reset, // <field>
					BoldUnderline, Reset, // Methods of Structs
					Cyan, Reset, // func
					Yellow, Reset, // <receiver>
					Yellow, Reset, // <Name>
					Yellow, Reset, // <methodName>
					Yellow, Reset, // <param>
					Yellow, Reset, // <type>
					Yellow, Reset, // <returnType>
					Cyan, Reset, // <receiver>
					Yellow, Reset, // <methodName>
					Yellow, Reset, // <param>
				)},
				{name: "Interfaces", content: fmt.Sprintf(
					("%sInterfaces%s:\n\n" +
					"\t// %sDefinition%s\n" +
					"\t%stype%s %s<Name>%s interface {\n" +
					"\t\t%s<Method1>%s(%s<param1>%s %s<type1>%s) %s<returnType1>%s\n" +
					"\t\t%s<Method2>%s(%s<param2>%s %s<type2>%s, %s<param3>%s %s<type3>%s) (%s<returnType2>%s, %s<returnType3>%s)\n" +
					"\t}\n\n" +
					"\t// %sEmpty Interface%s\n" +
					"\t%svar%s %s<anything>%s interface{}\n\n" +
					"\t// %sImplementation%s (implicit, no \"implements\" keyword)\n" +
					"\t%stype%s %s<StructName>%s struct {\n" +
					"\t\t// fields\n" +
					"\t}\n\n" +
					"\t%sfunc%s (%s<receiver>%s %s<StructName>%s) %s<Method1>%s(%s<param1>%s %s<type1>%s) %s<returnType1>%s {\n" +
					"\t\t// implementation\n" +
					"\t}\n\n" +
					"\t// %sType Assertion%s\n" +
					"\t%s<value>%s := %s<interfaceVar>%s.(%s<Type>%s) // Panics if wrong type\n" +
					"\t%s<value>%s, %s<ok>%s := %s<interfaceVar>%s.(%s<Type>%s) // Safe checking\n\n" +
					"\t// %sType Switch%s\n" +
					"\t%sswitch%s %s<v>%s := %s<interfaceVar>%s.(type) {\n" +
					"\t%scase%s %s<Type1>%s:\n" +
					"\t\t// v is Type1\n" +
					"\t%scase%s %s<Type2>%s:\n" +
					"\t\t// v is Type2\n" +
					"\t%sdefault%s:\n" +
					"\t\t// unknown type\n" +
					"\t}\n"),
					BoldItalic, Reset, // Interfaces
					BoldUnderline, Reset, // Definition
					Cyan, Reset, // type
					Yellow, Reset, // <Name>
					Yellow, Reset, // <Method1>
					Yellow, Reset, // <param1>
					Yellow, Reset, // <type1>
					Yellow, Reset, // <returnType1>
					Yellow, Reset, // <Method2>
					Yellow, Reset, // <param2>
					Yellow, Reset, // <type2>
					Yellow, Reset, // <param3>
					Yellow, Reset, // <type3>
					Yellow, Reset, // <returnType2>
					Yellow, Reset, // <returnType3>
					BoldUnderline, Reset, // Empty Interface
					Cyan, Reset, // var
					Yellow, Reset, // <anything>
					BoldUnderline, Reset, // Implementation
					Cyan, Reset, // type
					Yellow, Reset, // <StructName>
					Cyan, Reset, // func
					Yellow, Reset, // <receiver>
					Yellow, Reset, // <StructName>
					Yellow, Reset, // <Method1>
					Yellow, Reset, // <param1>
					Yellow, Reset, // <type1>
					Yellow, Reset, // <returnType1>
					BoldUnderline, Reset, // Type Assertion
					Yellow, Reset, // <value>
					Yellow, Reset, // <interfaceVar>
					Yellow, Reset, // <Type>
					Yellow, Reset, // <value>
					Yellow, Reset, // <ok>
					Yellow, Reset, // <interfaceVar>
					Yellow, Reset, // <Type>
					BoldUnderline, Reset, // Type Switch
					Cyan, Reset, // switch
					Yellow, Reset, // <v>
					Yellow, Reset, // <interfaceVar>
					Cyan, Reset, // case
					Yellow, Reset, // <Type1>
					Cyan, Reset, // case
					Yellow, Reset, // <Type2>
					Cyan, Reset, // default
				)},
			},
		},
		{
			name: "Channels",
			subsections: []subsection{
				{name: "Buffered", content: fmt.Sprintf(
					("%sBuffered Channels%s:\n\n"+
						"\t%sch%s := %smake%s(chan %sint%s, %s3%s)\n"+
						"\t%sch%s <- %s1%s  %s// Non-blocking until buffer full%s\n"+
						"\t%sch%s <- %s2%s\n"+
						"\t%sclose%s(%sch%s)\n"+
						"\t%sfor%s %sv%s := %srange%s %sch%s {\n"+
						"\t\t%sfmt.Println%s(%sv%s)\n"+
						"\t}\n"),
					BoldItalic, Reset, // Buffered Channels
					Yellow, Reset, Cyan, Reset, Yellow, Reset, Green, Reset, // ch make int 3
					Yellow, Reset, Green, Reset, Cyan, Reset, // ch 1 //...full
					Yellow, Reset, Green, Reset, // ch 2
					Cyan, Reset, Yellow, Reset, // close ch
					Cyan, Reset, Yellow, Reset, Cyan, Reset, Yellow, Reset, // for v range ch
					Cyan, Reset, Yellow, Reset, // fmt.Println v
				)},
				{name: "Select", content: fmt.Sprintf(
					("%sSelect Statement%s:\n\n"+
						"\t%sselect%s {\n"+
						"\t%scase%s %smsg%s := <-%sch1%s:\n"+
						"\t\t%sfmt.Println%s(%smsg%s)\n"+
						"\t%scase%s %sch2%s <- %s3%s:\n"+
						"\t\t%sfmt.Println%s(%s\"sent\"%s)\n"+
						"\t%sdefault%s:\n"+
						"\t\t%sfmt.Println%s(%s\"no activity\"%s)\n"+
						"\t}\n"),
					BoldItalic, Reset, // Select Statement
					Cyan, Reset, // select
					Cyan, Reset, Yellow, Reset, Yellow, Reset, // case msg := <-ch1:
					Cyan, Reset, Yellow, Reset, // fmt.Println msg
					Cyan, Reset, Yellow, Reset, Yellow, Reset, // case ch2 <- 3:
					Cyan, Reset, Cyan, Reset, // fmt.Println "sent"
					Cyan, Reset, Cyan, Reset, Green, Reset, // default:
				)},
				{name: "Looping", content: fmt.Sprintf(
					("%sLooping Through Channels%s:\n\n"+
						"\t%sfor%s {\n"+
						"\t\t%smsg%s, %sok%s := <-%sch%s\n"+
						"\t\t%sif%s !%sok%s {\n"+
						"\t\t\t%sbreak%s\n"+
						"\t\t}\n"+
						"\t\t%sfmt.Println%s(%smsg%s)\n"+
						"\t}\n\n"+
						"\t%sfor%s %smsg%s := %srange%s %sch%s {\n"+
						"\t\t%sfmt.Println%s(%smsg%s)\n"+
						"\t}\n"),
					BoldItalic, Reset, // Looping Through Channels
					Cyan, Reset, // for
					Yellow, Reset, Yellow, Reset, Yellow, Reset, // msg, ok := <-ch
					Cyan, Reset, Yellow, Reset, // if !ok
					BoldYellow, Reset, // break
					Cyan, Reset, Yellow, Reset, // fmt.Println msg
					Cyan, Reset, Yellow, Reset, Cyan, Reset, Yellow, Reset, // for msg := range ch
					Cyan, Reset, Yellow, Reset, // fmt.Println msg
				)},
			},
		},
		{
			name: "Goroutines",
			subsections: []subsection{
				{name: "Basic", content: fmt.Sprintf(
					("%sStarting Goroutines%s:\n\n"+
						"\t%sgo%s %sfunc%s() {\n"+
						"\t\t%sfmt.Println%s(%s\"Running\"%s)\n"+
						"\t}()\n"),
					BoldItalic, Reset,
					Cyan, Reset, Cyan, Reset,
					Cyan, Reset, Green, Reset,
				)},
				{name: "WaitGroups", content: fmt.Sprintf(
					("%sUsing WaitGroups%s:\n\n"+
						"\tvar %swg%s sync.WaitGroup\n"+
						"\t%swg%s.Add(%s1%s)\n"+
						"\t%sgo%s func() {\n"+
						"\t\t%sdefer%s %swg%s.Done()\n"+
						"\t\t%sfmt.Println%s(%s\"Done\"%s)\n"+
						"\t}()\n"+
						"\t%swg%s.Wait()\n"),
					BoldItalic, Reset, // Using WaitGroups
					Yellow, Reset, // wg
					Yellow, Reset, Green, Reset, // wg.Add(1)
					Cyan, Reset, // go
					Cyan, Reset, Yellow, Reset, // wg.Done()
					Cyan, Reset, Green, Reset, // fmt.Println "Done"
					Yellow, Reset, // wg
				)},
				{name: "Communication", content: fmt.Sprintf(
					("%sChannel Communication%s:\n\n"+
						"\t%sch%s := make(chan %sstring%s)\n"+
						"\t%sgo%s func() {\n"+
						"\t\t%sch%s <- %s\"ping\"%s\n"+
						"\t}()\n"+
						"\t%smsg%s := <-%sch%s\n"+
						"\t%sfmt.Println%s(%smsg%s)\n"),
					BoldItalic, Reset, // Channel Communication
					Yellow, Reset, Yellow, Reset, // ch make string
					Cyan, Reset, // go
					Yellow, Reset, Green, Reset, // ch <- "ping"
					Yellow, Reset, Yellow, Reset, // msg := <-ch
					Cyan, Reset, Yellow, Reset,	// fmt.Println msg
				)},
			},
		},
		{
			name: "Concurrency",
			subsections: []subsection{
				{name: "Mutex", content: fmt.Sprintf(
					"%sMutex Usage%s:\n\n"+
					"\tvar %smux%s sync.Mutex\n"+
					"\tvar %sval%s int\n\n"+
					"\t%smux%s.Lock()\n"+
					"\t%sval%s++\n"+
					"\t%smux%s.Unlock()",
					BoldItalic, Reset,
					Yellow, Reset, 
					Yellow, Reset,
					Yellow, Reset,
					Yellow, Reset,
					Yellow, Reset,
				)},
				{name: "WorkerPool", content: fmt.Sprintf(
					"%sWorker Pool%s:\n\n"+
					"\t%sworker%s := func(%sjobs%s <-chan int, %sresults%s chan<- int) {\n"+
					"\t\tfor %sj%s := range %sjobs%s {\n"+
					"\t\t\t%sresults%s <- %sj%s * 2\n"+
					"\t\t}\n"+
					"\t}",
					BoldItalic, Reset,
					Cyan, Reset, Yellow, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset,
				)},
			},
		},
		{
			name: "Pointers",
			subsections: []subsection{
				{name: "Basics", content: fmt.Sprintf(
					("%sPointer Basics%s:\n\n"+
						"\tvar %sp%s *%sint%s\n"+
						"\t%si%s := %s42%s\n"+
						"\t%sp%s = &%si%s\n"+
						"\t%sfmt.Println%s(*%sp%s)  %s// 42%s\n"),
					BoldItalic, Reset, // Pointer Basics
					Yellow, Reset, Yellow, Reset, // p *int
					Yellow, Reset, Green, Reset, // i := 42
					Yellow, Reset, Yellow, Reset, // p = &i
					Cyan, Reset, Yellow, Reset, Cyan, Reset, // fmt.Println *p
				)},
				{name: "Structs", content: fmt.Sprintf(
					("%sPointers to Structs%s:\n\n"+
						"\t%stype%s %sVertex%s struct { %sX%s, %sY%s float64 }\n"+
						"\t%sv%s := %sVertex%s{%s1%s, %s2%s}\n"+
						"\t%sp%s := &%sv%s\n"+
						"\t%sp%s.%sX%s = %s1e9%s\n"),
					BoldItalic, Reset, // Pointers to Structs
					Cyan, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset, // type Vertex X Y
					Yellow, Reset, Yellow, Reset, Green, Reset, Green, Reset, // v Vertex 1 2
					Yellow, Reset, Yellow, Reset, // p &v
					Yellow, Reset, Yellow, Reset, Green, Reset, // p X 1e9
				)},
				{name: "Functions", content: fmt.Sprintf(
					("%sFunction Parameters%s:\n\n"+
						"\tfunc %smodify%s(%sp%s *%sint%s) {\n"+
						"\t\t*%sp%s = %s2%s\n"+
						"\t}\n"+
						"\t%si%s := %s1%s\n"+
						"\t%smodify%s(&%si%s)\n"),
					BoldItalic, Reset, // Function Parameters
					Cyan, Reset, Yellow, Reset, Yellow, Reset, // modify p *int
					Yellow, Reset, Green, Reset, // p = 2
					Yellow, Reset, Green, Reset, // i := 1
					Cyan, Reset, Yellow, Reset, // modify &i
				)},
			},
		},
		{
			name: "ErrorHandling",
			subsections: []subsection{
				{name: "Basic", content: fmt.Sprintf(
					("%sBasic Error Handling%s:\n\n"+
						"\t%sfile%s, %serr%s := %sos.Open%s(%s\"file.txt\"%s)\n"+
						"\t%sif%s %serr%s != %snil%s {\n"+
						"\t\t%slog.Fatal%s(%serr%s)\n"+
						"\t}\n"+
						"\t%sdefer%s %sfile%s.Close()\n"),
					BoldItalic, Reset, // Basic Error Handling
					Yellow, Reset, Yellow, Reset, Cyan, Reset, Green, Reset, // file, err := os.Open
					Cyan, Reset, Yellow, Reset, Cyan, Reset, // if err != nil
					Cyan, Reset, Yellow, Reset, // log.Fatal err
					Cyan, Reset, Yellow, Reset, // defer file.Close
				)},
				{name: "Custom", content: fmt.Sprintf(
					("%sCustom Errors%s:\n\n"+
						"\t%stype%s %sMyError%s struct {\n"+
						"\t\t%sMsg%s string\n"+
						"\t}\n\n"+
						"\tfunc (%se%s *%sMyError%s) %sError%s() string {\n"+
						"\t\treturn %se%s.%sMsg%s\n"+
						"\t}\n"),
					BoldItalic, Reset, // Custom Errors
					Cyan, Reset, Yellow, Reset, // type MyError
					Yellow, Reset, // Msg
					Yellow, Reset, Yellow, Reset, Green, Reset, // e MyError Error
					Cyan, Reset, Yellow, Reset, // e Msg
				)},
				{name: "PanicRecover", content: fmt.Sprintf(
					("%sPanic and Recover%s:\n\n"+
						"\tfunc %smayPanic%s() {\n"+
						"\t\t%spanic%s(%s\"problem\"%s)\n"+
						"\t}\n\n"+
						"\t%sdefer%s func() {\n"+
						"\t\tif %sr%s := %srecover%s(); %sr%s != nil {\n"+
						"\t\t\t%sfmt.Println%s(%s\"Recovered:\"%s, %sr%s)\n"+
						"\t\t}\n"+
						"\t}()\n"+
						"\t%smayPanic%s()\n"),
					BoldItalic, Reset, // Panic and Recover
					Cyan, Reset, Green, Reset, // mayPanic
					Cyan, Reset, Green, Reset, // panic "problem"
					Cyan, Reset, // defer
					Yellow, Reset, Cyan, Reset, Yellow, Reset, // r recover r 
					Cyan, Reset, Green, Reset, // fmt.Println "Recovered:"
					Cyan, Reset, //mayPanic 
				)},
			},
		},
		{
			name: "Testing",
			subsections: []subsection{
				{name: "UnitTests", content: fmt.Sprintf(
					"%sUnit Test%s:\n\n"+
					"\tfunc %sTestAdd%s(%st%s *testing.T) {\n"+
					"\t\tgot := %sadd%s(2, 3)\n"+
					"\t\twant := 5\n"+
					"\t\t%sif%s got != want {\n"+
					"\t\t\t%st%s.Errorf(%s\"got %%d want %%d\"%s, got, want)\n"+
					"\t\t}\n"+
					"\t}",
					BoldItalic, Reset,
					Cyan, Reset, Yellow, Reset,
					Cyan, Reset,
					Cyan, Reset,
					Yellow, Reset, Green, Reset,
				)},
				{name: "Benchmarks", content: fmt.Sprintf(
					"%sBenchmark%s:\n\n"+
					"\tfunc %sBenchmarkAdd%s(%sb%s *testing.B) {\n"+
					"\t\tfor %si%s := 0; %si%s < %sb%s.N; %si%s++ {\n"+
					"\t\t\t%sadd%s(1, 2)\n"+
					"\t\t}\n"+
					"\t}",
					BoldItalic, Reset,
					Cyan, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset, Yellow, Reset, Yellow, Reset,
					Cyan, Reset,
				)},
			},
		},
		{
			name: "StringManipulation",
			subsections: []subsection{
				{name: "Basic", content: fmt.Sprintf(
					("%sBasic Operations%s:\n\n"+
						"\t%ss1%s := %s\"Hello\"%s\n"+
						"\t%ss2%s := %s\"World\"%s\n"+
						"\t%ss3%s := %ss1%s + %s\" \"%s + %ss2%s\n"+
						"\t%sfmt.Println%s(%slen%s(%ss3%s))  %s// 11%s\n"),
					BoldItalic, Reset, // Basic Operations
					Yellow, Reset, Green, Reset, // s1 := "Hello"
					Yellow, Reset, Green, Reset, // s2 := "World"
					Yellow, Reset, Yellow, Reset, Green, Reset, Yellow, Reset, // s3 := s1 + " " + s2
					Cyan, Reset, Cyan, Reset, Yellow, Reset, Cyan, Reset, // fmt.Println len s3
				)},
				{name: "StringsPackage", content: fmt.Sprintf(
					("%sStrings Package%s:\n\n"+
						"\t%sstrings.Split%s(%s\"a,b,c\"%s, %s\",\"%s)\n"+
						"\t%sstrings.ToUpper%s(%s\"test\"%s)\n"+
						"\t%sstrings.TrimSpace%s(%s\"  text  \"%s)\n"),
					BoldItalic, Reset, // Strings Package
					Cyan, Reset, Green, Reset, Green, Reset, // strings.Split
					Cyan, Reset, Green, Reset, // strings.ToUpper
					Cyan, Reset, Green, Reset, // strings.TrimSpace
				)},
				{name: "Conversions", content: fmt.Sprintf(
					("%sType Conversions%s:\n\n"+
						"\t%si%s, %s_%s := %sstrconv.Atoi%s(%s\"42\"%s)\n"+
						"\t%ss%s := %sstrconv.Itoa%s(%s42%s)\n"),
					BoldItalic, Reset, // Type Conversions
					Yellow, Reset, Yellow, Reset, Cyan, Reset, Green, Reset, // i, _ := strconv.Atoi
					Yellow, Reset, Cyan, Reset, Green, Reset, // s := strconv.Itoa
				)},
			},
		},
		{
			name: "PrintFormatting",
			subsections: []subsection{
				{name: "PrintFunctions", content: fmt.Sprintf(
					("%sPrint Functions%s:\n\n"+
						"\t%sfmt.Print%s(%s\"Hello\"%s)\n"+
						"\t%sfmt.Println%s(%s\"World\"%s)\n"+
						"\t%sfmt.Printf%s(%s\"Value: %%v\"%s, %s42%s)\n"),
					BoldItalic, Reset, // Print Functions
					Cyan, Reset, Green, Reset, // fmt.Print "Hello"
					Cyan, Reset, Green, Reset, // fmt.Println "World"
					Cyan, Reset, Green, Reset, Green, Reset, // fmt.Printf "Value: %v"
				)},
				{name: "FormatVerbs", content: fmt.Sprintf(
					("%sFormat Verbs%s:\n\n"+
						"\t%s%%v%s - Value\n"+
						"\t%s%%s%s - String\n"+
						"\t%s%%d%s - Integer\n"+
						"\t%s%%f%s - Float\n"+
						"\t%s%%t%s - Boolean\n"+
						"\t%s%%T%s - Type\n"),
					BoldItalic, Reset, // Format Verbs
					BoldPurple, Reset, // %v
					BoldPurple, Reset, // %s
					BoldPurple, Reset, // %d
					BoldPurple, Reset, // %f
					BoldPurple, Reset, // %t
					BoldPurple, Reset, // %T
				)},
				{name: "Sprintf", content: fmt.Sprintf(
					("%sString Formatting%s:\n\n"+
						"\t%ss%s := %sfmt.Sprintf%s(%s\"Name: %%s, Age: %%d\"%s, %s\"Alice\"%s, %s30%s)\n"+
						"\t%sfmt.Fprintf%s(%sos.Stderr%s, %s\"Error: %%v\"%s, %serr%s)\n"),
					BoldItalic, Reset, // String Formatting
					Yellow, Reset, Cyan, Reset, Green, Reset, Green, Reset, Green, Reset, // s := fmt.Sprintf
					Cyan, Reset, Cyan, Reset, Green, Reset, Yellow, Reset, // fmt.Fprintf
				)},
			},
		},
		{
			name: "FileIO",
			subsections: []subsection{
				{name: "ReadWrite", content: fmt.Sprintf(
					"%sRead/Write Files%s:\n\n"+
					"\t%sdata%s := []byte(%s\"hello\\nworld\"%s)\n"+
					"\t%serr%s := %sos.WriteFile%s(%s\"file.txt\"%s, %sdata%s, 0644)\n\n"+
					"\t%scontent%s, %serr%s := %sos.ReadFile%s(%s\"file.txt\"%s)",
					BoldItalic, Reset,
					Yellow, Reset, Green, Reset,
					Yellow, Reset, Cyan, Reset, Green, Reset, Yellow, Reset,
					Yellow, Reset, Yellow, Reset, Cyan, Reset, Green, Reset,
				)},
			},
		},
		{
			name: "Time",
			subsections: []subsection{
				{name: "Formatting", content: fmt.Sprintf(
					"%sTime Formatting%s:\n\n"+
					"\t%st%s := %stime.Now%s()\n"+
					"\t%sfmt.Println%s(%st%s.Format(%s\"2006-01-02 15:04:05\"%s))",
					BoldItalic, Reset,
					Yellow, Reset, Cyan, Reset,
					Cyan, Reset, Yellow, Reset, Green, Reset,
				)},
			},
		},
		{
			name: "HTTPServer",
			subsections: []subsection{
				{name: "BasicServer", content: fmt.Sprintf(
					"%sBasic Server%s:\n\n"+
					"\t%shttp.HandleFunc%s(%s\"/\"%s, func(%sw%s http.ResponseWriter, %sr%s *http.Request) {\n"+
					"\t\t%sfmt.Fprintf%s(%sw%s, %s\"Hello World\"%s)\n"+
					"\t})\n"+
					"\t%shttp.ListenAndServe%s(%s\":8080\"%s, nil)",
					BoldItalic, Reset,
					Cyan, Reset, Green, Reset,
					Yellow, Reset, Yellow, Reset,
					Cyan, Reset, Yellow, Reset, Green, Reset,
					Cyan, Reset, Green, Reset,
				)},
			},
		},
		{
			name: "PackageManagement",
			subsections: []subsection{
				{name: "GoMod", content: fmt.Sprintf(
					("%sgo.mod Example%s:\n\n"+
						"\tmodule %sgithub.com/yourname/project%s\n\n"+
						"\tgo %s1.21%s\n\n"+
						"\trequire (\n"+
						"\t\t%sgithub.com/pkg/errors%s %sv0.9.1%s\n"+
						"\t)"),
					BoldItalic, Reset, // go.mod Example
					Green, Reset, // 1.21
					Green, Reset, // github.com/pkg/errors
					Green, Reset, Green, Reset, // v0.9.1
				)},
				{name: "Dependencies", content: fmt.Sprintf(
					("%sDependency Management%s:\n\n"+
						"\t%sgo get%s %sgithub.com/pkg/errors@latest%s\n"+
						"\t%sgo mod tidy%s\n"+
						"\t%sgo list -m all%s\n"+
						"\t%sgo mod vendor%s"),
					BoldItalic, Reset, // Dependency Management
					Cyan, Reset, Green, Reset, // github.com/pkg/errors@latest
					Cyan, Reset, // go mod tidy
					Cyan, Reset, // go list -m all
					Cyan, Reset, // go mod vendor
				)},
				{name: "Vendoring", content: fmt.Sprintf(
					("%sVendor Directory%s:\n\n"+
						"\t%sgo mod vendor%s\n"+
						"\t%sgo build -mod=vendor%s\n"+
						"\t%s// vendor/modules.txt contains dependency info%s"),
					BoldItalic, Reset, // Vendor Directory
					Cyan, Reset, // go mod vendor
					Cyan, Reset, // go build -mod=vendor
					Italic, Reset, // vendor/modules.txt
				)},
			},
		},
		{
			name: "BuildRun",
			subsections: []subsection{
				{name: "Commands", content: fmt.Sprintf(
					("%sBuild Commands%s:\n\n"+
						"\t%sgo build%s %s./cmd/app%s\n"+
						"\t%sgo run%s %smain.go%s\n"+
						"\t%sgo install%s %sgithub.com/project/cmd/app%s\n"+
						"\t%sGOOS=linux GOARCH=amd64 go build%s"),
					BoldItalic, Reset, // Build Commands
					Cyan, Reset, Green, Reset, // go build ./cmd/app
					Cyan, Reset, Green, Reset, // go run main.go
					Cyan, Reset, Green, Reset, // go install github.com/project/cmd/app
					Cyan, Reset, Green, Reset, // GOOS=linux GOARCH=amd64 go build
				)},
				{name: "MultiModule", content: fmt.Sprintf(
					("%sLocal Modules%s:\n\n"+
						"\t// go.work file\n"+
						"\tuse (\n"+
						"\t\t./lib\n"+
						"\t\t./app\n"+
						"\t)\n\n"+
						"\t// go.mod\n"+
						"\treplace %slocal/mypackage%s => %s../mypackage%s"),
					BoldItalic, Reset, // Local Modules
					Green, Reset, Green, Reset, // local/mypackage
				)},
			},
		},
		{
			name: "Reflection",
			subsections: []subsection{
				{name: "Basic", content: fmt.Sprintf(
					("%sBasic Reflection%s:\n\n"+
						"\t%st%s := %sreflect.TypeOf%s(%s42%s)\n"+
						"\t%sv%s := %sreflect.ValueOf%s(%s\"hello\"%s)\n"+
						"\t%sfmt.Println%s(%st%s.Kind(), %sv%s.Len())"),
					BoldItalic, Reset, // Basic Reflection
					Yellow, Reset, Cyan, Reset, Green, Reset, // t := reflect.TypeOf
					Yellow, Reset, Cyan, Reset, Green, Reset, // v := reflect.ValueOf
					Cyan, Reset, Yellow, Reset, Yellow, Reset, // fmt.Println
				)},
				{name: "Structs", content: fmt.Sprintf(
					("%sStruct Reflection%s:\n\n"+
						"\t%stype%s %sPerson%s struct {\n"+
						"\t\t%sName%s string\n"+
						"\t}\n\n"+
						"\t%sp%s := %sPerson%s{%s\"Alice\"%s}\n"+
						"\t%sv%s := %sreflect.ValueOf%s(%sp%s)\n"+
						"\t%sf%s := %sv%s.FieldByName(%s\"Name\"%s)"),
					BoldItalic, Reset, // Struct Reflection
					Cyan, Reset, Yellow, Reset, // type Person
					Yellow, Reset, // Name 
					Yellow, Reset, Yellow, Reset, Green, Reset, // p Person "Alice"
					Yellow, Reset, Cyan, Reset, Yellow, Reset, // v reflect.ValueOf p
					Yellow, Reset, Yellow, Reset, Green, Reset), // f v.FieldByName
				},
			},
		},
		{
			name: "ImportsVisibility",
			subsections: []subsection{
				{name: "Imports", content: fmt.Sprintf(
					("%sImport Statements%s:\n\n"+
						"\timport (\n"+
						"\t\t%s\"fmt\"%s\n"+
						"\t\t%s\"github.com/user/pkg\"%s\n"+
						"\t\t%s\"./local\"%s\n"+
						"\t)"),
					BoldItalic, Reset, // Import Statements
					Green, Reset, // fmt
					Green, Reset, // github.com/user/pkg
					Green, Reset, // ./local
				)},
				{name: "Visibility", content: fmt.Sprintf(
					("%sPublic/Private%s:\n\n"+
						"\t// Public (exported)\n"+
						"\t%svar%s %sGlobalVar%s int\n"+
						"\t%sfunc%s %sPublicFunc%s() {}\n\n"+
						"\t// Private (unexported)\n"+
						"\t%svar%s %slocalVar%s int\n"+
						"\t%sfunc%s %sprivateFunc%s() {}"),
					BoldItalic, Reset, // Public/Private
					Cyan, Reset, Yellow, Reset, // var GlobalVar
					Cyan, Reset, Yellow, Reset, // func PublicFunc
					Cyan, Reset, Yellow, Reset, // var localVar
					Cyan, Reset, Yellow, Reset, // func privateFunc
				)},
			},
		},
		{
			name: "Generics",
			subsections: []subsection{
				{name: "Basic", content: fmt.Sprintf(
					("%sGeneric Function%s:\n\n"+
						"\tfunc %sPrintSlice%s[%sT%s %sany%s](%ss%s []%sT%s) {\n"+
						"\t\tfor _, %sv%s := range %ss%s {\n"+
						"\t\t\t%sfmt.Print%s(%sv%s)\n"+
						"\t\t}\n"+
						"\t}"),
					BoldItalic, Reset, // Generic Function
					Cyan, Reset, Yellow, Reset, Cyan, Reset, // PrintSlice T any
					Yellow, Reset, Yellow, Reset, // s []T
					Yellow, Reset, Yellow, Reset, // v s
					Cyan, Reset, Yellow, Reset, // fmt.Print v
				)},
				{name: "Constraints", content: fmt.Sprintf(
					("%sType Constraints%s:\n\n"+
						"\ttype %sNumber%s interface {\n"+
						"\t\t%sint%s | %sfloat64%s\n"+
						"\t}\n\n"+
						"\tfunc %sSum%s[%sT%s %sNumber%s](%snums%s []%sT%s) %sT%s {\n"+
						"\t\tvar total %sT%s\n"+
						"\t\tfor _, %sn%s := range %snums%s {\n"+
						"\t\t\ttotal += %sn%s\n"+
						"\t\t}\n"+
						"\t\treturn total\n"+
						"\t}"),
					BoldItalic, Reset, // Type Constraints
					Yellow, Reset, // Number
					Yellow, Reset, Yellow, Reset, // int float64
					Cyan, Reset, Yellow, Reset, Yellow, Reset, // Sum T Number
					Yellow, Reset, Yellow, Reset, Yellow, Reset, // nums []T
					Yellow, Reset, // n nums
					Yellow, Reset, Yellow, Reset, // total += n
					Yellow, Reset, // total
				)},
				{name: "GenericStruct", content: fmt.Sprintf(
					("%sGeneric Struct%s:\n\n"+
						"\ttype %sContainer%s[%sT%s %sany%s] struct {\n"+
						"\t\tValue %sT%s\n"+
						"\t}\n\n"+
						"\tfunc (%sc%s *%sContainer%s[%sT%s]) %sGet%s() %sT%s {\n"+
						"\t\treturn %sc%s.Value\n"+
						"\t}"),
					BoldItalic, Reset, // Generic Struct
					Yellow, Reset, Yellow, Reset, Cyan, Reset, // Container T any
					Yellow, Reset, Yellow, Reset, // Value T
					Yellow, Reset, Yellow, Reset, // c Container
					Cyan, Reset, Yellow, Reset, // Get T
					Yellow, Reset, // c Value
				)},
			},
		},
	}
	return sections
}