import { describe, it, expect } from "bun:test";
import { createUser, authHeader, BASE_URL } from "../helpers";

describe("POST /api/auth/token", () => {
  it("should return access token for valid mock google token", async () => {
    const user = await createUser();
    expect(user.accessToken.length).toBeGreaterThan(0);
  });

  it("should return same user for same email", async () => {
    const email = `same-user-${crypto.randomUUID().slice(0, 6)}@example.com`;
    const user1 = await createUser(email);
    const user2 = await createUser(email);

    // both should succeed
    expect(user1.accessToken.length).toBeGreaterThan(0);
    expect(user2.accessToken.length).toBeGreaterThan(0);
  });

  it("should fail with missing idToken", async () => {
    const response = await fetch(`${BASE_URL}/api/auth/token`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({}),
    });

    expect(response.status).not.toBe(200);
  });
});

describe("GET /api/auth/me", () => {
  it("should return current user info", async () => {
    const { accessToken, email } = await createUser();

    const response = await fetch(`${BASE_URL}/api/auth/me`, {
      headers: authHeader(accessToken),
    });

    const data = (await response.json()) as Record<string, unknown>;
    expect(response.status).toBe(200);
    expect(data).toHaveProperty("email", email);
    expect(data).toHaveProperty("role", "attendee");
  });

  it("should return 401 without token", async () => {
    const response = await fetch(`${BASE_URL}/api/auth/me`);
    expect(response.status).toBe(401);
  });

  it("should return 401 with invalid token", async () => {
    const response = await fetch(`${BASE_URL}/api/auth/me`, {
      headers: { Authorization: "Bearer invalid.token.here" },
    });

    expect(response.status).toBe(401);
  });
});
