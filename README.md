## <p style="text-align: center;">Strux</p>
Strux is a console client for the [api server](https://github.com/uwine4850/strux_api).

Документація на [українській](https://github.com/uwine4850/strux/blob/master/docs/eng/readme_ua.md).

## Groups of commands
* [Initialization commands](https://github.com/uwine4850/strux/blob/master/docs/eng/init_comm.md)
* [Commands for creating a package](https://github.com/uwine4850/strux/blob/master/docs/eng/create_comm.md)
* [Package manager commands](https://github.com/uwine4850/strux/blob/master/docs/eng/pkg_manager_comm.md)
* [User manager commands](https://github.com/uwine4850/strux/blob/master/docs/eng/user_comm.md)
* [User manager commands](https://github.com/uwine4850/strux/blob/master/docs/eng/info_comm.md)

## Getting started
First, you need to download the client from the repository using the link shown below. After downloading, you will find the file (main.go) in the directory
_strux/cmd/strux_ directory, you will find a file (main.go) for managing the client. **All commands will be executed using this file.**
```
https://github.com/uwine4850/strux
```
To work, the Golang language must already be installed. Or it will need to be installed after the client downloads it.

### Initializing the application
All created packages are stored in the same directory. Therefore, to get started, you need to create it.<br>
To do this, run the following command:
```
go run main.go --init <path>
```
It is recommended to use the absolute path format. This command creates the _strux_pkg_ directory at the selected path.

### Creating a package
To create a package, use the following command:
````
go run main.go --create <package_name>
````
If the command is successful, a new package will appear in the _<init/path>/strux_pkg_ path.<br>
So, the package has been created, and now you need to add data to it. This can be done with the following command:
````
go run main.go --package <package name> --add-dir <path/to/project>
````
The _--add-dir_ flag is the path to the project that is added to the selected package.


### Create a new user
To upload a package to the server, you need to be registered. To do this, run the following command:
````
go run main.go --user --register
````
After entering the required data, a new user will be created.


### Uploading a package to the server
To download the package, run the following command:
````
go run main.go --package <package name> --upload
````
The package is downloaded from the local machine. To download, you need to specify the name of the package, and after running the command
the username and password after running the command.

### Downloading the package
````
go run main.go --package <package name> -dwn <username> <version>
````
The command creates the _StruxDownloads_ directory in your home directory and downloads the package there.
