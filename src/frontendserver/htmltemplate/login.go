package htmltemplate

var Login = `
<h1>Log In</h1>

<form action="/postlogin" method="POST">
	Name: <input type="text" name="name"><br/>
	Password: <input type="text" name="password"><br/>
  	<input type="submit" value="submit">
</form>
`
