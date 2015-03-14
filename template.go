// 2015-03-06 Adam Bryt

package main

// textFormHTML definiuje template formatki do wprowadzania danych dla
// biorytmu w wersji tekstowej.
const textFormHTML = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<p>
	Podaj datę urodzenia i datę biorytmu w formacie yyyy-mm-dd
	</p>

	<form action="/text/display/">
		<table>
			<tr>
				<td style="text-align:right">data urodzenia:</td>
				<td><input type="text" name="born" value="{{ .BornString }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">data aktualna:</td>
				<td><input type="text" name="date" value="{{ .DateString }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">ilość dni:</td>
				<td><input type="text" name="range" value="{{ .Drange }}"></td>
			</tr>
			<tr>
				<td></td>
				<td style="text-align:right"><input type="submit"></td>
			</tr>
		</table>
	</form>
</body>
</html>
`

// textDisplayHTML definiuje template strony wyświetlającej biorytm w
// postaci tekstowej.
const textDisplayHTML = `
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
	<pre>{{ . }}</pre>
	</p>
</body>
</html>
`

// graphFormHTML definiuje template formatki do wprowadzania danych dla
// biorytmu w wersji graficznej.
const graphFormHTML = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<p>
	Podaj datę urodzenia i datę biorytmu w formacie yyyy-mm-dd
	</p>

	<form action="/graph/display/">
		<table>
			<tr>
				<td style="text-align:right">data urodzenia:</td>
				<td><input type="text" name="born" value="{{ .BornString }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">data aktualna:</td>
				<td><input type="text" name="date" value="{{ .DateString }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">ilość dni:</td>
				<td><input type="text" name="range" value="{{ .Drange }}"></td>
			</tr>
			<tr>
				<td></td>
				<td style="text-align:right"><input type="submit"></td>
			</tr>
		</table>
	</form>
</body>
</html>
`

// graphDisplayHTML definiuje template strony wyświetlającej biorytm w
// postaci graficznej.
const graphDisplayHTML = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<h2>
	Biorytm
	</h2>

	<table>
		<tr>
			<td>Data urodzenia:</td>
			<td>{{ .BornString }}</td>
		</tr>
		<tr>
			<td>Data biorytmu:</td>
			<td>{{ .DateString }}</td>
		</tr>
	</table>

	<p>
	<img src="data:image/png;base64,{{ .ImageString }}" alt="Wykres biorytmu">
	</p>
</body>
</html>
`
