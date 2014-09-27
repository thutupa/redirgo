# Compilation Issues
* Create a file called appid.go with the following contents in the same directory as app.go

```
package app

const clientID = ""
```
replace the clientID with the one obtained from this page:
```
https://console.developers.google.com/project/your-app-id/apiui/api
```
where your-app-id is replaced by the id of your app, ofcourse.

* set GOPATH to the directory with app.go so that go can see the endpoints source.

# Missing Features
* Edit
* Delete

