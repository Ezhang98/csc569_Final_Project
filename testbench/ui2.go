package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

type UIWindow struct {
	TrainData  string
	TestData   string
	ModelCount int
	Models     []ModelConfig
}

type ModelConfig struct {
	ModelID int
	Name    string
	Model1Params
	Model2Params
	Model3Params
}

type Model1Params struct {
	Activation int
	Nodes      int
}

type Model2Params struct {
	Layers       int
	LearningRate string
}

type Model3Params struct {
	Trees    int
	MaxDepth int
}

var mainwin *ui.Window
var modelCount int = 0
var models []ui.Control = make([]ui.Control, 5)
var windowData UIWindow

func makeModelParam(m ModelConfig) ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	grid := ui.NewGrid()
	grid.SetPadded(true)

	hbox.Append(grid, false)

	form := ui.NewForm()
	form.SetPadded(true)
	grid.Append(form,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	initial := 0
	if m.Name == "Model 2" {
		initial = 1
	}
	if m.Name == "Model 3" {
		initial = 2
	}
	alignment := ui.NewCombobox()
	// note that the items match with the values of the uiDrawTextAlign values
	alignment.Append("Model 1")
	alignment.Append("Model 2")
	alignment.Append("Model 3")
	alignment.SetSelected(initial)
	alignment.OnSelected(func(*ui.Combobox) {
		s := alignment.Selected()
		if s == 0 {
			windowData.Models[m.ModelID].Name = "Model 1"
			hbox1 := ui.NewHorizontalBox()
			hbox1.SetPadded(true)
			grid.Append(hbox1,
				0, 1, 1, 1,
				false, ui.AlignFill, false, ui.AlignFill)

			activation := ui.NewCombobox()
			// note that the items match with the values of the uiDrawTextAlign values
			activation.Append("tanh")
			activation.Append("sigmoid")
			activation.Append("relu")
			activation.SetSelected(0)
			activation.OnSelected(func(*ui.Combobox) {
				windowData.Models[m.ModelID].Activation = activation.Selected()
			})
			// activation and nodes
			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("activation", activation, false)

			nodes := ui.NewSpinbox(0, 100)
			nodes.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].Nodes = nodes.Value()
			})

			form2 := ui.NewForm()
			form2.SetPadded(true)
			hbox1.Append(form2, false)

			form2.Append("numNodes", nodes, false)
		} else if s == 1 {
			windowData.Models[m.ModelID].Name = "Model 2"
			hbox1 := ui.NewHorizontalBox()
			hbox1.SetPadded(true)
			grid.Append(hbox1,
				0, 1, 1, 1,
				false, ui.AlignFill, false, ui.AlignFill)
			// layers and learning rate

			layers := ui.NewSpinbox(0, 100)
			layers.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].Layers = layers.Value()
			})

			form2 := ui.NewForm()
			form2.SetPadded(true)
			hbox1.Append(form2, false)
			form2.Append("numLayers", layers, false)

			lrate := ui.NewEntry()
			lrate.SetText("0.1")
			windowData.Models[m.ModelID].LearningRate = "0.1"
			lrate.OnChanged(func(*ui.Entry) {
				windowData.Models[m.ModelID].LearningRate = lrate.Text()
			})

			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("learning rate", lrate, false)
		} else {
			windowData.Models[m.ModelID].Name = "Model 3"
			hbox1 := ui.NewHorizontalBox()
			hbox1.SetPadded(true)
			grid.Append(hbox1,
				0, 1, 1, 1,
				false, ui.AlignFill, false, ui.AlignFill)
			// numTrees and max depth
			numTrees := ui.NewSpinbox(0, 100)
			numTrees.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].Trees = numTrees.Value()
			})

			form1 := ui.NewForm()
			form1.SetPadded(true)
			hbox1.Append(form1, false)
			form1.Append("numTrees", numTrees, false)

			maxDepth := ui.NewSpinbox(0, 100)
			maxDepth.OnChanged(func(*ui.Spinbox) {
				windowData.Models[m.ModelID].MaxDepth = maxDepth.Value()
			})

			form2 := ui.NewForm()
			form2.SetPadded(true)
			hbox1.Append(form2, false)
			form2.Append("maxDepth", maxDepth, false)
		}
	})

	if m.Name == "Model 1" {
		hbox1 := ui.NewHorizontalBox()
		hbox1.SetPadded(true)
		grid.Append(hbox1,
			0, 1, 1, 1,
			false, ui.AlignFill, false, ui.AlignFill)

		activation := ui.NewCombobox()
		// note that the items match with the values of the uiDrawTextAlign values
		activation.Append("tanh")
		activation.Append("sigmoid")
		activation.Append("relu")
		activation.SetSelected(windowData.Models[m.ModelID].Activation)
		activation.OnSelected(func(*ui.Combobox) {
			windowData.Models[m.ModelID].Activation = activation.Selected()
		})
		// activation and nodes
		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("activation", activation, false)

		nodes := ui.NewSpinbox(0, 100)
		nodes.SetValue(windowData.Models[m.ModelID].Nodes)
		nodes.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].Nodes = nodes.Value()
		})
		form2 := ui.NewForm()
		form2.SetPadded(true)
		hbox1.Append(form2, false)
		form2.Append("numNodes", nodes, false)
	} else if m.Name == "Model 2" {
		hbox1 := ui.NewHorizontalBox()
		hbox1.SetPadded(true)
		grid.Append(hbox1,
			0, 1, 1, 1,
			false, ui.AlignFill, false, ui.AlignFill)
		// layers and learning rate
		layers := ui.NewSpinbox(0, 100)
		layers.SetValue(windowData.Models[m.ModelID].Layers)
		layers.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].Layers = layers.Value()
		})

		form2 := ui.NewForm()
		form2.SetPadded(true)
		hbox1.Append(form2, false)
		form2.Append("numLayers", layers, false)

		lrate := ui.NewEntry()
		lrate.SetText(windowData.Models[m.ModelID].LearningRate)
		lrate.OnChanged(func(*ui.Entry) {
			windowData.Models[m.ModelID].LearningRate = lrate.Text()
		})

		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("learning rate", lrate, false)
	} else if m.Name == "Model 3" {
		hbox1 := ui.NewHorizontalBox()
		hbox1.SetPadded(true)
		grid.Append(hbox1,
			0, 1, 1, 1,
			false, ui.AlignFill, false, ui.AlignFill)
		// numTrees and max depth
		numTrees := ui.NewSpinbox(0, 100)
		numTrees.SetValue(windowData.Models[m.ModelID].Trees)
		numTrees.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].Trees = numTrees.Value()
		})

		form1 := ui.NewForm()
		form1.SetPadded(true)
		hbox1.Append(form1, false)
		form1.Append("numTrees", numTrees, false)

		maxDepth := ui.NewSpinbox(0, 100)
		maxDepth.SetValue(windowData.Models[m.ModelID].MaxDepth)
		maxDepth.OnChanged(func(*ui.Spinbox) {
			windowData.Models[m.ModelID].MaxDepth = maxDepth.Value()
		})

		form2 := ui.NewForm()
		form2.SetPadded(true)
		hbox1.Append(form2, false)
		form2.Append("maxDepth", maxDepth, false)
	}

	form.Append("Alignment", alignment, false)

	return hbox
}

