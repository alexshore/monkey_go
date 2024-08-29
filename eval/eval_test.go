package eval

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

// expressions :)
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.expected)
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`"hello world"`, "hello world"},
		{`"Hello" + " " + "World!"`, "Hello World!"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		str, ok := result.(*object.String)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", result, result)
		}

		if str.Value != tt.expected {
			t.Errorf("String has wrong value. got=`%q`", str.Value)
		}
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!false", true},
		{"!true", false},
		{"!!false", false},
		{"!!true", true},
		{"!5", false},
		{"!!5", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"Hello" == "World!"`, false},
		{`"Hello" != "World!"`, true},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testBooleanObject(t, result, tt.expected)
	}
}

func TestEvalIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, result, int64(integer))
		} else {
			testNullObject(t, result)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := `[1, 2 * 2, 3 + 3]`

	result := testEval(input)
	arr, ok := result.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", result, result)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d, expected=3", len(arr.Elements))
	}

	testIntegerObject(t, arr.Elements[0], 1)
	testIntegerObject(t, arr.Elements[1], 4)
	testIntegerObject(t, arr.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
		{
			"[][0]",
			nil,
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		expected, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, result, int64(expected))
		} else {
			testNullObject(t, result)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	result := testEval(input)
	hash, ok := result.(*object.Hash)
	if !ok {
		t.Fatalf("object is not Hash. got=%T (%+v)", result, result)
	}

	if len(hash.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d, expected=6", len(hash.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair found for given key")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, result, int64(expected))
		case nil:
			testNullObject(t, result)
		}
	}
}

// statements :)
func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.expected)
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

// function testing
func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	result := testEval(input)
	fn, ok := result.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", result, result)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		testIntegerObject(t, result, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = fn(x) {
			fn(y) { x + y };
		};
		let addTwo = newAdder(2);
		addTwo(2);
	`
	testIntegerObject(t, testEval(input), 4)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`len("four")`, 4},
		{`len("")`, 0},
		{`len(1)`, "argument type given to `len` not supported, got=INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, expected=1"},

		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first(1)`, "argument type given to `first` not supported, got=INTEGER"},
		{`first("one", "two")`, "wrong number of arguments. got=2, expected=1"},

		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, "argument type given to `last` not supported, got=INTEGER"},
		{`last("one", "two")`, "wrong number of arguments. got=2, expected=1"},

		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`rest(1)`, "argument type given to `rest` not supported, got=INTEGER"},
		{`rest("one", "two")`, "wrong number of arguments. got=2, expected=1"},

		{`push([1, 2, 3], 4)`, []int{1, 2, 3, 4}},
		{`push(1, 1)`, "argument type given to `push` not supported, got=INTEGER"},
		{`push(1, 1, 1)`, "wrong number of arguments. got=3, expected=2"},

		{`puts("hello world")`, nil},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case nil:
			testNullObject(t, result)

		case int:
			testIntegerObject(t, result, int64(expected))

		case []int:
			array, ok := result.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", result, result)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong number of elements. got=%d, expected=%d", len(array.Elements), len(expected))
				continue
			}

			for i, elem := range expected {
				testIntegerObject(t, array.Elements[i], int64(elem))
			}

		case string:
			errorObj, ok := result.(*object.Error)
			if !ok {
				t.Errorf("object type not expected. got=%T (%+v)", errorObj, errorObj)
			}

			if errorObj.Message != expected {
				t.Errorf("wrong error message. got=%q, expected=%q", errorObj.Message, expected)
			}
		}
	}
}

// errors :)
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable type given as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		result := testEval(tt.input)

		err, ok := result.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", result, result)
		}

		if err.Message != tt.expected {
			t.Errorf("wrong error message. got=%q, expected=%q", err.Message, tt.expected)
		}
	}
}

// generic funcs :)
func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong Value. got=%d, expected=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong Value. got=%t, expected=%t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
