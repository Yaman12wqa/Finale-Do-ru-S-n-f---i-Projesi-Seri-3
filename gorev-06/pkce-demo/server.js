const crypto = require("crypto");
const http = require("http");
const { URL, URLSearchParams } = require("url");

const clientId = process.env.GOOGLE_CLIENT_ID;
const clientSecret = process.env.GOOGLE_CLIENT_SECRET;
const redirectURI = process.env.OAUTH_REDIRECT_URI || "http://localhost:4006/callback";
const pending = new Map();

function base64url(buffer) {
  return buffer.toString("base64").replaceAll("+", "-").replaceAll("/", "_").replaceAll("=", "");
}

function randomValue(size = 32) {
  return base64url(crypto.randomBytes(size));
}

function codeChallenge(verifier) {
  return base64url(crypto.createHash("sha256").update(verifier).digest());
}

function send(res, status, body, headers = {}) {
  res.writeHead(status, { "Content-Type": "text/html; charset=utf-8", ...headers });
  res.end(body);
}

function redirect(res, location) {
  res.writeHead(302, { Location: location });
  res.end();
}

async function exchangeCode(code, verifier) {
  const body = new URLSearchParams({
    client_id: clientId,
    client_secret: clientSecret,
    code,
    code_verifier: verifier,
    grant_type: "authorization_code",
    redirect_uri: redirectURI
  });

  const response = await fetch("https://oauth2.googleapis.com/token", {
    method: "POST",
    headers: { "Content-Type": "application/x-www-form-urlencoded" },
    body
  });

  if (!response.ok) {
    throw new Error(`Token exchange failed with HTTP ${response.status}`);
  }
  return response.json();
}

async function getProfile(accessToken) {
  const response = await fetch("https://openidconnect.googleapis.com/v1/userinfo", {
    headers: { Authorization: `Bearer ${accessToken}` }
  });
  if (!response.ok) {
    throw new Error(`Userinfo failed with HTTP ${response.status}`);
  }
  return response.json();
}

const server = http.createServer(async (req, res) => {
  const requestURL = new URL(req.url, "http://localhost:4006");

  if (requestURL.pathname === "/") {
    send(res, 200, `<h1>OAuth PKCE Demo</h1><a href="/login">Login with Google</a>`);
    return;
  }

  if (requestURL.pathname === "/login") {
    if (!clientId || !clientSecret) {
      send(res, 500, "Set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET first.");
      return;
    }

    const state = randomValue(24);
    const verifier = randomValue(48);
    pending.set(state, { verifier, createdAt: Date.now() });

    const params = new URLSearchParams({
      client_id: clientId,
      redirect_uri: redirectURI,
      response_type: "code",
      scope: "openid email profile",
      state,
      code_challenge: codeChallenge(verifier),
      code_challenge_method: "S256"
    });

    redirect(res, `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`);
    return;
  }

  if (requestURL.pathname === "/callback") {
    const code = requestURL.searchParams.get("code");
    const state = requestURL.searchParams.get("state");
    const saved = state ? pending.get(state) : null;
    if (!code || !state || !saved) {
      send(res, 400, "Invalid OAuth callback state.");
      return;
    }
    pending.delete(state);

    try {
      const tokens = await exchangeCode(code, saved.verifier);
      const profile = await getProfile(tokens.access_token);
      send(res, 200, `<h1>Login OK</h1><pre>${escapeHTML(JSON.stringify(profile, null, 2))}</pre>`);
    } catch (err) {
      send(res, 500, escapeHTML(err.message));
    }
    return;
  }

  send(res, 404, "Not found");
});

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

server.listen(4006, () => {
  console.log("OAuth PKCE demo listening on http://localhost:4006");
});
