#!/bin/sh

echo "🚨 Resetting ChefOps database..."
rm -f db/chefops.db

echo "🔧 Recreating schema..."
sqlite3 db/chefops.db <schema.sql
sqlite3 db/chefops.db <views.sql

echo "🍏 Importing ingredients..."
chefops ingredient add --name "Lobster Meat" --unit kg --cost 0
chefops ingredient add --name "Hammour" --unit kg --cost 0
chefops ingredient add --name "Pasta" --unit kg --cost 0
chefops ingredient add --name "Breadcrumbs" --unit kg --cost 0
chefops ingredient add --name "Flour" --unit kg --cost 0
chefops ingredient add --name "Cornstarch" --unit kg --cost 0
chefops ingredient add --name "Pizza Base" --unit piece --cost 0
chefops ingredient add --name "Brioche Bun" --unit piece --cost 0
chefops ingredient add --name "Burger Bun" --unit piece --cost 0
chefops ingredient add --name "Cheddar Cheese" --unit kg --cost 0
chefops ingredient add --name "Mozzarella Cheese" --unit kg --cost 0
chefops ingredient add --name "Grana Padano Cheese" --unit kg --cost 0
chefops ingredient add --name "Butter" --unit kg --cost 0
chefops ingredient add --name "Milk" --unit kg --cost 0
chefops ingredient add --name "Egg Yolks" --unit piece --cost 0
chefops ingredient add --name "Whole Eggs" --unit piece --cost 0
chefops ingredient add --name "Dijon Mustard" --unit kg --cost 0
chefops ingredient add --name "Chives" --unit kg --cost 0
chefops ingredient add --name "Lemon" --unit kg --cost 0
chefops ingredient add --name "Lemon Juice" --unit liter --cost 0
chefops ingredient add --name "Lemon Zest" --unit kg --cost 0
chefops ingredient add --name "Olive Oil" --unit liter --cost 0
chefops ingredient add --name "Extra Virgin Olive Oil" --unit liter --cost 0
chefops ingredient add --name "Grapeseed Oil" --unit liter --cost 0
chefops ingredient add --name "Sunflower Oil" --unit liter --cost 0
chefops ingredient add --name "Ice" --unit kg --cost 0
chefops ingredient add --name "Pepper" --unit kg --cost 0
chefops ingredient add --name "Red Pepper Flakes" --unit kg --cost 0
chefops ingredient add --name "Salt" --unit kg --cost 0
chefops ingredient add --name "Garlic" --unit kg --cost 0
chefops ingredient add --name "Garlic Powder" --unit kg --cost 0
chefops ingredient add --name "Onion Powder" --unit kg --cost 0
chefops ingredient add --name "Onion" --unit kg --cost 0
chefops ingredient add --name "Tomato" --unit kg --cost 0
chefops ingredient add --name "Sugar" --unit kg --cost 0
chefops ingredient add --name "Vinegar" --unit liter --cost 0
chefops ingredient add --name "Dill" --unit kg --cost 0
chefops ingredient add --name "Basil" --unit kg --cost 0
chefops ingredient add --name "Beer" --unit liter --cost 0
chefops ingredient add --name "Ketchup" --unit kg --cost 0
chefops ingredient add --name "Pickles" --unit kg --cost 0
chefops ingredient add --name "Tomato Paste" --unit kg --cost 0
chefops ingredient add --name "Fries" --unit kg --cost 0
chefops ingredient add --name "Hot Sauce" --unit kg --cost 0

echo "🧀 Importing bulk recipes (1 kg normalized)..."

# Mac And Cheese Base
chefops recipe new --name "Mac And Cheese Base" --yield 1 --unit kg
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Pasta" --qty 0.4
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Milk" --qty 0.3
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Cheddar Cheese" --qty 0.2
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Butter" --qty 0.08
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Flour" --qty 0.04
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Dijon Mustard" --qty 0.01
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Salt" --qty 0.002
chefops recipe add-item --recipe "Mac And Cheese Base" --ingredient "Pepper" --qty 0.002

# Lobster Mac And Cheese
chefops recipe new --name "Lobster Mac And Cheese" --yield 1 --unit kg
chefops recipe add-subrecipe --recipe "Lobster Mac And Cheese" --sub "Mac And Cheese Base" --qty 0.666
chefops recipe add-item --recipe "Lobster Mac And Cheese" --ingredient "Lobster Meat" --qty 0.166
chefops recipe add-item --recipe "Lobster Mac And Cheese" --ingredient "Chives" --qty 0.016
chefops recipe add-item --recipe "Lobster Mac And Cheese" --ingredient "Lemon Zest" --qty 0.016
chefops recipe add-item --recipe "Lobster Mac And Cheese" --ingredient "Pepper" --qty 0.002
chefops recipe add-item --recipe "Lobster Mac And Cheese" --ingredient "Salt" --qty 0.002

# ... (All other bulk recipes continue similarly)
# I will include the rest in your final ready-to-save file:
# ﹣ Citrus Aioli
# ﹣ Lobster Salad
# ﹣ Basil Oil
# ﹣ Batter
# ﹣ Lemon Pepper Seasoning
# ﹣ Burger Sauce
# ﹣ Hot Honey Drizzle
# ﹣ Pizza Sauce
# ﹣ Tomato Onion Jam

echo "🍽 Importing dish recipes..."

# Pole Position Burger
chefops recipe new --name "Pole Position Burger" --yield 1 --unit portion
chefops recipe add-item --recipe "Pole Position Burger" --ingredient "Beef Patty" --qty 0.18
chefops recipe add-item --recipe "Pole Position Burger" --ingredient "Cheddar Cheese" --qty 0.02
chefops recipe add-item --recipe "Pole Position Burger" --ingredient "Tomato" --qty 0.03
chefops recipe add-item --recipe "Pole Position Burger" --ingredient "Onion" --qty 0.02
chefops recipe add-item --recipe "Pole Position Burger" --ingredient "Pickles" --qty 0.02
chefops recipe add-item --recipe "Pole Position Burger" --ingredient "Burger Bun" --qty 1
chefops recipe add-subrecipe --recipe "Pole Position Burger" --sub "Mac And Cheese Base" --qty 0.06
chefops recipe add-subrecipe --recipe "Pole Position Burger" --sub "Burger Sauce" --qty 0.03
chefops recipe add-subrecipe --recipe "Pole Position Burger" --sub "Tomato Onion Jam" --qty 0.02

echo "✔ Import finished."
