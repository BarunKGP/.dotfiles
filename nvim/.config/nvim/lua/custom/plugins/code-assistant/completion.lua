local cmp = require("cmp")
local curl = require("plenary.curl")

local source = {}

source.new = function()
	return setmetatable({}, { __index = source })
end

function source:get_debug_name()
	return "llm"
end

function source:is_available()
	return true
end

function source:complete(request, callback)
	local context = request.context.cursor_before_line
	curl.post("https://api.openai.com/v1/chat/completions", {
		headers = {
			["Content-Type"] = "application/json",
			["Authorization"] = "Bearer " .. os.getenv("OPENAI_API_KEY"),
		},
		body = vim.fn.json_encode({
			model = "gpt-4.1",
			messages = {
				{ role = "system", content = "Suggest concise completions for the given code context." },
				{ role = "user",   content = context },
			},
			max_tokens = 50,
		}),
		callback = function(response)
			if response.status ~= 200 then
				return
			end
			local data = vim.fn.json_decode(response.body)
			local completion = data.choices[1].message.content

			callback({ {
				label = completion,
				insertText = completion,
			} })
		end,
	})
end

cmp.register_source("llm", source.new())
