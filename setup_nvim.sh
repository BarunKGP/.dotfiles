# !/bin/bash

# get nvim.appimage and make executable
wget https://github.com/neovim/neovim/releases/download/stable/nvim.appimage
chmod u+x nvim.appimage

# Get kickstart.nvim
git clone https://github.com/nvim-lua/kickstart.nvim.git
mv kickstart.nvim $HOME/.config/nvim
export SQUASHFS=0

# # Extract appimage if you don't have FUSE
# ./nvim.appimage --appimage-extract
# export SQUASHFS=1 # Used to set alias in bashrc
# # Launch using
# ./squashfs-root/usr/bin/nvim


echo "Installed kickstart.nvim using Neovim AppImage. You may have to extract the appimage if your system does not have FUSE"
