--限流对象
local key = KEYS[1]
--窗口大小
local window = tonumber(AGV[1])
--限流阈值
local rate = tonumber(AGV[2])
local now = tonumber(ARGV[3])
--窗口的最小起始时间
local min  = now - window
--把小于这时间的数据都去掉
redis.call("ZREMRANGEBYSCORE",key,'-inf',min)
--统计窗口内还有多少请求
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
if cnt >= threshold then
    -- 执行限流
    return "true"
else
    -- 把 score 和 member 都设置成 now zet是有key对应一个set，并且set里面都有对应的分数
    redis.call('ZADD', key, now, now)
    --给的元素设置过期时间，下次在进来筛选的时候就会把过期的丢掉
    redis.call('PEXPIRE', key, window)
    return "false"
end
