function Usage(isfromservertoclient,method, num)
	local result = 1
	if string.sub(method,1,8) == "/solpipe" then
		result = 0
	end
	return result
end