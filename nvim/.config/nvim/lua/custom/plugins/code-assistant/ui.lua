local M = {}

function M.get_user_input(callback)
	local buf = vim.api.nvim_create_buf(false, true)

	local width = math.floor(vim.o.columns * 0.5)
	local height = 1
	local row = math.floor(vim.o.lines * 0.3)
	local col = math.floor((vim.o.columns - width) / 2)

	local win = vim.api.nvim_open_win(buf, true, {
		relative = "editor",
		width = width,
		height = height,
		row = row,
		col = col,
		style = "minimal",
		border = "rounded",
	})

	vim.api.nvim_buf_set_option(buf, "buftype", "prompt")
	vim.fn.prompt_setprompt(buf, "Ask Assistant: ")

	vim.fn.prompt_setcallback(buf, function(input)
		vim.api.nvim_win_close(win, true)
		if input and input ~= "" then
			callback(input)
		end
	end)

	vim.cmd("startinsert")
end

return M
