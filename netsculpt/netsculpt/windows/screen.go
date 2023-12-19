package windows

import (
	"NetSculpt/netsculpt/core"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

var (
	configData = []string{"Name", "Hosts"}
	offset     = 2
	t          = tview.NewTable()
	hflex      = tview.NewFlex()
	result     = false
)

func startingIpInput() *tview.Flex {
	first := tview.NewInputField().
		SetFieldWidth(inputWidth).
		SetPlaceholder("E.g. 10").
		SetAcceptanceFunc(ipInputValidator)
	second := tview.NewInputField().
		SetFieldWidth(inputWidth).
		SetPlaceholder("E.g. 0").
		SetAcceptanceFunc(ipInputValidator)
	third := tview.NewInputField().
		SetFieldWidth(inputWidth).
		SetPlaceholder("E.g. 0").
		SetAcceptanceFunc(ipInputValidator)
	fourth := tview.NewInputField().
		SetFieldWidth(inputWidth).
		SetPlaceholder("E.g. 0").
		SetAcceptanceFunc(ipInputValidator)
	prefix := tview.NewInputField().
		SetFieldWidth(inputWidth).
		SetPlaceholder("E.g. 8").
		SetAcceptanceFunc(prefixInputValidator)
	hflex.SetTitle("Starting IP (" + core.StartingIP.ToString() + ")")
	hflex.SetBorder(true)
	set := tview.NewButton("Set")
	set.SetSelectedFunc(func() {
		if result {
			return
		}
		a, _ := strconv.Atoi(first.GetText())
		b, _ := strconv.Atoi(second.GetText())
		c, _ := strconv.Atoi(third.GetText())
		d, _ := strconv.Atoi(fourth.GetText())
		e, _ := strconv.Atoi(prefix.GetText())
		core.SetStartingIP([]int{a, b, c, d}, e)
		first.SetText("")
		second.SetText("")
		third.SetText("")
		fourth.SetText("")
		prefix.SetText("")
		hflex.SetTitle("Starting IP (" + core.StartingIP.ToString() + ")")
	})
	hflex.
		AddItem(first, 0, 2, false).
		AddItem(second, 0, 2, false).
		AddItem(third, 0, 2, false).
		AddItem(fourth, 0, 2, false).
		AddItem(newText("/"), 0, 1, true).
		AddItem(prefix, 0, 2, true).
		AddItem(set, 0, 1, true)
	return hflex
}

func controlsField(a *tview.Application) *tview.Flex {
	quit := tview.NewButton("Quit")
	quit.SetSelectedFunc(func() {
		a.Stop()
	})

	reset := tview.NewButton("Reset")
	reset.SetSelectedFunc(func() {
		core.SetStartingIP([]int{0, 0, 0, 0}, 0)
		hflex.SetTitle("Starting IP (" + core.StartingIP.ToString() + ")")
		configData = []string{"Name", "Hosts"}
		offset = 2
		result = false
		setTableData()

	})

	compute := tview.NewButton("Compute")
	compute.SetSelectedFunc(func() {
		if result {
			return
		}
		t.SetSelectable(false, false)
		if core.StartingIP.IsAllZeros() {
			return
		}
		newConfig := []string{"Name", "Network address", "Broadcast address", "Amount of ips", "Required ips", "Message"}
		for i := 2; i < len(configData); i += 2 {
			name := configData[i]
			r, _ := strconv.Atoi(configData[i+1])
			d := core.CalculateSubnet(r)
			newConfig = append(newConfig, name,
				d.NetworkAddress.ToString(),
				d.BroadcastAddress.ToString(),
				fmt.Sprintf("%d", d.Ips),
				configData[i+1],
				d.Message,
			)
		}
		configData = newConfig
		offset = 6
		setTableData()
		result = true

	})
	var hflex = tview.NewFlex()
	hflex.SetTitle("Actions")
	hflex.SetBorder(true)
	hflex.AddItem(quit, 0, 2, false)
	hflex.AddItem(newText(""), 0, 1, false)
	hflex.AddItem(reset, 0, 2, false)
	hflex.AddItem(newText(""), 0, 1, false)
	hflex.AddItem(compute, 0, 2, false)
	return hflex
}

func setTableData() {
	row := 0
	column := 0
	t.Clear()
	for i, v := range configData {
		color := tcell.ColorWhite
		if row == 0 {
			color = tcell.ColorBlue
		}
		t.SetCell(row, column, tview.NewTableCell(v).SetTextColor(color))

		column++
		if (i+1)%offset == 0 {
			row++
			column = 0
		}
	}
}

func overView() *tview.Table {
	t.SetBorders(true)
	t.SetTitle("Data")
	setTableData()
	t.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter && !result {
			t.SetSelectable(true, false)
		}
		if key == tcell.KeyEscape && !result {
			t.SetSelectable(false, false)
		}
	}).SetSelectedFunc(func(row int, column int) {
		if (row == 0 && offset == 2) || result {
			return
		}
		configData = append(configData[:row*2], configData[(row*2)+2:]...)
		t.Select(row-1, column)
		setTableData()
	})
	return t
}

func addInput(a *tview.Application) *tview.Flex {
	hosts := tview.NewInputField().
		SetFieldWidth(inputWidth).
		SetPlaceholder("Required Hosts").
		SetAcceptanceFunc(tview.InputFieldInteger)
	Name := tview.NewInputField().
		SetFieldWidth(inputWidth).
		SetPlaceholder("Name")
	add := tview.NewButton("Add")
	add.SetSelectedFunc(func() {
		if result {
			return
		}
		if len(hosts.GetText()) == 0 || len(Name.GetText()) == 0 {
			return
		}
		configData = append(configData, Name.GetText(), hosts.GetText())
		setTableData()
		hosts.SetText("")
		Name.SetText("")
		//println("added")
		//a.Draw()
		//hostsNum, _ := strconv.Atoi(hosts.GetText())

	})
	var hflex = tview.NewFlex()
	hflex.SetTitle("Config")
	hflex.SetBorder(true)
	hflex.AddItem(Name, 0, 1, false)
	hflex.AddItem(hosts, 0, 1, false)
	hflex.AddItem(add, 0, 1, false)
	return hflex
}

func RunScreen() {
	app := tview.NewApplication()
	var vflex = tview.NewFlex()
	vflex.SetDirection(tview.FlexRow).
		AddItem(startingIpInput(), 0, 1, false).
		AddItem(addInput(app), 0, 1, false).
		AddItem(overView(), 0, 4, false).
		AddItem(controlsField(app), 0, 1, false)

	if err := app.SetRoot(vflex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
