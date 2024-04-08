package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/dominickp/go-hn-cli/client"
)

func padRight(str string, length int) string {
	for {
		str += " "
		if len(str) >= length {
			return str
		}
	}
}

func main() {

	selectedItem := 0

	topMenuResponse, err := client.GetTopMenuResponse(10)
	for _, item := range topMenuResponse.Items {
		fmt.Println(item.Id, item.Title)
	}

	huhOptions := make([]huh.Option[int], len(topMenuResponse.Items))
	for i, item := range topMenuResponse.Items {
		scoreKey := padRight(fmt.Sprintf("%d", item.Score), 4)
		optionTitle := fmt.Sprintf("%s %s", scoreKey, item.Title) // Convert item.Score to string using %d.
		huhOptions[i] = huh.NewOption(optionTitle, item.Id)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Top Topics").
				Options(huhOptions...).
				Value(&selectedItem),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(selectedItem)

	// Generate new form to show the topic and comments, then have a back option
	// That renders the first form again in a loop

}
