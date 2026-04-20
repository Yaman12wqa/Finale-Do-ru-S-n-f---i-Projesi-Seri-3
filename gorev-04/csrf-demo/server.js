const crypto = require("crypto");
const express = require("express");
const cookieParser = require("cookie-parser");

const app = express();
const sessions = new Map();

app.use(cookieParser());
app.use(express.urlencoded({ extended: false }));

app.use((req, res, next) => {
  let sessionId = req.cookies.session_id;
  if (!sessionId || !sessions.has(sessionId)) {
    sessionId = crypto.randomBytes(24).toString("hex");
    sessions.set(sessionId, { csrfToken: crypto.randomBytes(32).toString("hex") });
    res.cookie("session_id", sessionId, {
      httpOnly: true,
      sameSite: "lax",
      secure: false
    });
  }
  req.session = sessions.get(sessionId);
  next();
});

function requireCSRF(req, res, next) {
  const submitted = req.body.csrf_token;
  if (!submitted || submitted !== req.session.csrfToken) {
    res.status(403).send("Forbidden: invalid CSRF token");
    return;
  }
  next();
}

app.get("/transfer", (req, res) => {
  res.type("html").send(`
    <h1>Transfer Demo</h1>
    <form method="post" action="/transfer">
      <input type="hidden" name="csrf_token" value="${req.session.csrfToken}">
      <label>To <input name="to" value="alice"></label>
      <label>Amount <input name="amount" value="100"></label>
      <button type="submit">Send</button>
    </form>
  `);
});

app.post("/transfer", requireCSRF, (req, res) => {
  res.type("html").send(`<p>Transfer accepted for ${escapeHTML(req.body.amount)} to ${escapeHTML(req.body.to)}.</p>`);
});

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

app.listen(4004, () => {
  console.log("CSRF demo listening on http://localhost:4004/transfer");
});
