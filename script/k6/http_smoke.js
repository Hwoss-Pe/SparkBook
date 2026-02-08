import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    scenarios: {
        steady: {
            executor: "constant-arrival-rate",
            rate: __ENV.RATE ? parseInt(__ENV.RATE, 10) : 50, // 每秒请求数
            timeUnit: "1s",
            duration: __ENV.DURATION || "2m",
            preAllocatedVUs: __ENV.VUS ? parseInt(__ENV.VUS, 10) : 200,
            maxVUs: __ENV.MAX_VUS ? parseInt(__ENV.MAX_VUS, 10) : 2000,
        },
    },
    thresholds: {
        http_req_failed: ["rate<0.05"],
        http_req_duration: ["p(95)<5000"], // ms
    },
};

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080";

// 你自己选接口：支持写死路径，也支持带 query
const PATHS = (__ENV.PATHS || "/__probe__")
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);

function pickPath() {
    return PATHS[Math.floor(Math.random() * PATHS.length)];
}

export default function () {
    const url = `${BASE_URL}${pickPath()}`;

    // 如果你需要带 token/cookie，自己在 ENV 里传
    const headers = {};
    if (__ENV.AUTH) headers["Authorization"] = __ENV.AUTH;

    const res = http.get(url, { headers });

    check(res, {
        "status is not 0": (r) => r.status !== 0,
    });

    // 不要太高，否则会影响你想看的 QPS
    sleep(__ENV.SLEEP ? parseFloat(__ENV.SLEEP) : 0);
}
// k6 run -e BASE_URL=http://localhost:8080 -e PATHS=/articles/pub/ranking,/comment/list,/notifications/unread_counts -e RATE=300 -e DURATION=10m http_smoke.js