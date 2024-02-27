# Dotfile management and organization
This repo contains scripts that help with dotfile management and organization across a variety of my workspaces.

 ## Clone locally
 On a new machine, reproduce this using git:
    ```git clone https://github.com/BarunKGP/dotfiles.git```

This fetches all dotfiles.
Ideally you would have [GNU stow](https://www.gnu.org/software/stow/) as well to set up symlinks automatically.
If not, there are a few helper scripts to set things up
   
## Setting up Neovim as an IDE
Most server environments I work on are painfully old (some don't even have Git 2.0!)
In such cases, it is best to use an `nvim.appimage` to keep things lightweight and portable.
You can always install the fully decked out version of Neovim if you so choose.
   
### Setting up helper scripts
After cloning this repo, run the `setup_nvim.sh` helper script using
```chmod +x setup_nvim.sh
./setup_nvim.sh
```
   
The script is set to use bash but you could edit the first line to reflect the shell of your choice.
(I have to keep switching between bash, ksh).
It installs the `nvim.appimage`, makes it an executable and pulls in the config files from `kickstart.nvim`.
This also sets up the `$HOME/.config/nvim` directory separately which is needed if you aren't using stow.
It also sets an alias for `nvim` so that you can use Neovim as usual by calling `nvim <dir/file>`
   
---
   
## Other setup actions
This repo contains configurations for:
- tmux
- Neovim
