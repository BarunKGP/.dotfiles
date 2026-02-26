local indexing = require("custom.plugins.code-assistant.indexing")
local context_util = require("custom.plugins.code-assistant.context")
local ui = require("custom.plugins.code-assistant.ui")
local curl = require("plenary.curl")

local M = {}

local api_key = os.getenv("OPENAI_API_KEY")

local function query_llm(prompt)
	local response = curl.post("https://api.openai.com/v1/chat/completions", {
		headers = {
			["Content-Type"] = "application/json",
			["Authorization"] = "Bearer " .. api_key,
		},
		body = vim.fn.json_encode({
			model = "gpt-4o",
			messages = prompt,
			max_tokens = 800,
		}),
	})

	if response.status ~= 200 then
		return "Error: " .. response.status .. "\n" .. response.body
	end

	local data = vim.fn.json_decode(response.body)
	return data.choices[1].message.content
end

function M.ask(question, attachments)
	local workspace = indexing.index_workspace()
	local context = context_util.get_context(attachments or {})

	local system_prompt = "You are a coding assistant with full project context:\n"
	for fname, content in pairs(workspace) do
		system_prompt = system_prompt .. "\nFile: " .. fname .. "\n" .. content .. "\n"
	end

	if context ~= "" then
		system_prompt = system_prompt .. "\nUser attached context:\n" .. context
	end

	local prompt = {
		{ role = "system", content = system_prompt },
		{ role = "user", content = question },
	}

	local result = query_llm(prompt)

	local buf = vim.api.nvim_create_buf(false, true)
	vim.api.nvim_buf_set_lines(buf, 0, -1, false, vim.split(result, "\n"))

	vim.api.nvim_open_win(buf, true, {
		relative = "editor",
		width = math.floor(vim.o.columns * 0.7),
		height = math.floor(vim.o.lines * 0.5),
		col = math.floor(vim.o.columns * 0.15),
		row = math.floor(vim.o.lines * 0.25),
		border = "rounded",
	})
end

function M.prompt_user()
	ui.get_user_input(function(input)
		M.ask(input, {})
	end)
end

return M
