# gopssh
Utility for execute command at remote servers by ssh

## Usage

1. Clone repo [gopssh](https://github.com/unixzen/gopssh.git)

```
git clone https://github.com/unixzen/gopssh.git
```

2. Build binary

```
cd ./gopssh
go build gopssh.go
```

3. Create config file - `config.yaml`

Example config file:

```
---

ssh_key_path: "/home/user/.ssh/id_rsa"
ssh_user: "root"
ssh_password: "testpassword"
ssh_port: "22"
method: "key"
command: "uptime"
passphrase: "testphrase"
hosts: 
  - host1
  - host2
  - host3
```

4. Start `gopssh`

```
./gopssh
```

Example output:

```
user@localmachine:~/notes/repo/gopssh $ ./gopssh
host1
 14:10:38 up 13 days,  1:52,  1 user,  load average: 0.00, 0.00, 0.00
```

