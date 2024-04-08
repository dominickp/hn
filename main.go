package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/huh"
	"github.com/dominickp/go-hn-cli/client"
)

func main() {

	sauceLevel := 0

	topMenuResponse, err := client.GetTopMenuResponse(10)
	for _, item := range topMenuResponse.Items {
		fmt.Println(item.Id, item.Title)
	}

	huhOptions := make([]huh.Option[int], len(topMenuResponse.Items))
	for i, item := range topMenuResponse.Items {
		optionTitle := fmt.Sprintf("[%d]\t%s", item.Score, item.Title) // Convert item.Score to string using %d.
		huhOptions[i] = huh.NewOption(optionTitle, item.Id)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Top Topics").
				Options(huhOptions...).
				Value(&sauceLevel),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal(err)
	}

}
