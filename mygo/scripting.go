package main

import (
	"fmt"
	"io/ioutil"

	lua "github.com/yuin/gopher-lua"
)

var L *lua.LState

func luaInit() {
	L = lua.NewState()

	L.SetGlobal("doPrompt", L.NewFunction(doPrompt))
	L.SetGlobal("getPrompt", L.NewFunction(getPrompt))
	L.SetGlobal("getModel", L.NewFunction(getModel))
	L.SetGlobal("saveFile", L.NewFunction(saveFile))
	L.SetGlobal("save_image", L.NewFunction(luaSaveImage))
	L.SetGlobal("encode_image", L.NewFunction(luaEncodeImage))
	registerAPIClient(L)
}

func doPrompt(L *lua.LState) int {
	arg := L.ToString(1)
	result := "Do Prompt: " + arg
	L.Push(lua.LString(result))
	return 1
}

func getPrompt(L *lua.LState) int {
	arg := L.ToString(1)
	result, _ := getPromptTemplate(arg)
	L.Push(lua.LString(result))
	return 1
}

func getModel(L *lua.LState) int {
	arg := L.ToString(1)
	result := getModelFunction(arg)
	L.Push(lua.LString(result))
	return 1
}

func saveFile(L *lua.LState) int {
	content := L.ToString(1)
	filename := L.ToString(2)
	saveToFile(content, filename)
	L.Push(lua.LString("done"))
	return 1
}

func runScript(script string) {
	if err := L.DoString(script); err != nil {
		panic(err)
	}
}

func runCode(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if err := L.DoString(string(data)); err != nil {
		panic(err)
	}
}

func setLuaPromptOutput(output string) {
	L.SetGlobal("output", lua.LString(output))
}

// Initialize the API client with a base URL
func newAPIClient(L *lua.LState) int {
	baseURL := L.ToString(1)
	client := &APIClient{BaseURL: baseURL}
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("APIClient"))
	L.Push(ud)
	return 1
}

// MouseMove Lua binding
func clientMouseMove(L *lua.LState) int {
	ud := L.CheckUserData(1)
	x := L.CheckInt(2)
	y := L.CheckInt(3)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.MouseMove(x, y)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LString(fmt.Sprintf("MouseMove Point: X=%d, Y=%d", point.X, point.Y)))
		return 1
	}
	return 0
}

// Click Lua binding
func clientClick(L *lua.LState) int {
	ud := L.CheckUserData(1)
	x := L.CheckInt(2)
	y := L.CheckInt(3)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.Click(x, y)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LString(fmt.Sprintf("Click Point: X=%d, Y=%d", point.X, point.Y)))
		return 1
	}
	return 0
}

// MoveWheel Lua binding
func clientMoveWheel(L *lua.LState) int {
	ud := L.CheckUserData(1)
	amt := L.CheckInt(2)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.MoveWheel(amt)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LString(fmt.Sprintf("MoveWheel Point: Data=%v", point.Data)))
		return 1
	}
	return 0
}

// SendText Lua binding
func clientSendText(L *lua.LState) int {
	ud := L.CheckUserData(1)
	x := L.CheckInt(2)
	y := L.CheckInt(3)
	text := L.CheckString(4)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.SendText(x, y, text)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LString(fmt.Sprintf("SendText Point: Data=%v", point.Data)))
		return 1
	}
	return 0
}

// Type Lua binding
func clientType(L *lua.LState) int {
	ud := L.CheckUserData(1)
	text := L.CheckString(2)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.Type(text)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LString(fmt.Sprintf("Type Point: Data=%v", point.Data)))
		return 1
	}
	return 0
}

// FindText Lua binding
func clientFindText(L *lua.LState) int {
	ud := L.CheckUserData(1)
	text := L.CheckString(2)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.FindText(text)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LNumber(point.X))
		L.Push(lua.LNumber(point.Y))
		return 2
	}
	return 0
}

// GetText Lua binding
func clientGetText(L *lua.LState) int {
	ud := L.CheckUserData(1)
	x := L.CheckInt(2)
	y := L.CheckInt(3)
	w := L.CheckInt(4)
	h := L.CheckInt(5)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.GetText(x, y, w, h)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LString(point.Data["text"].(string)))
		return 1
	}
	return 0
}

// FindImage Lua binding
func clientFindImage(L *lua.LState) int {
	ud := L.CheckUserData(1)
	img := L.CheckString(2)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.FindImage(img)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LNumber(point.X))
		L.Push(lua.LNumber(point.Y))
		return 2
	}
	return 0
}

// Screenshot Lua binding
func clientScreenshot(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.Screenshot()
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LString((point.Data["img"].(string))))
		return 1
	}
	return 0
}

// GetMouseColor Lua binding
func clientGetMouseColor(L *lua.LState) int {
	ud := L.CheckUserData(1)
	x := L.CheckInt(2)
	y := L.CheckInt(3)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.GetMouseColor(x, y)
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LNumber(point.Data["rgb_r"].(float64)))
		L.Push(lua.LNumber(point.Data["rgb_g"].(float64)))
		L.Push(lua.LNumber(point.Data["rgb_b"].(float64)))
		return 3
	}
	return 0
}

// GetMouse Lua binding
func clientGetMouse(L *lua.LState) int {
	ud := L.CheckUserData(1)
	if client, ok := ud.Value.(*APIClient); ok {
		point, err := client.GetMouse()
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
			return 1
		}
		L.Push(lua.LNumber(point.X))
		L.Push(lua.LNumber(point.Y))
		return 2
	}
	return 0
}

// Register APIClient methods to Lua
func registerAPIClient(L *lua.LState) {
	mt := L.NewTypeMetatable("APIClient")
	L.SetGlobal("APIClient", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newAPIClient))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"mousemove":     clientMouseMove,
		"click":         clientClick,
		"movewheel":     clientMoveWheel,
		"sendtext":      clientSendText,
		"type":          clientType,
		"findtext":      clientFindText,
		"gettext":       clientGetText,
		"findimage":     clientFindImage,
		"screenshot":    clientScreenshot,
		"getmousecolor": clientGetMouseColor,
		"getmouse":      clientGetMouse,
	}))
}
func luaSaveImage(L *lua.LState) int {
	base64Image := L.CheckString(1)
	filePath := L.CheckString(2)

	if err := SaveBase64ImageToPNG(base64Image, filePath); err != nil {
		L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
		return 1
	}

	L.Push(lua.LString(fmt.Sprintf("Image saved successfully: %s", filePath)))
	return 1
}

func luaEncodeImage(L *lua.LState) int {
	filePath := L.CheckString(1)

	base64Data, err := readFileAsBase64(filePath)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("error: %v", err)))
		return 1
	}

	L.Push(lua.LString(base64Data))
	return 1
}