func generateFromState() *ui.Grid {
	grid := ui.NewGrid()
	grid.SetPadded(true)

	button := ui.NewButton("Training Data")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	entry.SetText(windowData.TrainData)
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry.SetText(filename)
		windowData.TrainData = filename
	})
	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry,
		1, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	button1 := ui.NewButton("Test Data")
	entry1 := ui.NewEntry()
	entry1.SetReadOnly(true)
	entry1.SetText(windowData.TestData)
	button1.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry1.SetText(filename)
		windowData.TestData = filename
	})
	grid.Append(button1,
		0, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry1,
		1, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)
	modelCount = windowData.ModelCount
	for i := 0; i < windowData.ModelCount; i++ {
		m := makeModelParam(windowData.Models[i])
		grid.Append(m,
			0, i+2, 2, 1,
			true, ui.AlignFill, false, ui.AlignFill)
	}
	return grid
}

func makeToolbar2() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	msggrid := ui.NewGrid()
	msggrid.SetPadded(true)
	vbox.Append(msggrid, false)

	grid := ui.NewGrid()
	grid.SetPadded(true)
	vbox.Append(grid, false)

	button := ui.NewButton("Import Config")
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename != "" {
			file, _ := ioutil.ReadFile(filename)
			temp := UIWindow{}
			_ = json.Unmarshal([]byte(file), &temp)
			windowData = temp
			vbox.Delete(1)
			grid = generateFromState()
			vbox.Append(grid, false)
		}
	})
	msggrid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	button = ui.NewButton("Add Model")
	button.OnClicked(func(*ui.Button) {
		if modelCount < 5 {
			var m ModelConfig
			m.Name = ""
			m.ModelID = modelCount
			windowData.Models[modelCount] = m
			model := makeModelParam(m)
			grid.Append(model,
				0, modelCount+2, 2, 1,
				true, ui.AlignFill, false, ui.AlignFill)
			modelCount++
			windowData.ModelCount++
		}
	})
	msggrid.Append(button,
		1, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Save Config")
	button.OnClicked(func(*ui.Button) {
		filename := ui.SaveFile(mainwin)
		file, _ := json.MarshalIndent(windowData, "", " ")
		_ = ioutil.WriteFile(filename, file, 0644)
	})
	msggrid.Append(button,
		2, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Clear All")
	button.OnClicked(func(*ui.Button) {
		temp := UIWindow{}
		windowData = temp
		windowData.Models = make([]ModelConfig, 5)
		windowData.ModelCount = 0
		vbox.Delete(1)
		grid = generateFromState()
		vbox.Append(grid, false)
	})
	msggrid.Append(button,
		3, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Run Models")
	button.OnClicked(func(*ui.Button) {
		ui.MsgBoxError(mainwin,
			"This message box describes an error.",
			"More detailed information can be shown here.")
	})
	msggrid.Append(button,
		4, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)

	button = ui.NewButton("Training Data")
	entry := ui.NewEntry()
	entry.SetReadOnly(true)
	button.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry.SetText(filename)
	})
	grid.Append(button,
		0, 0, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry,
		1, 0, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	button1 := ui.NewButton("Test Data")
	entry1 := ui.NewEntry()
	entry1.SetReadOnly(true)
	button1.OnClicked(func(*ui.Button) {
		filename := ui.OpenFile(mainwin)
		if filename == "" {
			filename = "(cancelled)"
		}
		entry1.SetText(filename)
	})
	grid.Append(button1,
		0, 1, 1, 1,
		false, ui.AlignFill, false, ui.AlignFill)
	grid.Append(entry1,
		1, 1, 1, 1,
		true, ui.AlignFill, false, ui.AlignFill)

	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("libui Control Gallery", 640, 480, true)
	windowData.Models = make([]ModelConfig, 5)
	windowData.ModelCount = 0
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})
	mainwin.SetMargined(true)
	mainwin.SetChild(makeToolbar2())

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}
