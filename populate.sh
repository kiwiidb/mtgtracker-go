id=Kwinten;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Alex;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Lorin;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Lucas;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Dries;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Arthur;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Braïn;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Jari;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Laura;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id
id=Cyril;http POST localhost:8080/player/v1/signup name=$id email=$id image=$id

http POST localhost:8080/player/v1/groups name="Cardboard Crack" creator_id:=1

http PUT localhost:8080/player/v1/groups/1/add/Braïn
http PUT localhost:8080/player/v1/groups/1/add/Kwinten
http PUT localhost:8080/player/v1/groups/1/add/Alex
http PUT localhost:8080/player/v1/groups/1/add/Lorin
http PUT localhost:8080/player/v1/groups/1/add/Lucas
http PUT localhost:8080/player/v1/groups/1/add/Dries
http PUT localhost:8080/player/v1/groups/1/add/Laura
http PUT localhost:8080/player/v1/groups/1/add/Arthur
http PUT localhost:8080/player/v1/groups/1/add/Jari
http PUT localhost:8080/player/v1/groups/1/add/Cyril
