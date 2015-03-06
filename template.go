// 2015-03-06 Adam Bryt

package main

// template dla formatki biorytmu tekstowego
const formTextTmplStr = `
<html>
<head>
	<title>Biorytm</title>
</head>
<body>
	<p>
	Podaj datę urodzenia i datę biorytmu w formacie yyyy-mm-dd
	</p>

	<form action="/text/biorytm/">
		<table>
			<tr>
				<td style="text-align:right">data urodzenia:</td>
				<td><input type="text" name="born"></td>
			</tr>
			<tr>
				<td style="text-align:right">data aktualna:</td>
				<td><input type="text" name="date" value="{{ .DateString }}"></td>
			</tr>
			<tr>
				<td style="text-align:right">ilość dni:</td>
				<td><input type="text" name="range" value="{{ .Days }}"></td>
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

// template dla wyniku biorytmu tekstowego
const outputTextTmplStr = `
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
