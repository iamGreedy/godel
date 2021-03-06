package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/essence/axis"
	"github.com/iamGreedy/essence/meter"
	"github.com/iamGreedy/essence/must"
	"github.com/iamGreedy/essence/prefix"
	"github.com/iamGreedy/gltf2"
	"github.com/iamGreedy/godel/back"
	"github.com/iamGreedy/godel/shader"
	"github.com/iamGreedy/godel"
	"os"
	"runtime"
)

func main() {
	var (
		width  = 640
		height = 480
	)
	runtime.LockOSThread()
	// GLFW, GL Init
	wnd := window(width, height)
	defer glfw.Terminate()
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version : ", version)
	// GLTF
	// Complete
	//f := must.MustGet(os.Open("./demo/models/damagedHelmet/damagedHelmet.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/2CylinderEngine/glTF/2CylinderEngine.gltf")).(*os.File)
	// TODO f := must.MustGet(os.Open("./demo/models/2CylinderEngine/glTF-Draco/2CylinderEngine.gltf")).(*os.File)
	// TODO f := must.MustGet(os.Open("./demo/models/2CylinderEngine/glTF-Draco/2CylinderEngine.gltf")).(*os.File)
	// f := must.MustGet(os.Open("./demo/models/2CylinderEngine/glTF-Embedded/2CylinderEngine.gltf")).(*os.File)
	// TODO f := must.MustGet(os.Open("./demo/models/2CylinderEngine/glTF-pbrSpecularGlossiness/2CylinderEngine.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/AnimatedCube/glTF/AnimatedCube.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/AnimatedMorphCube/glTF/AnimatedMorphCube.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/AnimatedTriangle/glTF/AnimatedTriangle.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/AnimatedTriangle/glTF-Embedded/AnimatedTriangle.gltf")).(*os.File)
	// TODO f := must.MustGet(os.Open("./demo/models/Avocado/glTF-pbrSpecularGlossiness/Avocado.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/Avocado/glTF/Avocado.gltf")).(*os.File)

	//f := must.MustGet(os.Open("./demo/models/Buggy/glTF/Buggy.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/BrainStem/glTF/BrainStem.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/CesiumMilkTruck/glTF/CesiumMilkTruck.gltf")).(*os.File)
	f := must.MustGet(os.Open("./demo/models/RiggedFigure/glTF/RiggedFigure.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/RiggedSimple/glTF/RiggedSimple.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/0gltfTutorial19/file.gltf")).(*os.File)
	defer f.Close()
	md := must.MustGet(gltf2.Parser().
		Reader(f).
		Logger(os.Stdout).
		Tasks(
			//gltf2.Tasks.HelloWorld,
			gltf2.Tasks.Caching,
			gltf2.Tasks.AutoBufferTarget,
			gltf2.Tasks.AccessorMinMax,
			gltf2.Tasks.ModelScale(axis.Y, meter.New(prefix.No, 8)),
			//gltf2.Tasks.ModelAlign(align.Center, align.Center, align.Center),
			//gltf2.Tasks.ModelAlign(align.No, align.No, align.No),
			//gltf2.Tasks.ByeWorld,
		).
		Strictness(gltf2.LEVEL1).
		Parse()).(*gltf2.GLTF)
	// godel Model
	godel.NewNode()
	//app := godel.NewApplication(shader.Standard, shader.Flat, godel.NewCamera(godel.Perspective), godel.NewLighting())
	app := godel.NewApplication(shader.Standard, shader.PBR, godel.NewCamera(godel.Perspective), godel.NewLighting())
	app.Camera.LookFrom(mgl32.Vec3{0, 11, -128})
	app.Camera.LookTo(mgl32.Vec3{0, 0, 0})
	app.Lighting.Global.Direction = mgl32.Vec3{
		0, -1, 4,
	}
	model := must.MustGet(app.BuildModel(md, true)).(*godel.Model)
	//
	obj := model.NewInstance(-1)
	baseRotation := mgl32.AnglesToQuat(mgl32.DegToRad(90),0,mgl32.DegToRad(0), mgl32.XYZ)
	//baseRotation := mgl32.QuatIdent()
	obj.Rotate(baseRotation)
	obj.NewPlayer(0, func(a *godel.Player) {
		a.Loop = true
		a.PlaySpeed = 1
	})
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	//
	for prev, curr := float32(0), float32(glfw.GetTime()); !wnd.ShouldClose(); prev, curr = curr, float32(glfw.GetTime()) {
		glfw.PollEvents()
		// application update
		app.Update(curr - prev)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		obj.Rotate(mgl32.QuatRotate(curr, mgl32.Vec3{0, 1, 0}).Mul(baseRotation))
		obj.Render()
		//
		wnd.SwapBuffers()
	}

}
func window(w, h int) *glfw.Window {
	// glfw3
	must.Must(glfw.Init())
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	// glfw3.window
	window := must.MustGet(glfw.CreateWindow(w, h, "Testing", nil, nil)).(*glfw.Window)
	window.MakeContextCurrent()
	// OpenGL
	must.Must(gl.Init())
	return window
}
