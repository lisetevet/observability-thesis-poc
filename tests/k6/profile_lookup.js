import http from "k6/http";
import { check } from "k6";

const rate = Number(__ENV.RATE || 10);
const duration = __ENV.DURATION || "1m";
const preAllocatedVUs = Number(__ENV.PREALLOCATED_VUS || 20);
const maxVUs = Number(__ENV.MAX_VUS || 50);
const p95ThresholdMs = Number(__ENV.P95_THRESHOLD_MS || 1000);

const baseUrl = __ENV.BASE_URL || "http://localhost:8082";
const requestPath = __ENV.REQUEST_PATH || "/api/v1/profile/chris";
const expectedStatus = Number(__ENV.EXPECTED_STATUS || 200);

export const options = {
  scenarios: {
    profile_lookup: {
      executor: "constant-arrival-rate",
      rate: rate,
      timeUnit: "1s",
      duration: duration,
      preAllocatedVUs: preAllocatedVUs,
      maxVUs: maxVUs,
    },
  },
  thresholds: {
    checks: ["rate>0.99"],
    http_req_duration: [`p(95)<${p95ThresholdMs}`],
  },
};

export default function () {
  const url = `${baseUrl}${requestPath}`;
  const res = http.get(url, {
    tags: {
      endpoint: "profile_lookup",
      request_path: requestPath,
      expected_status: String(expectedStatus),
    },
  });

  check(res, {
    [`status is ${expectedStatus}`]: (r) => r.status === expectedStatus,
    [`duration is below ${p95ThresholdMs}ms`]: (r) =>
      r.timings.duration < p95ThresholdMs,
  });
}