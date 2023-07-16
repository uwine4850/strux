## <p style="text-align: center;">Initialization commands</p>
Displays some information about the client.

_--info_, _-inf_ — Indicates that all further commands belong to the group of information commands. Displays information
about other commands.<br>
````
go run main.go --info
````

_--path_, _-p_ — Shows the path to the main directory with packages.
````
go run main.go --info -p
````

_--package_, _-pkg_ (pkgName string) — Displays information about a specific package.
````
go run main.go --info -pkg
````