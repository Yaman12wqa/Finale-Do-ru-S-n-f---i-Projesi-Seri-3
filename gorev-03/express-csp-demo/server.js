const express = require("express");

const app = express();

app.use((req, res, next) => {
  res.setHeader("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'; base-uri 'self'");
  res.setHeader("X-Content-Type-Options", "nosniff");
  res.setHeader("X-Frame-Options", "DENY");
  next();
});

app.get("/", (req, res) => {
  res.type("html").send(`
    <h1>XSS + CSP Demo</h1>
    <p>Try /vulnerable?name=&lt;script&gt;alert(1)&lt;/script&gt;</p>
    <p>Then try /safe?name=&lt;script&gt;alert(1)&lt;/script&gt;</p>
  `);
});

app.get("/vulnerable", (req, res) => {
  const name = String(req.query.name || "guest");
  res.type("html").send(`<h1>Hello ${name}</h1>`);
});

app.get("/safe", (req, res) => {
  const name = escapeHTML(String(req.query.name || "guest"));
  res.type("html").send(`<h1>Hello ${name}</h1>`);
});

function escapeHTML(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

app.listen(4003, () => {
  console.log("XSS/CSP demo listening on http://localhost:4003");
});
