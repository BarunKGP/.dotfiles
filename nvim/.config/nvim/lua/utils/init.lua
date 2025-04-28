------------------------------- Keymaps ------------------------------- 
vim.keymap.set("v", "k", ":m '>-2<CR>gv=gv", {noremap=true, silent=true})
vim.keymap.set("v", "j", ":m '>+1<CR>gv=gv", {noremap=true, silent=true})
