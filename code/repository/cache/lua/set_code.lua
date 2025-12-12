--key，也就是 code:业务:手机号码
local key = KEYS[1]
--创建使用次数key
local cntKey = key..":cnt"
local code = ARGV[1]
-- 验证码的有效时间是十分钟，600 秒
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    --key存在但是没有过期时间，key冲突了
    return -2
elseif ttl < 540 or ttl ==-2 then
    --    没有发过或者已经过了一分钟
    redis.call("SET", key, code)
    redis.call("EXPIRE", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    --发送频繁
    return -1
end

