## <p style="text-align: center;">Initialization commands</p>
These commands are used for initial project setup. Namely, to create the necessary directories and database.

_--init_, _-i_ (path string) — Create a package catalog and configure the database<br>
````
go run main.go --init <path>
````
_--new_, _-n_ — Optional command for _--init_.Recreating the package catalog and changes to the database. Перестворення каталогу пакетів та зміни у базі даних.<br>
````
go run main.go --init <path> --new
````
_--database_, _-db_ —Optional command for _--init_. Change the entry responsible for saving the path to the package directory.
````
go run main.go --init <path> -db
````