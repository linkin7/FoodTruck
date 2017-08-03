package htmltemplate

var Index = `
<a href="/register">Register</a>
<a href="/login">Log In</a>

<br/>
<b>Enter your current address to find nearest Food Truck service</b>
<form action="/findnearest" method="POST">
	Address: <input type="text" name="address"><br/>
  	<input type="submit" name="find"><br/><br/>
</form>
`
