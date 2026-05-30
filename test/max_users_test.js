import http from 'k6/http';
import { check, sleep } from 'k6';

// --- Configuration ---
const TALENTS_URL = 'http://192.168.1.110:8011/api';
const AUTH_TOKEN = __ENV.API_KEY;

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
    'X-API-KEY': AUTH_TOKEN,
  },
};

// Random search terms for query variation
const searchTerms = ['strike', 'damage', 'roll', 'skill', 'bonus', 'strength', 'magic', 'combat', 'test'];

export default function () {
  const randomSearch = searchTerms[Math.floor(Math.random() * searchTerms.length)];
  const query = `?q=${randomSearch}&limit=100`;

  // http://192.168.1.110:8011/api/talents?q=roll&limit=100
  const res = http.get(`${TALENTS_URL}/talents${query}`, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.1);
}
