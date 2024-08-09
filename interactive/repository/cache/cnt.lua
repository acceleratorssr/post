local key = KEYS[1]
local field = ARGV[1]
local delta = tonumber(ARGV[2])
local exists = redis.call('EXISTS', key)

if exists == 1 then
    redis.call('HINCRBY', key, field, delta)
else
    redis.call('HSET', key, field, delta)
end
return 1
