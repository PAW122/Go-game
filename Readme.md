# webs

> https://github.com/gen2brain/raylib-go > https://cupnooble.itch.io/sprout-lands-asset-pack

# todo ->

1. itemy
    + scrolowalny pasek itemów bez wchodzenia w eq
    + lista objektów/structow z danymi itemów (assetFile, name, size it)
        > tak żeby można było je użyć w eq,postać mogła je trzymać itp
        > można było pobrać ich dane z nazwy np GetItemData(<name>)

    + na tej liście mozna też zrobić te stronkę z itemList w Book.
    + udostępnianie danych
        > ma posiadać swoje zmienne lokalne zamiast globalnych
        > ma udostęniać je funkcjami np GetSelectedItem()

* warstwa np dla dachy która jest wyśwetlana tylko gdy gracza pod nią nie ma

* jak user będzie miał w ręce bron i kliknie lewy myszy to odpala się animacja
    + tworzony jest box (range ataku broni w danym kierunku) sprawdzający czy będzie mial kolizję z czymś co można atakować
        > jeżeli tak to zada obrażenia

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

- gracz

  - po połączeniu posiada swoje ID

- serwer


# gameplay:
może coś pod styl albiona?
dungerony, mapy, pve, pvp

drzewka umiejek
gildie
opcje zarabiania:
farmy,
zbieractwo,
kraftowanie,
zlecenia od graczy,
zabijanie mobow,

# sklep:
- waluta premium,
- skiny (własne / przedmiotow)
    > mozna kupić skina który będzie dostępny na danym itemie ale jak dropniesz item to zniknie skin
    (tylko osoba płacąca go używa (nie przypisuje się do itemu))
        > opcja przypisania takiego skina do konkretnego itemu, ale wtedy znikaa z eq gracza i przypisany jest do itema
    
    > 2 rodzaj skina który jest przypisany do przedmiotu (taki item może mieć wartość runkową na steam)
    
- txtPacki zmieniające wygląd mapy

# ComunityContent:
- opcja kupienia/wynajmu własnego serwera gdzie właściciel może wgrywać swoje pluginy (lua albo dll)
- opcja tworzenia pluginów do serwerów
- market:
    + zrobić opcję sprzedawania/handlu itemami
    + pozowlić graczą wystawiać itemy na sprzedaż