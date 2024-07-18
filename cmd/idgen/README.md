### ID-GEN
Short explanation: this will generate a list of appIDs from a Steam account that:
- Have at least 1 achievement
    - That isn't unlocked

To use the id-gen, you will need to generate a [Steam WebAPI Key](https://steamcommunity.com/dev/apikey) from https://steamcommunity.com/dev/apikey.
(If you don't have one already. For the domain name, just put `localhost`), then run `go run main.go --key=STEAM_WEBAPI_KEY (--steamid=YOUR_STEAM_ID)` (from this directory), then sit back and wait. (--steamid is optional if you don't want to initialize Steam just to get your id, in that case, --libraryPath must also be set)
The process can take anywhere from 1 minute to 20 minutes, depending on how many games you have, once it's finished, you can use the generated `id.txt` with the wrapper in the root directory. (`go run main.go --idPath="./cmd/idgen/id.txt"`)