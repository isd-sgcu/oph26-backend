import { describe, it, expect } from "bun:test";
import {
  createUser,
  createAttendee,
  authHeader,
  BASE_URL,
} from "../helpers";

describe("POST /api/attendees", () => {
  it("should create a student attendee successfully", async () => {
    const { accessToken } = await createUser();
    const { response, data } = await createAttendee(accessToken);

    expect(response.status).toBe(200);
    expect(data).toHaveProperty("ok", true);
  });

  it("should create parent attendee without interested_faculty", async () => {
    const { accessToken } = await createUser();
    const { response, data } = await createAttendee(accessToken, {
      attendee_type: "parent",
      date_of_birth: "1985-03-20",
      firstname: "วิชัย",
      surname: "สุขสันต์",
      province: "เชียงใหม่",
      district: "เมืองเชียงใหม่",
      news_sources_selected: ["other"],
      news_sources_other: "จากเพื่อน",
      objective_selected: ["preparefordecision", "other"],
      objective_other: "ดูสถานที่ทำงาน",
      transportation_method: "รถยนต์ส่วนตัว",
      interested_faculty: undefined,
      study_level: undefined,
      school_name: undefined,
    });

    if (response.status !== 200) {
      console.log("❌ failed:", data);
    }
    expect(response.status).toBe(200);
    expect(data).toHaveProperty("ok", true);
  });

  it("should fail when student has no interested_faculty", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      interested_faculty: undefined,
    });

    expect(response.status).toBe(400);
  });

  it("should fail with invalid faculty value", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      interested_faculty: ["notafaculty"],
    });

    expect(response.status).toBe(400);
  });

  it("should fail with more than 4 faculties", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      interested_faculty: ["eng", "sci", "md", "law", "econ"],
    });

    expect(response.status).toBe(400);
  });

  it("should fail with invalid province", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      province: "ดาวอังคาร",
    });

    expect(response.status).toBe(400);
  });

  it("should fail with invalid attendee_type", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      attendee_type: "alien",
    });

    expect(response.status).toBe(400);
  });

  it("should fail when creating attendee twice (conflict)", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken);
    const { response } = await createAttendee(accessToken);

    expect(response.status).toBe(409);
  });

  it("should fail with invalid news source", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      news_sources_selected: ["carrier_pigeon"],
    });

    expect(response.status).toBe(400);
  });

  it("should fail when news_sources_selected has 'other' but no news_sources_other", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      news_sources_selected: ["other"],
    });

    expect(response.status).toBe(400);
  });

  it("should fail with invalid objective", async () => {
    const { accessToken } = await createUser();
    const { response } = await createAttendee(accessToken, {
      objective_selected: ["world_domination"],
    });

    expect(response.status).toBe(400);
  });
});

describe("GET /api/attendees/me", () => {
  it("should return attendee data for registered user", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken);

    const response = await fetch(`${BASE_URL}/api/attendees/me`, {
      headers: authHeader(accessToken),
    });

    const data = (await response.json()) as Record<string, unknown>;
    expect(response.status).toBe(200);
    expect(data).toHaveProperty("firstname", "สมชาย");
    expect(data).toHaveProperty("surname", "ใจดี");
    expect(data).toHaveProperty("attendee_type", "student");
    expect(data).toHaveProperty("ticket_code");
    expect(data).toHaveProperty("interested_faculty");
  });

  it("should return 404 for user without attendee data", async () => {
    const { accessToken } = await createUser();

    const response = await fetch(`${BASE_URL}/api/attendees/me`, {
      headers: authHeader(accessToken),
    });

    expect(response.status).toBe(404);
  });
});

describe("PUT /api/attendees/me", () => {
  it("should update attendee firstname", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken);

    const response = await fetch(`${BASE_URL}/api/attendees/me`, {
      method: "PUT",
      headers: authHeader(accessToken),
      body: JSON.stringify({ firstname: "สมหญิง" }),
    });

    expect(response.status).toBe(200);

    // verify the update persisted
    const getRes = await fetch(`${BASE_URL}/api/attendees/me`, {
      headers: authHeader(accessToken),
    });
    const getData = (await getRes.json()) as Record<string, unknown>;
    expect(getData).toHaveProperty("firstname", "สมหญิง");
  });

  it("should update interested_faculty", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken);

    const response = await fetch(`${BASE_URL}/api/attendees/me`, {
      method: "PUT",
      headers: authHeader(accessToken),
      body: JSON.stringify({ interested_faculty: ["md", "law", "arch"] }),
    });

    expect(response.status).toBe(200);

    const getRes = await fetch(`${BASE_URL}/api/attendees/me`, {
      headers: authHeader(accessToken),
    });
    const getData = (await getRes.json()) as Record<string, unknown>;
    expect(getData).toHaveProperty("interested_faculty", ["md", "law", "arch"]);
  });

  it("should return 400 with empty body", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken);

    const response = await fetch(`${BASE_URL}/api/attendees/me`, {
      method: "PUT",
      headers: authHeader(accessToken),
      body: JSON.stringify({}),
    });

    expect(response.status).toBe(400);
  });

  it("should return 404 for user without attendee data", async () => {
    const { accessToken } = await createUser();

    const response = await fetch(`${BASE_URL}/api/attendees/me`, {
      method: "PUT",
      headers: authHeader(accessToken),
      body: JSON.stringify({ firstname: "test" }),
    });

    expect(response.status).toBe(404);
  });

  it("should reject invalid faculty in update", async () => {
    const { accessToken } = await createUser();
    await createAttendee(accessToken);

    const response = await fetch(`${BASE_URL}/api/attendees/me`, {
      method: "PUT",
      headers: authHeader(accessToken),
      body: JSON.stringify({ interested_faculty: ["fake"] }),
    });

    expect(response.status).toBe(400);
  });
});
