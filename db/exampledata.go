package db

import (
	"foodbar/views/viewutils"

	"github.com/google/uuid"
)



func MkRecipeItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Recipe,
		ItemID:   uuid.New(),
	}
}

func MkPantryItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Pantry,
		ItemID:   uuid.New(),
	}
}

func MkMenuItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Menu,
		ItemID:   uuid.New(),
	}
}

func MkShoppingItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Shopping,
		ItemID:   uuid.New(),
	}
}

func MkPreplistItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Preplist,
		ItemID:   uuid.New(),
	}
}

func MkEarningsItem(name string) viewutils.TabItem {
	return viewutils.TabItem{
		ItemName: name,
		Ttype:    viewutils.Earnings,
		ItemID:   uuid.New(),
	}
}

func NewExampleRecipeTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			MkRecipeItem("Chicken"),
			MkRecipeItem("turd sandwich"),
			MkRecipeItem("chicken masala"),
			MkRecipeItem("tacos caliente"),
		},
		Ttype: viewutils.Recipe,
	}
}

func NewExamplePantryTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			MkPantryItem("Chicken"),
			MkPantryItem("turd sandwich"),
			MkPantryItem("chicken masala"),
			MkPantryItem("tacos caliente"),
		},
		Ttype: viewutils.Pantry,
	}
}

func NewExampleMenuTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			MkMenuItem("Chicken"),
			MkMenuItem("turd sandwich"),
			MkMenuItem("chicken masala"),
			MkMenuItem("tacos caliente"),
		},
		Ttype: viewutils.Menu,
	}
}

func NewExampleShoppingTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			MkShoppingItem("Chicken"),
			MkShoppingItem("turd sandwich"),
			MkShoppingItem("chicken masala"),
			MkShoppingItem("tacos caliente"),
		},
		Ttype: viewutils.Shopping,
	}
}

func NewExamplePreplistTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			MkPreplistItem("Chicken"),
			MkPreplistItem("turd sandwich"),
			MkPreplistItem("chicken masala"),
			MkPreplistItem("tacos caliente"),
		},
		Ttype: viewutils.Preplist,
	}
}

func NewExampleEarningsTabData(isActive bool) viewutils.TabData {
	return viewutils.TabData{
		Active:   isActive,
		Items: []viewutils.TabItem{
			MkEarningsItem("Chicken"),
			MkEarningsItem("turd sandwich"),
			MkEarningsItem("chicken masala"),
			MkEarningsItem("tacos caliente"),
		},
		Ttype: viewutils.Earnings,
	}
}
