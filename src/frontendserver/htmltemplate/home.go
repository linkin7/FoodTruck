package htmltemplate

var Home = `
<a href="/logout">Log Out</a>
<form action="/update" method="POST">
	<b>Want to start the food service? Enter the location of the truck:</b><br/>
	Latitude: <input type="text" name="latitude"><br/>
	Longitude: <input type="text" name="longitude"><br/>
  	<input type="submit" name="start" value="start"><br/><br/>

  	<b>Closing the service?</b>
  	<input type="submit" name="close" value="yes"><br/>
</form>
`
