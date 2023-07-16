## <p style="text-align: center;">Команди менеджеру пакетів</p>
Команди які управляють уже створеними пакетами.

_--package_, _-pkg_ (pkgName string) — Сама лише ця команда не виконує конкретні операції, вона лише вказує до якого
пакету потрібно звертатись.<br>

_--add-dir_, _-ad_ (addPath string) — Додає у вибраний пакет каталог по вказаному шляху.
````
go run main.go -pkg <pkgName> --add-dir <path/to/project>
````

_--rewrite_, _-rv_ — Використовується лише з командою _--add-dir_. Повністю переписує вміст пакету.
````
go run main.go -pkg <pkgName> --add-dir <path/to/project> -rv
````

_--upload_, _-upl_ — Використовується лише з командою _--package_. Завантажує пакет на сервер.
````
go run main.go -pkg <pkgName> -upl
````

_--download_, _-dwn_ (user string, version string) — Використовується лише з командою _--package_. Завантажує пакет із сервер на локальну машину.
````
go run main.go -pkg <pkgName> -dwn <username> <version>
````