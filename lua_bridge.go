// Copyright (c) 2026 Kartendsy
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.



package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func fromGo(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case float64:
		return lua.LNumber(v)
	case int:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case bool:
		return lua.LBool(v)
	case map[string]interface{}:
		tbl := L.NewTable()
		for key, val := range v {
			tbl.RawSetString(key, fromGo(L, val))
		}
		return tbl
	case []interface{}:
		tbl := L.NewTable()
		for _, val := range v {
			tbl.Append(fromGo(L, val))
		}
		return tbl
	case nil:
		return lua.LNil
	default:
		return lua.LNil
	}
}

func NewSafeState() *lua.LState {
	L := lua.NewState(lua.Options{SkipOpenLibs: true})

	libs := []struct {
		Name string
		Func lua.LGFunction
	}{
		{lua.LoadLibName, lua.OpenPackage},
		{lua.BaseLibName, lua.OpenBase},
		{lua.StringLibName, lua.OpenString},
		{lua.TabLibName, lua.OpenTable},
	}
	for _, lib := range libs {
		L.Push(L.NewFunction(lib.Func))
		L.Push(lua.LString(lib.Name))
		L.Call(1, 0)
	}

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"get": func(l *lua.LState) int {
			url := L.CheckString(1)
			resp, err := http.Get(url)
			if err != nil {
				L.Push(lua.LNil)
				return 1
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			L.Push(lua.LString(string(body)))
			return 1
		},
		"log": func(l *lua.LState) int {
			fmt.Println("[LUA]", l.CheckString(1))
			return 0
		},
		"json_decode": func(l *lua.LState) int {
			str := L.CheckString(1)
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(str), &result); err != nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(fromGo(L, result))
			return 1
		},
		"IsCommand": func(l *lua.LState) int {
			cmd := L.CheckString(1)
			input := L.CheckString(2)

			match := strings.HasPrefix(strings.ToLower(input), strings.ToLower(cmd))
			L.Push(lua.LBool(match))
			return 1
		},
	})
	L.SetGlobal("bot", mod)
	return L
}
