local Job = require("plenary.job")

local M = {}

function M.index_workspace()
	local files = {}
	Job:new({
		command = "rg",
		args = { "--files", "--iglob", "!.git", "." },
		cwd = vim.loop.cwd(),
		on_stdout = function(_, line)
			table.insert(files, line)
		end,
	}):sync()

	local workspace_content = {}
	for _, file in ipairs(files) do
		local content = vim.fn.readfile(file)
		workspace_content[file] = table.concat(content, "\n")
	end
	return workspace_content
end

return M
