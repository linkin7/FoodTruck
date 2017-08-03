package htmltemplate

var Register = `
<h1>Registration for Food Truck</h1>

<form action="/postregister" method="POST">
	Name: <input type="text" name="name"><br/>
	Password: <input type="text" name="password"><br/>
	<select name="cuisine">
    	<option value="" disabled="disabled" selected="selected">Please select a cuisine</option>
  		<option value="italian">Italian</option>
  		<option value="french">French</option>
  		<option value="mexican">Mexican</option>
	</select>
  	<input type="submit" value="Submit">
</form>
`
