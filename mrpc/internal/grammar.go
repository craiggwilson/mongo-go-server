package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/craiggwilson/mongo-go-server/mrpc/tree"
)

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Tree",
			pos:  position{line: 14, col: 1, offset: 219},
			expr: &actionExpr{
				pos: position{line: 14, col: 9, offset: 227},
				run: (*parser).callonTree1,
				expr: &seqExpr{
					pos: position{line: 14, col: 9, offset: 227},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 14, col: 9, offset: 227},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 14, col: 11, offset: 229},
							label: "v",
							expr: &zeroOrMoreExpr{
								pos: position{line: 14, col: 13, offset: 231},
								expr: &choiceExpr{
									pos: position{line: 14, col: 14, offset: 232},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 14, col: 14, offset: 232},
											name: "Attribute",
										},
										&ruleRefExpr{
											pos:  position{line: 14, col: 26, offset: 244},
											name: "Command",
										},
										&ruleRefExpr{
											pos:  position{line: 14, col: 36, offset: 254},
											name: "Service",
										},
										&ruleRefExpr{
											pos:  position{line: 14, col: 46, offset: 264},
											name: "Struct",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 14, col: 55, offset: 273},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Attribute",
			pos:  position{line: 34, col: 1, offset: 744},
			expr: &actionExpr{
				pos: position{line: 34, col: 14, offset: 757},
				run: (*parser).callonAttribute1,
				expr: &seqExpr{
					pos: position{line: 34, col: 14, offset: 757},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 34, col: 14, offset: 757},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 34, col: 16, offset: 759},
							val:        "@",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 34, col: 20, offset: 763},
							label: "lang",
							expr: &zeroOrOneExpr{
								pos: position{line: 34, col: 25, offset: 768},
								expr: &ruleRefExpr{
									pos:  position{line: 34, col: 25, offset: 768},
									name: "AttributeLang",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 34, col: 40, offset: 783},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 34, col: 45, offset: 788},
								name: "ID",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 34, col: 48, offset: 791},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 34, col: 51, offset: 794},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 34, col: 57, offset: 800},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 34, col: 64, offset: 807},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 34, col: 66, offset: 809},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AttributeLang",
			pos:  position{line: 42, col: 1, offset: 977},
			expr: &actionExpr{
				pos: position{line: 42, col: 18, offset: 994},
				run: (*parser).callonAttributeLang1,
				expr: &seqExpr{
					pos: position{line: 42, col: 18, offset: 994},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 42, col: 18, offset: 994},
							label: "lang",
							expr: &ruleRefExpr{
								pos:  position{line: 42, col: 23, offset: 999},
								name: "ID",
							},
						},
						&litMatcher{
							pos:        position{line: 42, col: 26, offset: 1002},
							val:        ":",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Command",
			pos:  position{line: 46, col: 1, offset: 1045},
			expr: &actionExpr{
				pos: position{line: 46, col: 12, offset: 1056},
				run: (*parser).callonCommand1,
				expr: &seqExpr{
					pos: position{line: 46, col: 12, offset: 1056},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 46, col: 12, offset: 1056},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 46, col: 14, offset: 1058},
							val:        "command",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 46, col: 24, offset: 1068},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 46, col: 27, offset: 1071},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 46, col: 32, offset: 1076},
								name: "ID",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 46, col: 35, offset: 1079},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 46, col: 37, offset: 1081},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 46, col: 41, offset: 1085},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 46, col: 43, offset: 1087},
							label: "attrs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 46, col: 49, offset: 1093},
								expr: &ruleRefExpr{
									pos:  position{line: 46, col: 49, offset: 1093},
									name: "Attribute",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 46, col: 60, offset: 1104},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 46, col: 62, offset: 1106},
							label: "req",
							expr: &zeroOrOneExpr{
								pos: position{line: 46, col: 66, offset: 1110},
								expr: &ruleRefExpr{
									pos:  position{line: 46, col: 66, offset: 1110},
									name: "Request",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 46, col: 75, offset: 1119},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 46, col: 77, offset: 1121},
							label: "resp",
							expr: &ruleRefExpr{
								pos:  position{line: 46, col: 82, offset: 1126},
								name: "Response",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 46, col: 91, offset: 1135},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 46, col: 93, offset: 1137},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Request",
			pos:  position{line: 61, col: 1, offset: 1428},
			expr: &actionExpr{
				pos: position{line: 61, col: 12, offset: 1439},
				run: (*parser).callonRequest1,
				expr: &seqExpr{
					pos: position{line: 61, col: 12, offset: 1439},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 61, col: 12, offset: 1439},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 14, offset: 1441},
							val:        "request",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 24, offset: 1451},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 26, offset: 1453},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 30, offset: 1457},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 32, offset: 1459},
							label: "attrs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 61, col: 38, offset: 1465},
								expr: &ruleRefExpr{
									pos:  position{line: 61, col: 38, offset: 1465},
									name: "Attribute",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 49, offset: 1476},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 61, col: 51, offset: 1478},
							label: "fields",
							expr: &zeroOrMoreExpr{
								pos: position{line: 61, col: 58, offset: 1485},
								expr: &ruleRefExpr{
									pos:  position{line: 61, col: 58, offset: 1485},
									name: "Field",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 61, col: 65, offset: 1492},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 61, col: 67, offset: 1494},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Response",
			pos:  position{line: 75, col: 1, offset: 1753},
			expr: &actionExpr{
				pos: position{line: 75, col: 13, offset: 1765},
				run: (*parser).callonResponse1,
				expr: &seqExpr{
					pos: position{line: 75, col: 13, offset: 1765},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 75, col: 13, offset: 1765},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 75, col: 15, offset: 1767},
							val:        "response",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 26, offset: 1778},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 75, col: 28, offset: 1780},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 32, offset: 1784},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 75, col: 34, offset: 1786},
							label: "attrs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 75, col: 40, offset: 1792},
								expr: &ruleRefExpr{
									pos:  position{line: 75, col: 40, offset: 1792},
									name: "Attribute",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 51, offset: 1803},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 75, col: 53, offset: 1805},
							label: "fields",
							expr: &zeroOrMoreExpr{
								pos: position{line: 75, col: 60, offset: 1812},
								expr: &ruleRefExpr{
									pos:  position{line: 75, col: 60, offset: 1812},
									name: "Field",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 75, col: 67, offset: 1819},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 75, col: 69, offset: 1821},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Service",
			pos:  position{line: 89, col: 1, offset: 2084},
			expr: &actionExpr{
				pos: position{line: 89, col: 12, offset: 2095},
				run: (*parser).callonService1,
				expr: &seqExpr{
					pos: position{line: 89, col: 12, offset: 2095},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 89, col: 12, offset: 2095},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 89, col: 14, offset: 2097},
							val:        "service",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 24, offset: 2107},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 89, col: 27, offset: 2110},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 89, col: 32, offset: 2115},
								name: "ID",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 35, offset: 2118},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 89, col: 37, offset: 2120},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 41, offset: 2124},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 89, col: 43, offset: 2126},
							label: "attrs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 89, col: 49, offset: 2132},
								expr: &ruleRefExpr{
									pos:  position{line: 89, col: 49, offset: 2132},
									name: "Attribute",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 60, offset: 2143},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 89, col: 62, offset: 2145},
							label: "cmds",
							expr: &zeroOrMoreExpr{
								pos: position{line: 89, col: 67, offset: 2150},
								expr: &ruleRefExpr{
									pos:  position{line: 89, col: 67, offset: 2150},
									name: "CommandRef",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 89, col: 79, offset: 2162},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 89, col: 81, offset: 2164},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CommandRef",
			pos:  position{line: 103, col: 1, offset: 2441},
			expr: &actionExpr{
				pos: position{line: 103, col: 15, offset: 2455},
				run: (*parser).callonCommandRef1,
				expr: &seqExpr{
					pos: position{line: 103, col: 15, offset: 2455},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 103, col: 15, offset: 2455},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 103, col: 17, offset: 2457},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 103, col: 22, offset: 2462},
								name: "ID",
							},
						},
						&litMatcher{
							pos:        position{line: 103, col: 25, offset: 2465},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 107, col: 1, offset: 2508},
			expr: &actionExpr{
				pos: position{line: 107, col: 11, offset: 2518},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 107, col: 11, offset: 2518},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 107, col: 11, offset: 2518},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 107, col: 13, offset: 2520},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 22, offset: 2529},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 107, col: 25, offset: 2532},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 107, col: 30, offset: 2537},
								name: "ID",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 33, offset: 2540},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 107, col: 35, offset: 2542},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 39, offset: 2546},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 107, col: 41, offset: 2548},
							label: "attrs",
							expr: &zeroOrMoreExpr{
								pos: position{line: 107, col: 47, offset: 2554},
								expr: &ruleRefExpr{
									pos:  position{line: 107, col: 47, offset: 2554},
									name: "Attribute",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 58, offset: 2565},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 107, col: 60, offset: 2567},
							label: "fields",
							expr: &zeroOrMoreExpr{
								pos: position{line: 107, col: 67, offset: 2574},
								expr: &ruleRefExpr{
									pos:  position{line: 107, col: 67, offset: 2574},
									name: "Field",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 74, offset: 2581},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 107, col: 76, offset: 2583},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 121, col: 1, offset: 2853},
			expr: &actionExpr{
				pos: position{line: 121, col: 10, offset: 2862},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 121, col: 10, offset: 2862},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 121, col: 10, offset: 2862},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 121, col: 12, offset: 2864},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 121, col: 17, offset: 2869},
								name: "ID",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 121, col: 20, offset: 2872},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 121, col: 23, offset: 2875},
							label: "r",
							expr: &ruleRefExpr{
								pos:  position{line: 121, col: 25, offset: 2877},
								name: "FieldTypeRef",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 121, col: 38, offset: 2890},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 121, col: 40, offset: 2892},
							val:        ";",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FieldTypeRef",
			pos:  position{line: 127, col: 1, offset: 2990},
			expr: &actionExpr{
				pos: position{line: 127, col: 17, offset: 3006},
				run: (*parser).callonFieldTypeRef1,
				expr: &seqExpr{
					pos: position{line: 127, col: 17, offset: 3006},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 127, col: 17, offset: 3006},
							expr: &litMatcher{
								pos:        position{line: 127, col: 17, offset: 3006},
								val:        "[]",
								ignoreCase: false,
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 127, col: 23, offset: 3012},
							expr: &litMatcher{
								pos:        position{line: 127, col: 23, offset: 3012},
								val:        "*",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 127, col: 28, offset: 3017},
							name: "TypeRef",
						},
					},
				},
			},
		},
		{
			name: "ID",
			pos:  position{line: 131, col: 1, offset: 3065},
			expr: &actionExpr{
				pos: position{line: 131, col: 7, offset: 3071},
				run: (*parser).callonID1,
				expr: &seqExpr{
					pos: position{line: 131, col: 7, offset: 3071},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 131, col: 7, offset: 3071},
							val:        "[a-zA-Z]",
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 131, col: 16, offset: 3080},
							expr: &charClassMatcher{
								pos:        position{line: 131, col: 16, offset: 3080},
								val:        "[a-zA-Z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "TypeRef",
			pos:  position{line: 135, col: 1, offset: 3134},
			expr: &actionExpr{
				pos: position{line: 135, col: 12, offset: 3145},
				run: (*parser).callonTypeRef1,
				expr: &seqExpr{
					pos: position{line: 135, col: 12, offset: 3145},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 135, col: 12, offset: 3145},
							expr: &litMatcher{
								pos:        position{line: 135, col: 12, offset: 3145},
								val:        "$",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 135, col: 17, offset: 3150},
							name: "ID",
						},
						&zeroOrMoreExpr{
							pos: position{line: 135, col: 20, offset: 3153},
							expr: &seqExpr{
								pos: position{line: 135, col: 21, offset: 3154},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 135, col: 21, offset: 3154},
										val:        ".",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 135, col: 25, offset: 3158},
										name: "ID",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "String",
			pos:  position{line: 139, col: 1, offset: 3203},
			expr: &actionExpr{
				pos: position{line: 139, col: 11, offset: 3213},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 139, col: 11, offset: 3213},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 139, col: 11, offset: 3213},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 139, col: 15, offset: 3217},
							expr: &choiceExpr{
								pos: position{line: 139, col: 17, offset: 3219},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 139, col: 17, offset: 3219},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 139, col: 17, offset: 3219},
												expr: &ruleRefExpr{
													pos:  position{line: 139, col: 18, offset: 3220},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 139, col: 30, offset: 3232,
											},
										},
									},
									&seqExpr{
										pos: position{line: 139, col: 34, offset: 3236},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 139, col: 34, offset: 3236},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 139, col: 39, offset: 3241},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 139, col: 57, offset: 3259},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 144, col: 1, offset: 3382},
			expr: &charClassMatcher{
				pos:        position{line: 144, col: 16, offset: 3397},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 145, col: 1, offset: 3413},
			expr: &choiceExpr{
				pos: position{line: 145, col: 19, offset: 3431},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 145, col: 19, offset: 3431},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 145, col: 38, offset: 3450},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 146, col: 1, offset: 3465},
			expr: &charClassMatcher{
				pos:        position{line: 146, col: 21, offset: 3485},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 147, col: 1, offset: 3498},
			expr: &seqExpr{
				pos: position{line: 147, col: 18, offset: 3515},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 147, col: 18, offset: 3515},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 147, col: 22, offset: 3519},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 147, col: 31, offset: 3528},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 147, col: 40, offset: 3537},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 147, col: 49, offset: 3546},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 149, col: 1, offset: 3558},
			expr: &zeroOrMoreExpr{
				pos: position{line: 149, col: 19, offset: 3576},
				expr: &charClassMatcher{
					pos:        position{line: 149, col: 19, offset: 3576},
					val:        "[ \\r\\n\\t]",
					chars:      []rune{' ', '\r', '\n', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name:        "__",
			displayName: "\"whitespace\"",
			pos:         position{line: 150, col: 1, offset: 3588},
			expr: &oneOrMoreExpr{
				pos: position{line: 150, col: 20, offset: 3607},
				expr: &charClassMatcher{
					pos:        position{line: 150, col: 20, offset: 3607},
					val:        "[ \\r\\n\\t]",
					chars:      []rune{' ', '\r', '\n', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 151, col: 1, offset: 3619},
			expr: &seqExpr{
				pos: position{line: 151, col: 8, offset: 3626},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 151, col: 8, offset: 3626},
						name: "_",
					},
					&notExpr{
						pos: position{line: 151, col: 10, offset: 3628},
						expr: &anyMatcher{
							line: 151, col: 11, offset: 3629,
						},
					},
				},
			},
		},
	},
}

