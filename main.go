package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strconv"
)

// The application.
var app = tview.NewApplication()
var types =  []string{"専門科目", "人文・社会科学及び英語科目群", "数学・自然科学・情報技術系科目群"}

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
	syllabuses, err := parseSyllabuses()
	if err != nil {
		log.Fatal(err)
	}
	for _, syllabus := range syllabuses{
		table := CreateTable(syllabus)
		app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEsc {
				app.Stop()
				return nil
			}
			if event.Key() == tcell.KeyCtrlC {
				app.Stop()
				return nil
			}
			return event
		})
		if err := app.SetRoot(table, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
	}


	fmt.Println("success fully finished")
}
func parseSyllabuses() (map[string][]Syllabus, error) {
	syllabuses := make(map[string][]Syllabus, 3)
	for _, t := range types {
		buf, err := ioutil.ReadFile(fmt.Sprintf("./original-syllabus/%s.yaml", t))
		if err != nil {
			return nil, err
		}

		syllabus := []Syllabus{}
		if err := yaml.Unmarshal(buf, &syllabus); err != nil {
			return nil, err
		}
		syllabuses[t] = syllabus
	}

	//for _, syllabus := range syllabuses["専門科目"] {
	//	fmt.Printf("--- %s:\n%v\n\n", syllabus.Group, syllabus)
	//}
	return syllabuses, nil
}


func CreateTable(syllabuses []Syllabus) tview.Primitive{
	var tables  []*tview.Table
	fmt.Println(syllabuses)
	for _, syllabus := range syllabuses {
		table := tview.NewTable().
			SetFixed(1, 1).
			SetSelectable(true, false)
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
		for i, subject := range syllabus.Subjects {
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
			if key == tcell.KeyEnter {
				table.SetSelectable(true, true)
			}
		}).SetSelectedFunc(func(row int, column int) {
			if syllabus.Subjects[row-1].EarnCredit {
				table.GetCell(row, 4).SetTextColor(tcell.ColorWhite).SetText("false")

			} else {
				table.GetCell(row, 4).SetTextColor(tcell.ColorRed).SetText("true")

			}
			syllabus.Subjects[row-1].EarnCredit = !syllabus.Subjects[row-1].EarnCredit
			table.SetSelectable(true, false)
		})
		table.
			SetBorder(true).
			SetTitle(syllabus.Group)
		tables = append(tables, table)
	}

	flex := tview.NewFlex()
	var rows []*tview.Flex
	row := tview.NewFlex().SetDirection(tview.FlexRow)
	for i, table := range tables {
		row.AddItem(table, 0, 1, false)
		if i % 3 == 2 {
			rows = append(rows, row)
			row = tview.NewFlex().SetDirection(tview.FlexRow)
		}
	}
	if len(rows) == 0 {
		rows = append(rows, row)
	}
	fmt.Println(rows)

	for _, row := range rows {
		flex.AddItem(row, 0, 2, true)
	}

	return 	flex
}