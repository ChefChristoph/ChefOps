#!/bin/sh
set -e

echo "🚨 Resetting ChefOps database..."
rm -rf db
mkdir -p db

echo "🔧 Recreating schema and views..."
sqlite3 db/chefops.db <schema.sql
sqlite3 db/chefops.db <views.sql

echo "🍏 Importing ingredients (with UAE estimated prices)..."

./chefops ingredient add --name "Lobster Meat" --unit kg --cost 150
./chefops ingredient add --name "Hammour" --unit kg --cost 38
./chefops ingredient add --name "Beef Patty" --unit kg --cost 32
./chefops ingredient add --name "Chicken Strips" --unit kg --cost 22

./chefops ingredient add --name "Pasta" --unit kg --cost 6
./chefops ingredient add --name "Breadcrumbs" --unit kg --cost 8
./chefops ingredient add --name "Flour" --unit kg --cost 4
./chefops ingredient add --name "Cornstarch" --unit kg --cost 6

./chefops ingredient add --name "Pizza Base" --unit piece --cost 2.5
./chefops ingredient add --name "Brioche Bun" --unit piece --cost 2.0
./chefops ingredient add --name "Burger Bun" --unit piece --cost 1.5

./chefops ingredient add --name "Cheddar Cheese" --unit kg --cost 26
./chefops ingredient add --name "Mozzarella Cheese" --unit kg --cost 22
./chefops ingredient add --name "Grana Padano Cheese" --unit kg --cost 48

./chefops ingredient add --name "Butter" --unit kg --cost 22
./chefops ingredient add --name "Milk" --unit kg --cost 4

./chefops ingredient add --name "Egg Yolks" --unit piece --cost 0.45
./chefops ingredient add --name "Whole Eggs" --unit piece --cost 0.38

./chefops ingredient add --name "Mayonnaise" --unit kg --cost 12
./chefops ingredient add --name "Dijon Mustard" --unit kg --cost 16

./chefops ingredient add --name "Chives" --unit kg --cost 65
./chefops ingredient add --name "Dill" --unit kg --cost 48
./chefops ingredient add --name "Basil" --unit kg --cost 42

./chefops ingredient add --name "Lemon" --unit kg --cost 7
./chefops ingredient add --name "Lemon Juice" --unit liter --cost 8
./chefops ingredient add --name "Lemon Zest" --unit kg --cost 25

./chefops ingredient add --name "Olive Oil" --unit liter --cost 18
./chefops ingredient add --name "Extra Virgin Olive Oil" --unit liter --cost 22
./chefops ingredient add --name "Grapeseed Oil" --unit liter --cost 15
./chefops ingredient add --name "Sunflower Oil" --unit liter --cost 6

./chefops ingredient add --name "Ice" --unit kg --cost 0.5
./chefops ingredient add --name "Pepper" --unit kg --cost 32
./chefops ingredient add --name "Red Pepper Flakes" --unit kg --cost 28
./chefops ingredient add --name "Salt" --unit kg --cost 2

./chefops ingredient add --name "Garlic" --unit kg --cost 6
./chefops ingredient add --name "Garlic Powder" --unit kg --cost 22
./chefops ingredient add --name "Onion Powder" --unit kg --cost 22

./chefops ingredient add --name "Onion" --unit kg --cost 3
./chefops ingredient add --name "Tomato" --unit kg --cost 5
./chefops ingredient add --name "Tomato Paste" --unit kg --cost 12

./chefops ingredient add --name "Sugar" --unit kg --cost 3
./chefops ingredient add --name "Vinegar" --unit liter --cost 4

./chefops ingredient add --name "Fries" --unit kg --cost 7
./chefops ingredient add --name "Beer" --unit liter --cost 9

./chefops ingredient add --name "Ketchup" --unit kg --cost 6
./chefops ingredient add --name "Pickles" --unit kg --cost 10

./chefops ingredient add --name "Hot Sauce" --unit kg --cost 18
./chefops ingredient add --name "Honey" --unit kg --cost 26
./chefops ingredient add --name "Oregano" --unit kg --cost 30

