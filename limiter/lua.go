package limiter

const luaScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local ttl = tonumber(ARGV[2])

local exist = redis.call("EXISTS", key)

if exist == 0 then
    redis.call("SET", key, 1, "EX", ttl)
    return 0
else
    local val = tonumber(redis.call("GET", key))
    if val and val >= limit then
        return 1
    else
        redis.call("INCR", key)
        return 0
    end
end
`