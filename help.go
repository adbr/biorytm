// 2015-04-21 Adam Bryt

package main

const usageStr = `Sposób użycia:
	biorytm [flagi] -born <data urodzenia>
	biorytm [flagi] -http <host:port>

Flagi:
	-born="": data urodzenia w formacie yyyy-mm-dd
	-date="": data biorytmu w formacie yyyy-mm-dd (domyślnie: dzisiaj)
	-http="": adres usługi HTTP (np. ':5050')
	-days=15: liczba dni biorytmu
	-fonts="": katalog z fontami (domyślnie ./fonts lub $VGFONTPATH)
	-help=false: wyświetla help`

// Wartość helpStr jest kopią dokumentacji z doc.go.
const helpStr = `Program wyświetla biorytm dla podanego zakresu dni.

Sposób użycia:

	biorytm [flagi] -born <data urodzenia>
	biorytm [flagi] -http <host:port>

Flagi:

	-born=""
		data urodzenia w formacie yyyy-mm-dd
	-date=""
		data biorytmu w formacie yyyy-mm-dd (domyślnie: dzisiaj)
	-http=""
		adres usługi HTTP (np. ':5050')
	-days=15
		liczba dni biorytmu
	-fonts=""
		katalog z fontami (domyślnie ./fonts lub $VGFONTPATH)
	-help=false
		wyświetla help

Program ma dwa tryby pracy.

Jeśli nie podano opcji -http, biorytm jest drukowany na stdout.  Opcja -days
określa ile dni biorytmu wyświetlić, przy czym data podana w opcji -date, lub
data aktualna, znajduje się w środku tego zakresu.

Jeśli podano opcję -http to program działa jako serwer HTTP.  Parametry
biorytmu można zmienić w formatce na stronie, a wyniki są prezentowane w
postaci tekstowej lub w postaci wykresu.

Informacje na temat biorytmów można znaleźć np. na stronie
http://en.wikipedia.org/wiki/Biorhythm

Przykłady

Uruchomienie w trybie CLI:

	biorytm -born 1970-01-01

Na standardowe wyjście zostanie wydrukowany biorytm w postaci:

	Data urodzenia: 1970-01-01
	Data biorytmu:  2015-04-21
	Liczba dni:     16546

	Data         Fizyczny           Psychiczny         Intelektualny
	2015-04-14   F: +0.52 ( 2/23)   P: -0.90 (19/28)   I: +0.91 ( 6/33)
	2015-04-15   F: +0.73 ( 3/23)   P: -0.97 (20/28)   I: +0.97 ( 7/33)
	2015-04-16   F: +0.89 ( 4/23)   P: -1.00 (21/28)m  I: +1.00 ( 8/33)M
	[...]

Uruchomienie w trybie serwera HTTP:

	biorytm -http :5050

Z tak uruchomionym serwerem można się połączyć przeglądarką WWW podając adres
'localhost:5050'. Następnie w formatce na stronie można wprowadzić parametry
biorytmu i wybrać sposób prezentacji.`
