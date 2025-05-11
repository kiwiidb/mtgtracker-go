host=https://mtgtracker.kwintendebacker.com
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
