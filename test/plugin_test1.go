package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bilibili/gengine/context"

)

func Test_plugin_with_gengine(t *testing.T)  {

	dir, err := os.Getwd()
	if err!=nil {
		panic(err)
	}

	dc := context.NewDataContext()
	//3.load plugin into apiName, exportApi
	_, _, e := dc.PluginLoader( dir + "/plugin_M_m.so")
	if e != nil {
		panic(e)
	}

	dc.Add("println", fmt.Println)
	ruleBuilder := builder.NewRuleBuilder(dc)
	err = ruleBuilder.BuildRuleFromString(`
	rule "1"
		begin
	
			//this method is defined in plugin
			err = m.SaveLive()
		
			if isNil(err) {
			   println("err is nil")
			}
		end
	`)

	if err != nil {
		panic(err)
	}
	gengine := engine.NewGengine()
	err = gengine.Execute(ruleBuilder, false)

	if err!=nil {
		panic(err)
	}
}
