# `hictl`

> A developer tool extension for Go [ent](https://github.com/ent/ent) framework.

-
  1. Add database reverse engineering.

## 1.`Install`

```shell
$ go install github.com/photowey/hictl/cmd/hictl
```

## 2.`Command`

### 2.1.`version`

> show the current `hictl` `cmd ` version.

```shell
$ hictl version
```

### 2.2.`config`

> show the current `hictl` `cmd` `config` info.

```shell
$ hictl config
```

## 2.3.`init`

#### 2.3.1.`create project`

```shell
$ go mod init <project>

$ go mod init userapp
```

#### 2.3.2.`init ent schema `

> `init ` the **`ent`**schema by database reverse engineering.

##### 2.3.2.1.`config database`

> `hictl` home `dir`: `$HOME`/`.hictl`
>
> `hictl` `config` file: `hictl.json`

> original config

```json
{
  "databases": {}
}
```

> databases.users -> the key of flag: --database

```json
{
  "databases": {
    "users": {
      "database": "users",
      "host": "192.168.1.11",
      "port": "3306",
      "username": "root",
      "password": "root"
    }
  }
}
```

##### 2.3.2.2.`reverse engineering`

```shell
# 1.
$ cd ./userapp

$ hictl init -d[--database] users  // full tables
$ hictl init -d users User Group   // only given table user[s] and group[s].
```

