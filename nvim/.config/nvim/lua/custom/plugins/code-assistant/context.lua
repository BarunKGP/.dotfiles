local M = {}

function M.get_context(attachments)
	local context = ""
	for _, att in ipairs(attachments) do
		local content = vim.fn.readfile(att.file)
		local selected_lines = vim.list_slice(content, att.start_line, att.end_line)
		context = context .. att.file .. ":\n" .. table.concat(selected_lines, "\n") .. "\n\n"
	end
	return context
end

return M
