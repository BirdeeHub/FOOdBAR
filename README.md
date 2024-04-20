# This repo is an extreme work in progress, it is NOT EVEN CLOSE TO FINISHED!!!!

## It is my first time makind a webpage, and first time using go, htmx, sql, js, html, css, templ, routing, etc.

This repo will be a tracker app for a small catering buisness when it is finished.

Mostly though, it is a way for me to learn web development as I have not done much of it.


### TODO AND IDEAS:
database efficient infinite scroll for all list style tabs

shopping price search feature?

styling/animation?

recipe list and entry

suggestion based on stock

sort recipe by category, dietary, ingredients, % ingredients in stock, (total ingredient price?)

generate prep list, and shopping list (accounts for current stock), the menu, and projected and actual portions ordered.

track what was bought and then used after the day to keep track of the total stock, as actual amounts may differ based on prices for different quantities at stores.

total workflow, 3 interaction stages. Generate menu, input orders and generate the lists, input actual stock bought and used to correct suggested values.

6 tab interface? recipe database, menu tab where you have the menu and track orders, prep list, shopping list, stock list, profit review

All are list style tabs with expandable elements

Would be nice to have these able to be dynamically resized with the mouse

recipes will need category, dietary, ingredients, estimated cost, fields, ingredients should have an optional scaling factor for scaling recipes accurately

ingredients will need amounts in stock, time since purchased, storage method, (last queried price?)

when all data for stock and menu and orders have been finalized for the week, generate profit report, with option for adding extra incidental costs.

Interface will be written as a web page, because I want to learn how to make an interactive web page, and I will be using go + htmx to do it.

Run the program, and then visit localhost:42069 (or other configured value)

recipe entry should have completion for existing items in ingredients dietary, category, etc but not for name or instructions or amounts or prices.
