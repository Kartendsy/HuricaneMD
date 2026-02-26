# HuricaneMD

A modular WhatsApp bot engine built with Go. Use Lua scripts as plugins to add features without recompiling the entire binary. Stop wasting time on restarts; just save your `.lua` file and you're good to go.

## Why this?
Developing bots in Go is fast, but waiting for compile-build cycles every time you change a small response logic is a pain. This project bridges **Whatsmeow** with **Gopher-Lua**, allowing you to keep the core performance of Go while having the flexibility of a scripting language.

## Key Highlights
* **Modular Architecture**: Every feature lives in `./plugins/`.
* **Hot Reloading**: Scripts are executed on-the-fly. Update your logic while the bot is online.
* **Sandboxed VM**: Uses a restricted Lua environment for safety. No `os.execute` or `io` access by default.
* **Thread-Safe**: Uses an LState pool to handle concurrent WhatsApp messages efficiently.

## Core Stack
* [Whatsmeow](https://github.com/tulir/whatsmeow) - WhatsApp MD library.
* [Gopher-Lua](https://github.com/yuin/gopher-lua) - Lua 5.1 VM in Go.

## Getting Started

1.  **Clone & Install**
    ```bash
    git clone [https://github.com/kartendsy/HuricaneMD.git](https://github.com/Kartendsy/HuricaneMD.git)
    cd repo
    go mod tidy
    ```

2.  **Configuration**
    Populate your `./plugins` folder with your custom logic. There's an `example.lua` to get you started.

3.  **Run**
    ```bash
    go run .
    ```
    Scan the QR code, and you're officially live.

## Plugin API
Every script in the `/plugins` folder has access to:
- `Message`: (string) The incoming message text.
- `Sender`: (string) JID of the sender.
- `Reply`: (string) Set this variable to send a response.
- `bot.get(url)`: Simple HTTP GET helper.
- `bot.json_decode(str)`: Converts JSON string to a Lua table.

### Example: `plugins/ping.lua`
```lua
if Message == ".ping" then
    Reply = "Pong! Bot is breathing."
end
