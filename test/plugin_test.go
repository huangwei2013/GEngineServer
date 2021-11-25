package test

import (
	"fmt"
	"os"
	"plugin"
	"testing"
)

func Test_plugin(t *testing.T) {

	dir, err := os.Getwd()
	if err!=nil {
		panic(err)
	}

	// load module 插件您也可以使用go http.Request从远程下载到本地,在加载做到动态的执行不同的功能
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(dir + "/../plugins/plugin_M_m.so")
	if err != nil {
		panic(err)
	}
	println("plugin opened")

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Greeter
	m, err := plug.Lookup("M") //大写
	if err != nil {
		panic(err)
	}
	fmt.Println(m)

	//// 3. Assert that loaded symbol is of a desired type
	//man, ok := m.(Man)
	//if !ok {
	//	fmt.Println("unexpected type from module symbol")
	//	os.Exit(1)
	//}
	//
	//// 4. use the module
	//if err := man.SaveLive(); err != nil {
	//	println("use plugin man failed, ", err)
	//}
}

