features

TODO:
maximize button for tabs
efficient infinite scroll for all list style tabs
database
all the tabs


recipe list and entry
suggestion based on stock
sort recipe by category, dietary, ingredients, % ingredients in stock, (total ingredient price?)
finalize menu.

generate prep list, and shopping list (accounts for current stock), the menu, and projected and actual portions ordered.

track what was bought and then used after the day to keep track of the total stock, as actual amounts may differ based on prices for different quantities at stores.

total workflow, 3 interaction stages. Generate menu, input orders and generate the lists, input actual stock bought and used to correct suggested values.

6 tab interface? recipe database, menu tab where you have the menu and track orders, prep list, shopping list, stock list, profit review

Would be nice to have these able to be dynamically resized and put side by side or above and below.

persistent data consists of recipes, and stock, and configuration.

Will need to store recipe objects, and ingredient objects.

recipes will need category dietary, ingredients, last price sold for, fields

ingredients will need amounts in stock, time since purchased, storage method, (last queried price?)

when all data for stock and menu and orders have been finalized for the week, generate profit report, with option for adding extra incidental costs.


Interface will be written as a web page, because I want to learn how to make an interactive web page, and I will be using go + htmx to do it.

Run the program, and then visit localhost:42069 (or other configured value)

recipe entry should have completion for existing items in ingredients dietary, category, etc but not for name or instructions or amounts or prices.