./chefops ingredient add --name "Mustard Seeds" --unit kg --cost 16
./chefops ingredient add --name "Water" --unit liter --cost 0.1
./chefops ingredient add --name "Apple Cider Vinegar" --unit liter --cost 10
./chefops ingredient add --name "Shallots" --unit kg --cost 14
./chefops ingredient add --name "Bay Leaves" --unit kg --cost 38

./chefops ingredient add --name "Potatoes" --unit kg --cost 3.5
./chefops ingredient add --name "Malt Vinegar" --unit liter --cost 6
./chefops ingredient add --name "Gherkin Brine" --unit liter --cost 5
./chefops ingredient add --name "Vegetable Stock" --unit liter --cost 7
./chefops ingredient add --name "Parsley" --unit kg --cost 38

./chefops ingredient add --name "Potato Chips" --unit kg --cost 16

./chefops ingredient add --name "Lettuce" --unit kg --cost 7
./chefops ingredient add --name "Brioche Dog Bun" --unit piece --cost 2.8
./chefops ingredient add --name "Stracciatella" --unit kg --cost 52
./chefops ingredient add --name "Pepperoni" --unit kg --cost 28
./chefops ingredient add --name "Chicken Goujons" --unit kg --cost 22

echo "🧀 Importing bulk recipes (normalized to 1 kg)..."

# ---------------------------------------------------------------------
# SUB Mac And Cheese Base (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "SUB Mac And Cheese Base" --yield 1 --unit kg
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Pasta" --qty 0.4
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Milk" --qty 0.3
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Cheddar Cheese" --qty 0.2
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Butter" --qty 0.08
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Flour" --qty 0.04
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Dijon Mustard" --qty 0.01
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Salt" --qty 0.002
./chefops recipe add-item --recipe "SUB Mac And Cheese Base" --ingredient "Pepper" --qty 0.002

# ---------------------------------------------------------------------
# BULK Lobster Mac And Cheese (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Lobster Mac And Cheese" --yield 1 --unit kg
./chefops recipe add-subrecipe --recipe "BULK Lobster Mac And Cheese" --sub "SUB Mac And Cheese Base" --qty 0.666
./chefops recipe add-item --recipe "BULK Lobster Mac And Cheese" --ingredient "Lobster Meat" --qty 0.166
./chefops recipe add-item --recipe "BULK Lobster Mac And Cheese" --ingredient "Chives" --qty 0.016
./chefops recipe add-item --recipe "BULK Lobster Mac And Cheese" --ingredient "Lemon Zest" --qty 0.016
./chefops recipe add-item --recipe "BULK Lobster Mac And Cheese" --ingredient "Pepper" --qty 0.002
./chefops recipe add-item --recipe "BULK Lobster Mac And Cheese" --ingredient "Salt" --qty 0.002

# ---------------------------------------------------------------------
# BULK Lobster Salad (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Lobster Salad" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Lobster Salad" --ingredient "Lobster Meat" --qty 0.5
./chefops recipe add-item --recipe "BULK Lobster Salad" --ingredient "Dijon Mustard" --qty 0.033
./chefops recipe add-item --recipe "BULK Lobster Salad" --ingredient "Lemon Juice" --qty 0.033
./chefops recipe add-item --recipe "BULK Lobster Salad" --ingredient "Chives" --qty 0.033
./chefops recipe add-item --recipe "BULK Lobster Salad" --ingredient "Pepper" --qty 0.003
./chefops recipe add-item --recipe "BULK Lobster Salad" --ingredient "Salt" --qty 0.003

# ---------------------------------------------------------------------
# BULK Citrus Aioli (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Citrus Aioli" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Citrus Aioli" --ingredient "Mayonnaise" --qty 0.85
./chefops recipe add-item --recipe "BULK Citrus Aioli" --ingredient "Lemon Juice" --qty 0.09
./chefops recipe add-item --recipe "BULK Citrus Aioli" --ingredient "Lemon Zest" --qty 0.015
./chefops recipe add-item --recipe "BULK Citrus Aioli" --ingredient "Garlic" --qty 0.022
./chefops recipe add-item --recipe "BULK Citrus Aioli" --ingredient "Dijon Mustard" --qty 0.015
./chefops recipe add-item --recipe "BULK Citrus Aioli" --ingredient "Pepper" --qty 0.003
./chefops recipe add-item --recipe "BULK Citrus Aioli" --ingredient "Salt" --qty 0.005

