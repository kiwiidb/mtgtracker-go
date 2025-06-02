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

# Create a game
# Kwinten plays with Alania, Divergent Storm
# Lorin plays with Massacre Girl, Known Killer
# Alex plays with Teysa, Envoy of Ghosts
# Lucas plays with Jodah, The Unifier

# type CreateGameRequest struct {
# 	Rankings []Ranking  `json:"rankings"`
# }
# type Ranking struct {
# 	PlayerID       uint   `json:"player_id"`
# 	Commander      string `json:"commander"`
# 	Position       int    `json:"position"`
# }
  http POST $host/game/v1/games \
  rankings:='[{"player_id":1,"commander":"Edgar Markov","position":1},{"player_id":2,"commander":"Tymna the Weaver","position":2},{"player_id":3,"commander":"Isshin, Two Heavens as One","position":3},{"player_id":4,"commander":"Jodah, Archmage Eternal","position":4}]'
  http PUT localhost:8080/game/v1/games/1 \
  rankings:='[{"player_id":1,"commander":"Edgar Markov","position":2},{"player_id":2,"commander":"Tymna the Weaver","position":1},{"player_id":3,"commander":"Isshin, Two Heavens as One","position":3},{"player_id":4,"commander":"Jodah, Archmage Eternal","position":4}]';