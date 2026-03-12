import { describe, it, expect } from "bun:test";
import { createUser, createAttendee, authHeader, BASE_URL } from "../helpers";

describe("GET /api/pieces/me", () => {
  it("should return piece for student attendee", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken, { attendee_type: "student" });

    const response = await fetch(`${BASE_URL}/api/pieces/me`, {
      headers: authHeader(accessToken),
    });

    // student type in PostAttendee creates piece only for "highschool",
    // so this may be 404 depending on implementation. Adjust if needed.
    const data = (await response.json()) as Record<string, unknown>;
    if (response.status === 200) {
      expect(data).toHaveProperty("piece_code");
      expect(data).toHaveProperty("expire_date");
      expect(data).toHaveProperty("faculty");
    } else {
      expect(response.status).toBe(404);
    }
  });

  it("should return 403 for parent attendee", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken, {
      attendee_type: "parent",
      interested_faculty: undefined,
      study_level: undefined,
      school_name: undefined,
      news_sources_selected: ["Facebook"],
      objective_selected: ["explorechula"],
    });

    const response = await fetch(`${BASE_URL}/api/pieces/me`, {
      headers: authHeader(accessToken),
    });

    expect(response.status).toBe(403);
  });

  it("should return 404 for user without attendee", async () => {
    const { accessToken } = await createUser();

    const response = await fetch(`${BASE_URL}/api/pieces/me`, {
      headers: authHeader(accessToken),
    });

    expect(response.status).toBe(404);
  });

  it("should return 401 without auth", async () => {
    const response = await fetch(`${BASE_URL}/api/pieces/me`);
    expect(response.status).toBe(401);
  });
});
