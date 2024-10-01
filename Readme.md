# webs
> https://github.com/gen2brain/raylib-go
> https://cupnooble.itch.io/sprout-lands-asset-pack

# todo ->
warstwa np dla dachy która jest wyśwetlana tylko gdy gracza pod nią nie ma

# map data
line1: width, height
line2-x:
<nr_warstwy> <oznaczenie Sprita> <id boxa>
np: 1g05, 2g14

im większy nr warstwy tym póżniej = tym wyżej jest rysowany element.
> można potem napisać kod który sprawdza czy dany sprite na wyżej warstiw nie ma nic przeżroczystego, jeżeli tak to nie rysować niższych warstw

# map creator:
1. lista plików textór
2. jak wybierzemy plik to lista bloków
3. na ekranie wyświtlany grid
4. kliknięcie z wybranym blokiem kładzie go
5. na podstawie id bloktów jest generowana mapa

# multiplayer
1. gracz przesyła do serwera: pozycję x,y,playerDir,animationFrame
2. serwer przypisuje dane każdego gracza do jego ID
3. serwer emituje dane do graczy (playerID: dane)
4. gracze rysują innych graczy u siebie
    > jeżeli gracz miał laga to zostaje cofnięty do danego miejsca

+ gracz
    - po połączeniu posiada swoje ID

+ serwer