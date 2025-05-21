host=localhost:8080
id=Kwinten;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Alex;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Lorin;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Lucas;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Dries;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Arthur;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Braïn;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Jari;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Laura;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Cyril;http POST $host/player/v1/signup name=$id email=$id image=$id
id=Heather;http POST $host/player/v1/signup name=$id email=$id image=$id

http POST $host/player/v1/groups name="Cardboard Crack" creator_id:=1

http PUT $host/player/v1/groups/1/add/Braïn
http PUT $host/player/v1/groups/1/add/Kwinten
http PUT $host/player/v1/groups/1/add/Alex
http PUT $host/player/v1/groups/1/add/Lorin
http PUT $host/player/v1/groups/1/add/Lucas
http PUT $host/player/v1/groups/1/add/Dries
http PUT $host/player/v1/groups/1/add/Laura
http PUT $host/player/v1/groups/1/add/Arthur
http PUT $host/player/v1/groups/1/add/Jari
http PUT $host/player/v1/groups/1/add/Cyril

# Create a game
# Kwinten plays with Alania, Divergent Storm
# Lorin plays with Massacre Girl, Known Killer
# Alex plays with Teysa, Envoy of Ghosts
# Lucas plays with Jodah, The Unifier

# type CreateGameRequest struct {
# 	GroupID  uint       `json:"group_id"`
# 	Rankings []Ranking  `json:"rankings"`
# }
# type Ranking struct {
# 	PlayerID       uint   `json:"player_id"`
# 	Commander      string `json:"commander"`
# 	Position       int    `json:"position"`
# }
http POST $host/game/v1/games \
  group_id:=1 \
  rankings:='[{"player_id":1,"commander":"Alania, Divergent Storm","position":1},{"player_id":2,"commander":"Massacre Girl, Known Killer","position":2},{"player_id":3,"commander":"Teysa, Envoy of Ghosts","position":3},{"player_id":4,"commander":"Jodah, The Unifier","position":4}]'

  http POST $host/game/v1/games \
  group_id:=1 \
  rankings:='[{"player_id":1,"commander":"Edgar Markov","position":1},{"player_id":2,"commander":"Tymna the Weaver","position":2},{"player_id":3,"commander":"Isshin, Two Heavens as One","position":3},{"player_id":4,"commander":"Jodah, Archmage Eternal","position":4}]'