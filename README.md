# Gop - one-time password agent
Agent can generate only 6 digit time-based passwords.
See [TOTP](https://en.wikipedia.org/wiki/Time-based_one-time_password) to read more about algorithm.

## Usage

    Usage:
        gop [command]
        
        Available Commands:
        add         Add new key
        generate    Generate password and copy to clipboard
        help        Help about any command
        list        List all keys

`gop {add|a} <name> <secret key>` adds a new key to gop.
`name` is a case-sensitive name of a key. 
Two-factor key `<secret key>` are short case-insensitive strings of letters A-Z and
digits 2-7.

`gop {list|l}` prints list of available key's names

`gop {generate|g} <name>` generates a new time-based one-time password,
copy to clipboard and prints it.

### Aliases

There is no such feature as 'aliases' but you can do some workaround:

    gop a cryptoexhange L7TXVT75PU52SE3G
    gop a c L7TXVT75PU52SE3G
    gop g c
    $> 570553
    gop g cryptoexhange
    $> 570553

### Keychain
Keys are stored in a home directory in `.gop.keychain` file.

## License
`gop` is licenced under [MIT License](./LICENSE) 