# ---------------------------------------------------------------------
# BULK Hot Honey Drizzle (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Hot Honey Drizzle" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Hot Honey Drizzle" --ingredient "Honey" --qty 0.9
./chefops recipe add-item --recipe "BULK Hot Honey Drizzle" --ingredient "Red Pepper Flakes" --qty 0.05
./chefops recipe add-item --recipe "BULK Hot Honey Drizzle" --ingredient "Lemon Juice" --qty 0.04
./chefops recipe add-item --recipe "BULK Hot Honey Drizzle" --ingredient "Butter" --qty 0.01

# ---------------------------------------------------------------------
# BULK Lemon Pepper Seasoning (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Lemon Pepper Seasoning" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Lemon Pepper Seasoning" --ingredient "Lemon Zest" --qty 0.18
./chefops recipe add-item --recipe "BULK Lemon Pepper Seasoning" --ingredient "Salt" --qty 0.5
./chefops recipe add-item --recipe "BULK Lemon Pepper Seasoning" --ingredient "Pepper" --qty 0.3
./chefops recipe add-item --recipe "BULK Lemon Pepper Seasoning" --ingredient "Garlic Powder" --qty 0.02

# ---------------------------------------------------------------------
# BULK Basil Oil (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Basil Oil" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Basil Oil" --ingredient "Basil" --qty 0.12
./chefops recipe add-item --recipe "BULK Basil Oil" --ingredient "Grapeseed Oil" --qty 0.78
./chefops recipe add-item --recipe "BULK Basil Oil" --ingredient "Ice" --qty 0.1

# ---------------------------------------------------------------------
# BULK Batter (Beer Batter, 1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Batter" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Batter" --ingredient "Flour" --qty 0.45
./chefops recipe add-item --recipe "BULK Batter" --ingredient "Cornstarch" --qty 0.09
./chefops recipe add-item --recipe "BULK Batter" --ingredient "Beer" --qty 0.4
./chefops recipe add-item --recipe "BULK Batter" --ingredient "Salt" --qty 0.02
./chefops recipe add-item --recipe "BULK Batter" --ingredient "Pepper" --qty 0.002
./chefops recipe add-item --recipe "BULK Batter" --ingredient "Garlic Powder" --qty 0.002
./chefops recipe add-item --recipe "BULK Batter" --ingredient "Onion Powder" --qty 0.002

# ---------------------------------------------------------------------
# BULK Burger Sauce (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Burger Sauce" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Mayonnaise" --qty 0.6
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Ketchup" --qty 0.2
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Dijon Mustard" --qty 0.05
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Pickles" --qty 0.05
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Onion" --qty 0.05
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Garlic" --qty 0.01
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Pepper" --qty 0.003
./chefops recipe add-item --recipe "BULK Burger Sauce" --ingredient "Salt" --qty 0.003

# ---------------------------------------------------------------------
# BULK Pizza Sauce (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Pizza Sauce" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Tomato" --qty 0.75
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Tomato Paste" --qty 0.1
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Olive Oil" --qty 0.05
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Garlic" --qty 0.04
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Basil" --qty 0.02
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Oregano" --qty 0.02
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Salt" --qty 0.01
./chefops recipe add-item --recipe "BULK Pizza Sauce" --ingredient "Pepper" --qty 0.01

# ---------------------------------------------------------------------
# BULK Tomato Onion Jam (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Tomato Onion Jam" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Tomato Onion Jam" --ingredient "Onion" --qty 0.5
./chefops recipe add-item --recipe "BULK Tomato Onion Jam" --ingredient "Tomato" --qty 0.3
./chefops recipe add-item --recipe "BULK Tomato Onion Jam" --ingredient "Sugar" --qty 0.1
./chefops recipe add-item --recipe "BULK Tomato Onion Jam" --ingredient "Vinegar" --qty 0.08
./chefops recipe add-item --recipe "BULK Tomato Onion Jam" --ingredient "Salt" --qty 0.01
./chefops recipe add-item --recipe "BULK Tomato Onion Jam" --ingredient "Pepper" --qty 0.01

