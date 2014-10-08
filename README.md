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
* Delete
Next steps
* Prase and store the domain of the link separately.
* Make domain a top level feature
* Remove filter field and merge functionality to add
  * Leverage this duality to instantly surface any duplicates, even before add
* Implement redirect.

## Development Log
* Spent of a lot of time figuring out endpoint routing (it seems to drop the google cookie when routing for endopoints, so user.CurrentUser() does not work.
* Spent a bunch of time figuring out integration of angularjs and gapi code. Finally figured out angular.bootstrap. Still don't undersatnd providers.
* Had to wipe out the datastore due to the above issue (was using u.ID in the ancestor key, cannot do that).
