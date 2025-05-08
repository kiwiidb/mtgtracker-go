id=kwinten;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=alex;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=lorin;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=lucas;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=dries;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=arthur;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=brain;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=jari;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=laura;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id

http POST localhost:8080/player/v1/groups name="Cardboard Crack" creator_id:=1

http PUT localhost:8080/player/v1/groups/1/add/brain
http PUT localhost:8080/player/v1/groups/1/add/kwinten
http PUT localhost:8080/player/v1/groups/1/add/alex
http PUT localhost:8080/player/v1/groups/1/add/lorin
http PUT localhost:8080/player/v1/groups/1/add/lucas
http PUT localhost:8080/player/v1/groups/1/add/dries
http PUT localhost:8080/player/v1/groups/1/add/laura
http PUT localhost:8080/player/v1/groups/1/add/arthur
http PUT localhost:8080/player/v1/groups/1/add/jari