# ---------------------------------------------------------------------
# BULK Pickled Mustard Seeds (1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Pickled Mustard Seeds" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Mustard Seeds" --qty 0.3
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Water" --qty 1.5
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Apple Cider Vinegar" --qty 0.5
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Sugar" --qty 0.15
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Salt" --qty 0.02
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Shallots" --qty 0.05
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Pepper" --qty 0.005
./chefops recipe add-item --recipe "BULK Pickled Mustard Seeds" --ingredient "Bay Leaves" --qty 0.002

# ---------------------------------------------------------------------
# BULK Potato Salad (German-style, 1 kg)
# ---------------------------------------------------------------------
./chefops recipe new --name "BULK Potato Salad" --yield 1 --unit kg
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Potatoes" --qty 0.6
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Apple Cider Vinegar" --qty 0.03
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Mayonnaise" --qty 0.16
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Dijon Mustard" --qty 0.02
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Pickles" --qty 0.08
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Gherkin Brine" --qty 0.02
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Onion" --qty 0.06
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Vegetable Stock" --qty 0.04
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Extra Virgin Olive Oil" --qty 0.02
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Salt" --qty 0.006
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Pepper" --qty 0.001
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Chives" --qty 0.004
./chefops recipe add-item --recipe "BULK Potato Salad" --ingredient "Parsley" --qty 0.004

echo "🍽 Importing dish recipes (1 portion each)..."

# ---------------------------------------------------------------------
# DISH Pole Position Burger
# ---------------------------------------------------------------------
./chefops recipe new --name "DISH Pole Position Burger" --yield 1 --unit portion
./chefops recipe add-item --recipe "DISH Pole Position Burger" --ingredient "Beef Patty" --qty 0.18
./chefops recipe add-item --recipe "DISH Pole Position Burger" --ingredient "Cheddar Cheese" --qty 0.02
./chefops recipe add-item --recipe "DISH Pole Position Burger" --ingredient "Tomato" --qty 0.03
./chefops recipe add-item --recipe "DISH Pole Position Burger" --ingredient "Onion" --qty 0.02
./chefops recipe add-item --recipe "DISH Pole Position Burger" --ingredient "Pickles" --qty 0.02
./chefops recipe add-item --recipe "DISH Pole Position Burger" --ingredient "Burger Bun" --qty 1
./chefops recipe add-subrecipe --recipe "DISH Pole Position Burger" --sub "SUB Mac And Cheese Base" --qty 0.06
./chefops recipe add-subrecipe --recipe "DISH Pole Position Burger" --sub "BULK Burger Sauce" --qty 0.03
./chefops recipe add-subrecipe --recipe "DISH Pole Position Burger" --sub "BULK Tomato Onion Jam" --qty 0.02
./chefops recipe add-subrecipe --recipe "DISH Pole Position Burger" --sub "BULK Potato Salad" --qty 0.15

# ---------------------------------------------------------------------
# DISH Lobster Roll
# ---------------------------------------------------------------------
./chefops recipe new --name "DISH Lobster Roll" --yield 1 --unit portion
./chefops recipe add-item --recipe "DISH Lobster Roll" --ingredient "Brioche Bun" --qty 1
./chefops recipe add-subrecipe --recipe "DISH Lobster Roll" --sub "BULK Lobster Salad" --qty 0.12
./chefops recipe add-subrecipe --recipe "DISH Lobster Roll" --sub "BULK Citrus Aioli" --qty 0.03
./chefops recipe add-subrecipe --recipe "DISH Lobster Roll" --sub "BULK Batter" --qty 0.02

# ---------------------------------------------------------------------
# DISH Full Throttle Lobster, Mac And Cheese Croquette (5 pcs)
# ---------------------------------------------------------------------
./chefops recipe new --name "DISH Full Throttle Lobster, Mac And Cheese Croquette" --yield 1 --unit portion
./chefops recipe add-subrecipe --recipe "DISH Full Throttle Lobster, Mac And Cheese Croquette" --sub "BULK Lobster Mac And Cheese" --qty 0.30
./chefops recipe add-item --recipe "DISH Full Throttle Lobster, Mac And Cheese Croquette" --ingredient "Potato Chips" --qty 0.12
./chefops recipe add-item --recipe "DISH Full Throttle Lobster, Mac And Cheese Croquette" --ingredient "Lemon" --qty 0.05
./chefops recipe add-subrecipe --recipe "DISH Full Throttle Lobster, Mac And Cheese Croquette" --sub "BULK Pickled Mustard Seeds" --qty 0.01

