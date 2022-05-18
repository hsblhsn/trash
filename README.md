# Trash

Move your files to the trash (`$HOME/.Trash`).

## Usage

```txt
$ trash -h               

Usage: trash [options] [files...]
moves files to the trash (/Users/hsblhsn/.Trash)

  -I    prompt once before removing any removals
  -i    prompt before every removal
  -r    this flag is ignored
  -f    ignore nonexistent files, never prompt
  -rf   alias for -f
  -v    explain what is being done

```

Example:

```txt
$ trash -rf ~/Desktop/foo.txt ~/Desktop/bar.txt
```

### Status

Work in progress.