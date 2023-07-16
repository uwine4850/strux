## <p style="text-align: center;">Initialization commands</p>
Commands for user management.

_--user_, _-usr_ — Indicates that all commands from now on are user-specific.<br>

_--register_, _-reg_ — Register a new user.<br>
````
go run main.go --user -reg
````

_--update-password_, _-upd-pass_ — Update the user's password.<br>
````
go run main.go --user -upd-pass
````

_--delete_, _-del_ — Deleting a user. This will delete all its packages from the server.<br>
````
go run main.go --user -del
````