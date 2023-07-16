## <p style="text-align: center;">Commands to the package manager</p>
Commands that manage already created packages.

_--package_, _-pkg_ (pkgName string) — This command alone does not perform specific operations, it only indicates which
package you want to access.<br>

_--add-dir_, _-ad_ (addPath string) — Adds a directory to the selected package at the specified path.
````
go run main.go -pkg <pkgName> --add-dir <path/to/project>
````

_--rewrite_, _-rv_ — Використовується лише з командою _--add-dir_. Completely rewrites the contents of the package.
````
go run main.go -pkg <pkgName> --add-dir <path/to/project> -rv
````

_--upload_, _-upl_ — Used only with the _--package_ command. Uploads a package to the server.
````
go run main.go -pkg <pkgName> -upl
````

_--download_, _-dwn_ (user string, version string) — Used only with the _--package_ command. Downloads a package from the server to the local machine.
````
go run main.go -pkg <pkgName> -dwn <username> <version>
````