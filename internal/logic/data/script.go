package data

import "github.com/redis/go-redis/v9"

var getMappingScript = redis.NewScript(`
local res = {}
local t = redis.call('hgetall', KEYS[1])
for i, v in pairs(t) do
    if( i%2 == 1 )
    then
        local keyExists = redis.call('exists', 'key:'..v)
        if( keyExists == 0 )
        then
            redis.call('hdel', KEYS[1], v)
        else
			table.insert(res, v)
			table.insert(res, t[i+1])
        end
    end
end
return res
`)
