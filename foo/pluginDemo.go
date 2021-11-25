//package必须是main
package foo


//为了演示完整与gengine无关的go plugin使用实现，需要一个接口
type Man interface {
	SaveLive() error
}

type SuperMan struct {

}

func (g *SuperMan) SaveLive() error {

	println("execute finished...")
	return nil
}

// go build -buildmode=plugin -o=plugin_M_m.so pluginDemo.go

// exported as symbol named "M",必须大写开头
var M = SuperMan{}
