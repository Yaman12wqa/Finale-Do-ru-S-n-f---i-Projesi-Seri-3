const crypto = require("crypto");

function hashPassword(password) {
  return crypto.createHash("sha256").update(password, "utf8").digest("hex");
}

async function login(db, email, password) {
  if (typeof email !== "string" || typeof password !== "string") {
    throw new Error("Email and password are required");
  }

  const normalizedEmail = email.trim().toLowerCase();
  const passwordHash = hashPassword(password);

  return new Promise((resolve, reject) => {
    db.get(
      "SELECT id, email FROM users WHERE email = ? AND password_hash = ?",
      [normalizedEmail, passwordHash],
      (err, row) => {
        if (err) {
          reject(err);
          return;
        }
        resolve(row || null);
      }
    );
  });
}

module.exports = { login, hashPassword };
