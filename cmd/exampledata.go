package main

import (
	"foodbar/views/viewutils"

	"github.com/google/uuid"
)



func mkRecipeItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Recipe,
		ItemID:   uuid.New(),
	}
}

func mkPantryItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Pantry,
		ItemID:   uuid.New(),
	}
}

func mkMenuItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Menu,
		ItemID:   uuid.New(),
	}
}

func mkShoppingItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Shopping,
		ItemID:   uuid.New(),
	}
}

func mkPreplistItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Preplist,
		ItemID:   uuid.New(),
	}
}

func mkEarningsItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Earnings,
		ItemID:   uuid.New(),
	}
}

func newExampleRecipeTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			mkRecipeItem("Chicken"),
			mkRecipeItem("turd sandwich"),
			mkRecipeItem("chicken masala"),
			mkRecipeItem("tacos caliente"),
		},
		Ttype: viewutils.Recipe,
	}
}

func newExamplePantryTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			mkPantryItem("Chicken"),
			mkPantryItem("turd sandwich"),
			mkPantryItem("chicken masala"),
			mkPantryItem("tacos caliente"),
		},
		Ttype: viewutils.Pantry,
	}
}

func newExampleMenuTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			mkMenuItem("Chicken"),
			mkMenuItem("turd sandwich"),
			mkMenuItem("chicken masala"),
			mkMenuItem("tacos caliente"),
		},
		Ttype: viewutils.Menu,
	}
}

func newExampleShoppingTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			mkShoppingItem("Chicken"),
			mkShoppingItem("turd sandwich"),
			mkShoppingItem("chicken masala"),
			mkShoppingItem("tacos caliente"),
		},
		Ttype: viewutils.Shopping,
	}
}

func newExamplePreplistTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			mkPreplistItem("Chicken"),
			mkPreplistItem("turd sandwich"),
			mkPreplistItem("chicken masala"),
			mkPreplistItem("tacos caliente"),
		},
		Ttype: viewutils.Preplist,
	}
}

func newExampleEarningsTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			mkEarningsItem("Chicken"),
			mkEarningsItem("turd sandwich"),
			mkEarningsItem("chicken masala"),
			mkEarningsItem("tacos caliente"),
		},
		Ttype: viewutils.Earnings,
	}
}
