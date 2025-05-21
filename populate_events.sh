url="http://localhost:8080/game/v1/games/1/events"
r1=1
r2=2
r3=3
r4=4

# Initial life totals
life_r1=40
life_r2=40
life_r3=40
life_r4=40

# Players alive
alive_r1=1
alive_r2=1
alive_r3=1
alive_r4=1

# Helper function to pick a random delta between 1 and 20
rand_delta() {
  echo $(( ( RANDOM % 10 ) + 1 ))
}

# Helper function to pick a random alive player
pick_target() {
  local arr=("$@")
  local idx=$(( RANDOM % ${#arr[@]} ))
  echo "${arr[$idx]}"
}

# add an event for every player to set their life total at the start to 40
http POST $url life_total_after:=40 source_ranking_id:=$r1 target_ranking_id:=$r1
http POST $url life_total_after:=40 source_ranking_id:=$r2 target_ranking_id:=$r2
http POST $url life_total_after:=40 source_ranking_id:=$r3 target_ranking_id:=$r3
http POST $url life_total_after:=40 source_ranking_id:=$r4 target_ranking_id:=$r4
while true; do
  # Build array of alive rankings and their life totals
  alive_ids=()
  alive_life=()
  [ $alive_r1 -eq 1 ] && alive_ids+=($r1) && alive_life+=($life_r1)
  [ $alive_r2 -eq 1 ] && alive_ids+=($r2) && alive_life+=($life_r2)
  [ $alive_r3 -eq 1 ] && alive_ids+=($r3) && alive_life+=($life_r3)
  [ $alive_r4 -eq 1 ] && alive_ids+=($r4) && alive_life+=($life_r4)

  # Stop if only one player is alive
  if [ ${#alive_ids[@]} -le 1 ]; then
    break
  fi

  # Pick a random source and target (must be alive and not the same)
  while :; do
    src_idx=$(( RANDOM % ${#alive_ids[@]} ))
    tgt_idx=$(( RANDOM % ${#alive_ids[@]} ))
    [ $src_idx -ne $tgt_idx ] && break
  done
  src_id=${alive_ids[$src_idx]}
  tgt_id=${alive_ids[$tgt_idx]}
  tgt_life=${alive_life[$tgt_idx]}

  # Pick a random delta, but not more than target's life
  delta=$(rand_delta)
  if [ $delta -gt $tgt_life ]; then
    delta=$tgt_life
  fi

  # Calculate new life total
  new_life=$(( tgt_life - delta ))

  # Post event
  http POST $url damage_delta:=$delta life_total_after:=$new_life source_ranking_id:=$src_id target_ranking_id:=$tgt_id

  # Update local life totals and alive status
  if [ $tgt_id -eq $r1 ]; then
    life_r1=$new_life
    [ $life_r1 -le 0 ] && alive_r1=0
  elif [ $tgt_id -eq $r2 ]; then
    life_r2=$new_life
    [ $life_r2 -le 0 ] && alive_r2=0
  elif [ $tgt_id -eq $r3 ]; then
    life_r3=$new_life
    [ $life_r3 -le 0 ] && alive_r3=0
  elif [ $tgt_id -eq $r4 ]; then
	life_r4=$new_life
	[ $life_r4 -le 0 ] && alive_r4=0
  fi

  # Sleep 1-2 seconds
  sleep $(( ( RANDOM % 2 ) + 1 ))
done