func (c *current) onTree1(v interface{}) (interface{}, error) {

	t := tree.New()
	for _, i := range toIfaceSlice(v) {
		switch ti := i.(type) {
		case *tree.Attribute:
			t.AddAttribute(ti)
		case *tree.Command:
			t.AddCommand(ti)
		case *tree.Service:
			t.AddService(ti)
		case *tree.Struct:
			t.AddStruct(ti)
		default:
			return nil, errors.New("invalid top-level declaration")
		}
	}

	return t, nil
}

func (p *parser) callonTree1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTree1(stack["v"])
}

func (c *current) onAttribute1(lang, name, value interface{}) (interface{}, error) {

	langStr := ""
	if lang != nil {
		langStr = lang.(string)
	}
	return tree.NewAttribute(langStr, name.(string), value.(string)), nil
}

func (p *parser) callonAttribute1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAttribute1(stack["lang"], stack["name"], stack["value"])
}

func (c *current) onAttributeLang1(lang interface{}) (interface{}, error) {

	return lang.(string), nil
}

func (p *parser) callonAttributeLang1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAttributeLang1(stack["lang"])
}

func (c *current) onCommand1(name, attrs, req, resp interface{}) (interface{}, error) {

	t := tree.NewCommand(name.(string))

	for _, attr := range toIfaceSlice(attrs) {
		t.AddAttribute(attr.(*tree.Attribute))
	}

	if req != nil {
		t.Request = req.(*tree.Struct)
	}
	t.Response = resp.(*tree.Struct)

	return t, nil
}

