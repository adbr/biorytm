// 2015-03-06 Adam Bryt

package main

import "html/template"

// formHTML definiuje template formatki do wprowadzania danych dla
// biorytmu.
const formHTML = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<p>
	Podaj datę urodzenia i datę biorytmu w formacie yyyy-mm-dd
	</p>

	<form action="/display/">
		<table>
			<tr>
				<td style="text-align:right">data urodzenia:</td>
				<td><input type="text" name="born" value="{{ .Born }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">data aktualna:</td>
				<td><input type="text" name="date" value="{{ .Date }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">ilość dni:</td>
				<td><input type="text" name="days" value="{{ .Days }}"></td>
			</tr>
			<tr>
				<td style="text-align:right;vertical-align:top">prezentacja:</td>
				<td>
					<input type="radio" name="look" value="text" checked>tekst<br>
					<input type="radio" name="look" value="graph">grafika
				</td>
			</tr>
			<tr>
				<td></td>
				<td style="text-align:right"><input type="submit" value="OK"></td>
			</tr>
		</table>
	</form>
</body>
</html>
`

// Template formatki do wprowadzania danych dla biorytmu.
var formTmpl = template.New("form")

// formData zawiera dane dla template formHTML.
type formData struct {
	Born string // data urodzenia yyyy-mm-dd
	Date string // data biorytmu yyyy-mm-dd
	Days int    // zakres dni biorytmu
}

// textHTML definiuje template strony wyświetlającej biorytm w postaci
// tekstowej.
const textHTML = `
<html>
<head>
	<title>Biorytm</title>
	<!--
	<style>
		body {background-color:black; color:#FFC200}
	</style>
	-->
</head>
<body>
	<p>
	<pre>{{ .Text }}</pre>
	</p>
</body>
</html>
`

// Template strony wyświetlającej biorytm w postaci tekstowej.
var textTmpl = template.New("text")

// textData zawiera dane dla template textHTML
type textData struct {
	Text string // biorytm w postaci tekstuowej
}

// graphHTML definiuje template strony wyświetlającej biorytm w
// postaci graficznej.
const graphHTML = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<h2>
	Biorytm
	</h2>

	<p style="color:red">Not Yet Implemented</p>

	<table>
		<tr>
			<td>Data urodzenia:</td>
			<td>{{ .Born }}</td>
		</tr>
		<tr>
			<td>Data biorytmu:</td>
			<td>{{ .Date }}</td>
		</tr>
	</table>

	<p>
	<img src="data:image/png;base64,{{ .Image }}" alt="Wykres biorytmu">
	</p>
</body>
</html>
`

// Template strony wyświetlającej biorytm w postaci graficznej.
var graphTmpl = template.New("graph")

// graphData zawiera dane dla template graphHTML.
type graphData struct {
	Born  string // data urodzenie yyyy-mm-dd
	Date  string // data biorytmu yyyy-mm-dd
	Image string // obrazek PNG zakodowany w base64
}
