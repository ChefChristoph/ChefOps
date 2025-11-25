#!/bin/sh

CHEFOPS=./chefops
OUTDIR=forecasts
MASTER="$OUTDIR/f1_master_forecast.csv"

mkdir -p "$OUTDIR"

# Clean previous master
rm -f "$MASTER"

run_forecast() {
  dish="$1"
  portions="$2"
  slug="$3"

  tmp="$OUTDIR/${slug}.csv"

  echo "▶ Forecasting: $dish  ($portions portions)"

  # FLAGS FIRST, NAME LAST  ✔✔✔ FIXED
  "$CHEFOPS" forecast --out "$tmp" "$dish=$portions"

  if [ ! -f "$tmp" ]; then
    echo "  ⚠️ No CSV created for $dish, skipping."
    return
  fi

  if [ ! -f "$MASTER" ]; then
    cat "$tmp" >"$MASTER"
  else
    tail -n +2 "$tmp" >>"$MASTER"
  fi
}

# ----------------------------------------------------
# Edit portions here for F1 forecast
# ----------------------------------------------------

run_forecast "DISH Pole Position Burger" 300 pole_position_burger_300
run_forecast "DISH Lobster Roll" 250 lobster_roll_250
run_forecast "DISH Full Throttle Lobster, Mac And Cheese Croquette" 400 lobster_croquette_400
run_forecast "DISH Margherita Pizza" 350 margherita_350
run_forecast "DISH Hot Lap Honey & Pepperoni Pizza" 500 hot_lap_pepperoni_500
run_forecast "DISH Turbo Hammour Popcorn" 150 hammour_popcorn_150
run_forecast "DISH Chequered Flag Chicken Goujons" 200 chicken_goujons_200

echo
echo "✅ Combined forecast written to: $MASTER"
