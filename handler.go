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
	"fmt"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func EventHandler(client *whatsmeow.Client, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.IsFromMe {
			return
		}
		L := pool.Get()
		defer pool.Put(L)

		msgText := v.Message.GetConversation()
		L.SetGlobal("Message", lua.LString(msgText))
		L.SetGlobal("Sender", lua.LString(v.Info.Chat.String()))
		L.SetGlobal("Reply", lua.LNil)

		files, _ := filepath.Glob("./plugins/*.lua")
		for _, f := range files {
			_ = L.DoFile(f)

			res := L.GetGlobal("Reply")
			if res.Type() == lua.LTString {
				client.SendMessage(ctx, v.Info.Chat, &waE2E.Message{
					Conversation: proto.String(res.String()),
				})
				break
			}
		}
		fmt.Println("Dapat pesan", v.Message.GetConversation())
	}
}
