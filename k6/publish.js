import http from 'k6/http';
import {check, sleep} from 'k6';

export const options = {
    scenarios: {
        ramp_up: {
            executor: 'ramping-arrival-rate',
            timeUnit: '0.5s', // 以 1s 为单位增加 vus
            preAllocatedVUs: 2000,
            stages: [
                { target: 100, duration: '10s' },
                { target: 10, duration: '10s' },
                { target: 50, duration: '10s' },
            ],
        },
    },
};

// setup 阶段：登录并获取 token
export function setup() {
    const loginUrl = 'http://127.0.0.1:9301/sso/login';
    const loginPayload = JSON.stringify({
        username: 'testuser',
        password: '123456'
    });

    const loginParams = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const loginRes = http.post(loginUrl, loginPayload, loginParams);

    console.log('Login Response Body: ', loginRes.body);

    let loginBody;
    try {
        loginBody = JSON.parse(loginRes.body);
    } catch (e) {
        console.error('Failed to parse response body as JSON:', e);
        return null;
    }

    check(loginBody, {
        'logged in successfully': (body) => body.status === 200,
    });

    const accessToken = loginBody.data.access_token;

    check(accessToken, {
        'AccessToken exists': () => accessToken !== undefined && accessToken !== null,
    });

    return accessToken;
}

function getRandomChineseChar() {
    const baseChar = 0x4e00;
    const charRange = 0x9fff - baseChar;
    return String.fromCharCode(baseChar + Math.floor(Math.random() * charRange));
}

function getRandomChineseString(length) {
    let result = '';
    for (let i = 0; i < length; i++) {
        result += getRandomChineseChar();
    }
    return result;
}

function createRandomRequestBody() {
    return JSON.stringify({
        id: null,
        title: getRandomChineseString(10),
        content: getRandomChineseString(1000),
    });
}

export default function (accessToken) {
    const url = 'http://127.0.0.1:9301/articles/publish';
    const payload = createRandomRequestBody();

    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`,
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    // sleep(1);
}