func (p *parser) callonCommand1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCommand1(stack["name"], stack["attrs"], stack["req"], stack["resp"])
}

func (c *current) onRequest1(attrs, fields interface{}) (interface{}, error) {

	s := tree.NewStruct("")

	for _, attr := range toIfaceSlice(attrs) {
		s.AddAttribute(attr.(*tree.Attribute))
	}

	for _, f := range toIfaceSlice(fields) {
		s.AddField(f.(*tree.Field))
	}

	return s, nil
}

func (p *parser) callonRequest1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRequest1(stack["attrs"], stack["fields"])
}

func (c *current) onResponse1(attrs, fields interface{}) (interface{}, error) {

	s := tree.NewStruct("")

	for _, attr := range toIfaceSlice(attrs) {
		s.AddAttribute(attr.(*tree.Attribute))
	}

	for _, f := range toIfaceSlice(fields) {
		s.AddField(f.(*tree.Field))
	}

	return s, nil
}

func (p *parser) callonResponse1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onResponse1(stack["attrs"], stack["fields"])
}

func (c *current) onService1(name, attrs, cmds interface{}) (interface{}, error) {

	s := tree.NewService(name.(string))

	for _, attr := range toIfaceSlice(attrs) {
		s.AddAttribute(attr.(*tree.Attribute))
	}

	for _, cmd := range toIfaceSlice(cmds) {
		s.AddCommandRef(cmd.(string))
	}

	return s, nil
}

