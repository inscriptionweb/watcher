# Watcher #

Watch local files change according to modify date and copy it remotely through ssh

## Building ##

```bash
make

```

## Tests ##

```bash
make test
```

## Usage ##

```
usage: watcher [<flags>] <local-path> <remote-path>

Flags:
   --help                           Show help.
   --verbose                        Report every operation occuring
   --max-change-time=9s             Maximal change time
   --interval-time=10s              Interval between two check
   --excluded-paths=EXCLUDED-PATHS  Folder to exclude from lookup
                                    separated with comma
   --username=USERNAME              Ssh username
   --host=HOST                      Ssh host
   --key-file=KEY-FILE              Ssh keyfile

Args:
<local-path>   Local pathname to lookup
<remote-path>  Remote pathname to copy data in
```
