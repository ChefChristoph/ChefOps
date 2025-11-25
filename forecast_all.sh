#!/bin/sh

# ----------------------------------------------
# Bulk Forecast Script for ALL DISH Recipes
# ----------------------------------------------
# This script produces forecast CSV files for each dish.
# You can adjust the portion numbers below.
# ----------------------------------------------

CHEFOPS=./chefops
OUTDIR=forecasts

mkdir -p "$OUTDIR"

# -----------------------
# Dish: Pole Position Burger
# -----------------------
$CHEFOPS forecast --out "$OUTDIR/pole_position_burger_300.csv" "DISH Pole Position Burger=300"

# -----------------------
# Dish: Lobster Roll
# -----------------------
$CHEFOPS forecast --out "$OUTDIR/lobster_roll_250.csv" "DISH Lobster Roll=250"

# -----------------------
# Dish: Full Throttle Lobster, Mac And Cheese Croquette
# -----------------------
$CHEFOPS forecast --out "$OUTDIR/lobster_croquette_400.csv" "DISH Full Throttle Lobster, Mac And Cheese Croquette=400"

# -----------------------
# Dish: Margherita Pizza
# -----------------------
$CHEFOPS forecast --out "$OUTDIR/margherita_350.csv" "DISH Margherita Pizza=350"

# -----------------------
# Dish: Hot Lap Honey & Pepperoni Pizza
# -----------------------
$CHEFOPS forecast --out "$OUTDIR/hot_lap_pepperoni_500.csv" "DISH Hot Lap Honey & Pepperoni Pizza=500"

# -----------------------
# Dish: Turbo Hammour Popcorn
# -----------------------
$CHEFOPS forecast --out "$OUTDIR/hammour_popcorn_150.csv" "DISH Turbo Hammour Popcorn=150"

# -----------------------
# Dish: Chequered Flag Chicken Goujons
# -----------------------
$CHEFOPS forecast --out "$OUTDIR/chicken_goujons_200.csv" "DISH Chequered Flag Chicken Goujons=200"

echo "✅ Forecasts generated in: $OUTDIR/"
