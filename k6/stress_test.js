import http from "k6/http";
import { check, sleep } from "k6";

// ----------------------
// Test configuration
// ----------------------
export const options = {
  stages: [
    { duration: "20s", target: 500 }, // Ramp-up to 50 users
    { duration: "20s", target: 500 }, // Stay at 50 users
    { duration: "20s", target: 0 }, // Ramp-down to 0
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
    http_req_failed: ["rate<0.05"],
  },
};

// ----------------------
// Base URL
// ----------------------
const BASE_URL = "http://localhost:8080/api";

// ----------------------
// SHARED TOKEN for all virtual users
// ----------------------
let sharedToken = null;

// Setup function runs ONCE before the test
export function setup() {
  console.log("🔑 Getting shared token...");

  // Login once to get token
  const loginPayload = {
    email: "testuser@example.com", // Pre-created test user
    password: "123456", // Pre-created password
  };

  const signinRes = http.post(
    `${BASE_URL}/v1/auth/signin`,
    JSON.stringify(loginPayload),
    { headers: { "Content-Type": "application/json" } }
  );

  if (signinRes.status !== 200) {
    throw new Error(
      `Failed to get token: ${signinRes.status} ${signinRes.body}`
    );
  }

  const token = signinRes.json("token");
  console.log(`✅ Got token: ${token ? "Yes" : "No"}`);

  return { token };
}

// ----------------------
// Main test function - uses SHARED token
// ----------------------
export default function (data) {
  const { token } = data;

  if (!token) {
    console.log("❌ No token available");
    sleep(1);
    return;
  }

  // Health check (public endpoint)
  let healthRes = http.get(`${BASE_URL}/health`);
  check(healthRes, { "health OK": (r) => r.status === 200 });

  // ----------------------
  // TEST PROTECTED APIS WITH SHARED TOKEN
  // ----------------------

  // Get user profile
  // let profileRes = http.get(`${BASE_URL}/v1/users/me`, {
  //   headers: { Authorization: `Bearer ${token}` },
  // });
  // check(profileRes, { "profile OK": (r) => r.status === 200 });

  // Get recommendations
  let recRes = http.get(`${BASE_URL}/v1/recommendation/posts?user_id=test`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  check(recRes, { "recommendations OK": (r) => r.status === 200 });

  // List posts (protected)
  // let postsRes = http.get(`${BASE_URL}/v1/posts/`);
  // check(postsRes, { "posts OK": (r) => r.status === 200 });

  // Public endpoints (no auth needed)
  let publicPostsRes = http.get(`${BASE_URL}/v1/posts/`);
  check(publicPostsRes, { "public posts OK": (r) => r.status === 200 });

  sleep(Math.random() * 2 + 1); // simulate user think time
}
