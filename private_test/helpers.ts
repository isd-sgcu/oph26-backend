export const BASE_URL = "http://localhost:8080";

const prefix = crypto.randomUUID().slice(0, 8);
let emailCounter = 0;

export function randomEmail(): string {
  return `test-${prefix}-${++emailCounter}@example.com`;
}

export function generateMockToken(email: string): string {
  const header = btoa(JSON.stringify({ alg: "HS256", typ: "JWT" })).replace(
    /=+$/,
    "",
  );
  const payload = btoa(JSON.stringify({ email })).replace(/=+$/, "");
  return `${header}.${payload}.fake_signature`;
}

export async function createUser(
  email?: string,
): Promise<{ email: string; accessToken: string }> {
  const userEmail = email ?? randomEmail();
  const mockGoogleToken = generateMockToken(userEmail);

  const response = await fetch(`${BASE_URL}/api/auth/token`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ idToken: mockGoogleToken }),
  });

  if (!response.ok) {
    throw new Error(`Auth failed: ${response.status} ${await response.text()}`);
  }

  const data = (await response.json()) as { accessToken: string };
  return { email: userEmail, accessToken: data.accessToken };
}

export function authHeader(accessToken: string) {
  return {
    "Content-Type": "application/json",
    Authorization: `Bearer ${accessToken}`,
  };
}

export async function createAttendee(
  accessToken: string,
  overrides: Record<string, unknown> = {},
) {
  const body = {
    date_of_birth: "2008-05-15",
    attendee_type: "student",
    firstname: "สมชาย",
    surname: "ใจดี",
    province: "กรุงเทพมหานคร",
    district: "บางรัก",
    news_sources_selected: ["Facebook", "Instagram"],
    objective_selected: ["learnaboutfaculties", "explorechula"],
    transportation_method: "รถไฟฟ้า",
    interested_faculty: ["eng", "sci"],
    study_level: "matthayom_plai",
    school_name: "โรงเรียนตัวอย่าง",
    ...overrides,
  };

  const response = await fetch(`${BASE_URL}/api/attendees`, {
    method: "POST",
    headers: authHeader(accessToken),
    body: JSON.stringify(body),
  });

  return { response, data: (await response.json()) as Record<string, unknown> };
}
