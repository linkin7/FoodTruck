package htmltemplate

var Index = `
<a href="/register">Register</a>
<a href="/login">Log In</a>

<br/>
<b>Enter your current location</b>
<form action="/findnearest" method="POST">
	Latitude: <input type="text" name="latitude"><br/>
	Longitude: <input type="text" name="longitude"><br/>
  	<input type="submit" name="find"><br/><br/>
</form>
`
