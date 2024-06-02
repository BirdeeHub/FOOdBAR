# This repo is an extreme work in progress, it is NOT EVEN CLOSE TO FINISHED!!!!

## This is my first attempt at making a webpage.
## It is also my first time using go, htmx, sql, js, html, css, templ, authentication, routing, etc.

This program will be a tracker app for a friend's small catering/personal chef buisness when it is finished.

Mostly though, it is a way for me to learn web development as I have not done any before.

It is a single page application with a monolitic architecture,
designed equally to be hosted as a web page,
or as a completely offline local application viewed through a browser.

There will likely be some spaghetti that I will figure out how to handle better as I go.

### Install and development building instructions

#### nix

For now it is packaged for nix only so install nix package manager on linux, mac, or WSL, by following the instructions here [nixos.org/download](https://nixos.org/download/)

Then to build or run use

```nix build github:BirdeeHub/FOOdBAR```

or ```nix run github:BirdeeHub/FOOdBAR```

or ```nix shell github:BirdeeHub/FOOdBAR```

and to hack around with it with hot reload run the following 4 commands:

```bash
git clone https://github.com/BirdeeHub/FOOdBAR
cd FOOdBAR
nix develop .
# then inside the shell:
air
```
Then go to localhost:42069 in a web browser

#### nixless

If you dont want to install nix, the following should work, and the air command might work.

Assuming you have go, templ, sqlite3, tailwindcss, and air installed already and in your PATH, run the following commands:
```bash
git clone https://github.com/BirdeeHub/FOOdBAR
cd FOOdBAR
go mod tidy && templ generate && tailwindcss -o ./static/tailwind.css -c ./tailwind.config.js && go build -o bin/FOOdBAR main.go
```
-OR-
```bash
git clone https://github.com/BirdeeHub/FOOdBAR
cd FOOdBAR
air
```

### Name:

It's FU'd Beyond Any Recognition

The FOOdBAR tracks the food in the FOOdb and tells you how much FOOd you have and sold and what it cost,
and reminds you of recipes that fit those criteria best. Allowing you to cook more variety with less headache

I also missed the chance to put dreambird in the name so I technically do not own this repo.

---

### Done so far:

setup go, htmx, templ, echo project

made it be buildable by nix, and included a dev shell

defined the basic api routes and templates needed for the tab skeleton

set up authentication via JWT and auth database

defined skeleton of runtime state type objects

defined part of the skeleton for the tabs themselves

---

### TODO AND IDEAS:

make the tab flip to edit work nicely

expandable elements for tab items

Would be nice to have dynamic resizing with the mouse

infinite scroll for all list style tabs

write some damn tests

make it forced http*S*

suggestion based on stock

sort recipe by category, dietary, ingredients, % ingredients in stock, (total ingredient price?)

generate prep list, and shopping list (accounts for current stock), the menu, and projected and actual portions ordered.

track what was bought and then used after the day to keep track of the total stock, as actual amounts may differ based on prices for different quantities at stores.

total workflow, 3 interaction stages. Generate menu, input orders and generate the lists, input actual stock bought and used to correct suggested values.

recipes will need category, dietary, ingredients, estimated cost, fields, ingredients should have an optional scaling factor for scaling recipes accurately

ingredients will need amounts in stock, time since purchased, storage method, (last queried price?)

when all data for stock and menu and orders have been finalized for the week, generate profit report, with option for adding extra incidental costs.

styling/animation?

recipe entry should have completion for existing items in ingredients dietary, category, etc but not for name or instructions or amounts or prices.

shopping price search feature?
