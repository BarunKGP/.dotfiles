---@type LazyPluginSpec
local cwd = vim.fn.getcwd()
print("Current CWD: " .. cwd)

return {
	dir = vim.fn.stdpath("config") .. "/lua/custom/plugins/code-assistant", -- use relative if within same structure
	name = "code-assistant", -- must be unique
	dev = true, -- tells Lazy to treat this as a local dev plugin
	lazy = false, -- or true if you want to lazy load on command
	dependencies = {
		"nvim-lua/plenary.nvim",
		-- "hrsh7th/nvim-cmp",
	},
	config = function()
		-- -- Optional: any global setup, like setting commands or keymap
		-- vim.api.nvim_create_user_command("LLMChat", function(opts)
		-- 	require("custom.plugins.code-assistant.assistant").ask(opts.args, {}) -- Empty attachments for now
		-- end, { nargs = "*" })

		vim.keymap.set("n", "<leader>xa", function()
			require("code-assistant.assistant").prompt_user()
		end, { desc = "Code Assistant: Chat" })
	end,
}
