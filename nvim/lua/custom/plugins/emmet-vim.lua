return {
  {
    'mattn/emmet-vim',
    lazy = false, -- Set to false to ensure it loads immediately and sets up correctly
    config = function()
      -- Optional: Configure Emmet mode (e.g., enable in all modes)
      vim.g.user_emmet_mode = 'a' -- 'a' for all modes
    end,
  },
}
