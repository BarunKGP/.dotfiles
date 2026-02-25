# Common shell functions

# Create directory and cd into it
mkdir_and_cd() {
    mkdir -p "$1" && cd "$1"
}

alias mkcd='mkdir_and_cd'
