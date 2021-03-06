// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strings"

	"github.com/google/gxui"
	"github.com/google/gxui/drivers/gl"
	"github.com/google/gxui/themes/dark"
)

func Fonts() ([]string, error) {
	b, err := exec.Command("fc-list").Output()
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanLines)

	fonts := make(map[string]struct{})
	for sc.Scan() {
		t := sc.Text()
		s := strings.Split(t, ": ")[0]
		fonts[s] = struct{}{}
	}
	res := make([]string, 0, len(fonts))
	for key, _ := range fonts {
		res = append(res, key)
	}
	sort.Sort(sort.StringSlice(res))
	return res, nil
}

func CreateList(theme gxui.Theme) gxui.List {
	adapter := gxui.CreateDefaultAdapter()
	items, err := Fonts()
	if err != nil {
		panic(err)
	}
	adapter.SetItems(items)

	list := theme.CreateList()
	list.SetAdapter(adapter)
	list.SetOrientation(gxui.Vertical)

	return list
}

func CreateText(theme gxui.Theme) gxui.TextBox {
	text := theme.CreateTextBox()
	text.SetText(`あのイーハトーヴォの
すきとおった風、
夏でも底に冷たさをもつ青いそら、
うつくしい森で飾られたモーリオ市、
郊外のぎらぎらひかる草の波。
祇辻飴葛蛸鯖鰯噌庖箸
ABCDEFGHIJKLM
abcdefghijklm
1234567890`)

	return text
}

func setFont(theme gxui.Theme, driver gxui.Driver, text gxui.TextBox, path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	font, err := driver.CreateFont(b, 32)
	if err != nil {
		return err
	}

	text.SetFont(font)
	return nil
}

func appMain(driver gxui.Driver) {
	theme := dark.CreateTheme(driver)

	window := theme.CreateWindow(1920, 1080, "Hi")

	splitter := theme.CreateSplitterLayout()
	splitter.SetOrientation(gxui.Horizontal)

	list := CreateList(theme)
	text := CreateText(theme)
	list.OnSelectionChanged(func(item gxui.AdapterItem) {
		s, ok := item.(string)
		if !ok {
			log.Printf("%q is not string", item)
			return
		}
		err := setFont(theme, driver, text, s)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("set %s", s)
	})

	splitter.AddChild(list)
	splitter.AddChild(text)

	layout := theme.CreateLinearLayout()
	layout.AddChild(splitter)
	window.AddChild(layout)

	window.OnClose(driver.Terminate)
}

func main() {
	gl.StartDriver(appMain)
}
