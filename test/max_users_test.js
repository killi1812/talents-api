import http from 'k6/http';
import { check, sleep } from 'k6';

// --- Configuration ---
const TALENTS_URL = 'https://192.168.1.110:8011/api';
// TODO: read from env
const AUTH_TOKEN = '';

export const options = {
  executor: 'ramping-vus',
  stages: [
    { duration: '5m', target: 1000 },
    { duration: '2m', target: 1000 },
    { duration: '1m', target: 0 },
  ],
  thresholds: {
    http_req_failed: [{ threshold: 'rate<0.01', abortOnFail: true }],
    http_req_duration: ['p(99)<1000'],
  },
};

const params = {
  headers: {
    'Content-Type': 'application/json',
    'X-API-KEY ': AUTH_TOKEN,
  },
};

export default function () {
  // TODO: random string
  const query = '?q=a';

  const res = http.get(`${TALENTS_URL}/talents/${query}`, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.1);
}
