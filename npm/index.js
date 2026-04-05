#!/usr/bin/env node

const { execFileSync } = require("child_process");
const path = require("path");
const fs = require("fs");

const ext = process.platform === "win32" ? ".exe" : "";
const binary = path.join(__dirname, "bin", `starscope-cli${ext}`);

if (!fs.existsSync(binary)) {
  console.error(
    "starscope-cli binary not found. Try reinstalling: npm install -g @strscp/cli"
  );
  process.exit(1);
}

try {
  const result = execFileSync(binary, process.argv.slice(2), {
    stdio: "inherit",
  });
} catch (err) {
  process.exit(err.status || 1);
}