func (p *parser) callonService1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onService1(stack["name"], stack["attrs"], stack["cmds"])
}

func (c *current) onCommandRef1(name interface{}) (interface{}, error) {

	return name.(string), nil
}

func (p *parser) callonCommandRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCommandRef1(stack["name"])
}

func (c *current) onStruct1(name, attrs, fields interface{}) (interface{}, error) {

	s := tree.NewStruct(name.(string))

	for _, attr := range toIfaceSlice(attrs) {
		s.AddAttribute(attr.(*tree.Attribute))
	}

	for _, f := range toIfaceSlice(fields) {
		s.AddField(f.(*tree.Field))
	}

	return s, nil
}

func (p *parser) callonStruct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStruct1(stack["name"], stack["attrs"], stack["fields"])
}

func (c *current) onField1(name, r interface{}) (interface{}, error) {

	f := tree.NewField(name.(string))
	f.TypeRef = r.(string)
	return f, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["name"], stack["r"])
}

func (c *current) onFieldTypeRef1() (interface{}, error) {

	return string(c.text), nil
}

func (p *parser) callonFieldTypeRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldTypeRef1()
}

func (c *current) onID1() (interface{}, error) {

	return string(c.text), nil
}

func (p *parser) callonID1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onID1()
}

func (c *current) onTypeRef1() (interface{}, error) {

	return string(c.text), nil
}

func (p *parser) callonTypeRef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeRef1()
}

func (c *current) onString1() (interface{}, error) {

	c.text = bytes.Replace(c.text, []byte(`\/`), []byte(`/`), -1)
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