# ---------------------------------------------------------------------
# DISH Margherita Pizza
# ---------------------------------------------------------------------
./chefops recipe new --name "DISH Margherita Pizza" --yield 1 --unit portion
./chefops recipe add-item --recipe "DISH Margherita Pizza" --ingredient "Pizza Base" --qty 1
./chefops recipe add-item --recipe "DISH Margherita Pizza" --ingredient "Mozzarella Cheese" --qty 0.12
./chefops recipe add-item --recipe "DISH Margherita Pizza" --ingredient "Basil" --qty 0.01
./chefops recipe add-item --recipe "DISH Margherita Pizza" --ingredient "Olive Oil" --qty 0.005
./chefops recipe add-subrecipe --recipe "DISH Margherita Pizza" --sub "BULK Pizza Sauce" --qty 0.08

# ---------------------------------------------------------------------
# DISH Hot Lap Honey & Pepperoni Pizza
# ---------------------------------------------------------------------
./chefops recipe new --name "DISH Hot Lap Honey & Pepperoni Pizza" --yield 1 --unit portion
./chefops recipe add-subrecipe --recipe "DISH Hot Lap Honey & Pepperoni Pizza" --sub "BULK Hot Honey Drizzle" --qty 0.020
./chefops recipe add-subrecipe --recipe "DISH Hot Lap Honey & Pepperoni Pizza" --sub "BULK Pizza Sauce" --qty 0.100
./chefops recipe add-item --recipe "DISH Hot Lap Honey & Pepperoni Pizza" --ingredient "Pepperoni" --qty 0.080
./chefops recipe add-item --recipe "DISH Hot Lap Honey & Pepperoni Pizza" --ingredient "Mozzarella Cheese" --qty 0.120

# ---------------------------------------------------------------------
# DISH Turbo Hammour Popcorn
# ---------------------------------------------------------------------
./chefops recipe new --name "DISH Turbo Hammour Popcorn" --yield 1 --unit portion
./chefops recipe add-subrecipe --recipe "DISH Turbo Hammour Popcorn" --sub "BULK Lemon Pepper Seasoning" --qty 0.010
./chefops recipe add-subrecipe --recipe "DISH Turbo Hammour Popcorn" --sub "BULK Citrus Aioli" --qty 0.030
./chefops recipe add-subrecipe --recipe "DISH Turbo Hammour Popcorn" --sub "BULK Batter" --qty 0.100
./chefops recipe add-subrecipe --recipe "DISH Turbo Hammour Popcorn" --sub "BULK Pickled Mustard Seeds" --qty 0.010
./chefops recipe add-item --recipe "DISH Turbo Hammour Popcorn" --ingredient "Hammour" --qty 0.150
./chefops recipe add-item --recipe "DISH Turbo Hammour Popcorn" --ingredient "Lemon" --qty 0.050

# ---------------------------------------------------------------------
# DISH Chequered Flag Chicken Goujons
# ---------------------------------------------------------------------
./chefops recipe new --name "DISH Chequered Flag Chicken Goujons" --yield 1 --unit portion
./chefops recipe add-subrecipe --recipe "DISH Chequered Flag Chicken Goujons" --sub "BULK Lemon Pepper Seasoning" --qty 0.010
./chefops recipe add-subrecipe --recipe "DISH Chequered Flag Chicken Goujons" --sub "BULK Citrus Aioli" --qty 0.030
./chefops recipe add-item --recipe "DISH Chequered Flag Chicken Goujons" --ingredient "Chicken Goujons" --qty 0.150
./chefops recipe add-item --recipe "DISH Chequered Flag Chicken Goujons" --ingredient "Lemon" --qty 0.050
./chefops recipe add-item --recipe "DISH Chequered Flag Chicken Goujons" --ingredient "Potato Chips" --qty 0.100

echo "✅ Import finished. You can now run ChefOps on the F1 dataset."
# END OF master_import_F1.sh
