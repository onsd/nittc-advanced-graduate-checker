package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Syllabus struct {
	Group    string    `yaml:"group"`
	Subjects []Subject `yaml:"subjects"`
}
type Subject struct {
	Name       string `yaml:"name"`
	Credits    int    `yaml:"credits"`
	Required   bool   `yaml:"required,omitempty"`
	JABEE      bool   `yaml:"JABEE,omitempty"`
	EarnCredit bool   `yaml:"earn_credit,omitempty"`
}

func main() {
	types := []string{"専門科目", "人文・社会科学及び英語科目群", "数学・自然科学・情報技術系科目群"}
	syllabuses := make(map[string][]Syllabus, 3)
	for _, t := range types {
		buf, err := ioutil.ReadFile(fmt.Sprintf("./original-syllabus/%s.yaml", t))
		if err != nil {
			fmt.Println(err)
			return
		}
		syllabus := []Syllabus{}
		if err := yaml.Unmarshal(buf, &syllabus); err != nil {
			log.Fatalf("error: %v", err)
		}
		syllabuses[t] = syllabus
	}

	for key, syllabus := range syllabuses {
		fmt.Printf("--- %s:\n%v\n\n", key, syllabus)
	}

	app := tview.NewApplication()
	table := tview.NewTable().
		SetFixed(1, 1).SetSelectable(true, false)

	syllabus := syllabuses["専門科目"]
	table.SetCell(0, 0,
		tview.NewTableCell("Name").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter),
	)
	table.SetCell(0, 1,
		tview.NewTableCell("Credits").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter),
	)
	table.SetCell(0, 2,
		tview.NewTableCell("Require ").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter),
	)
	table.SetCell(0, 3,
		tview.NewTableCell("JABEE Required").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter),
	)
	table.SetCell(0, 4,
		tview.NewTableCell("Earned Credit").
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter),
	)
	for i, subject := range syllabus[0].Subjects {
		table.SetCell(i+1, 0,
			tview.NewTableCell(subject.Name).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter),
		)
		table.SetCell(i+1, 1,
			tview.NewTableCell(strconv.Itoa(subject.Credits)).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter),
		)
		table.SetCell(i+1, 2,
			tview.NewTableCell(strconv.FormatBool(subject.Required)).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter),
		)
		table.SetCell(i+1, 3,
			tview.NewTableCell(strconv.FormatBool(subject.JABEE)).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter),
		)
		table.SetCell(i+1, 4,
			tview.NewTableCell(strconv.FormatBool(subject.EarnCredit)).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter),
		)
	}
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
		}
	}).SetSelectedFunc(func(row int, column int) {
		if syllabus[0].Subjects[row-1].EarnCredit {
			table.GetCell(row, 4).SetTextColor(tcell.ColorWhite).SetText("false")

		} else {
			table.GetCell(row, 4).SetTextColor(tcell.ColorRed).SetText("true")

		}
		syllabus[0].Subjects[row-1].EarnCredit = !syllabus[0].Subjects[row-1].EarnCredit
		table.SetSelectable(true, false)
	})
	if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

