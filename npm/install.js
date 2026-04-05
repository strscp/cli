#!/usr/bin/env node

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const os = require("os");

const VERSION = require("./package.json").version;
const REPO = "strscp/cli";

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

function getBinaryName() {
  const platform = PLATFORM_MAP[process.platform];
  const arch = ARCH_MAP[process.arch];

  if (!platform || !arch) {
    console.error(
      `Unsupported platform: ${process.platform}-${process.arch}`
    );
    process.exit(1);
  }

  const ext = platform === "windows" ? "zip" : "tar.gz";
  return `starscope-cli_${VERSION}_${platform}_${arch}.${ext}`;
}

function download(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          return download(res.headers.location).then(resolve).catch(reject);
        }
        if (res.statusCode !== 200) {
          reject(new Error(`Download failed: HTTP ${res.statusCode}`));
          return;
        }
        const chunks = [];
        res.on("data", (chunk) => chunks.push(chunk));
        res.on("end", () => resolve(Buffer.concat(chunks)));
        res.on("error", reject);
      })
      .on("error", reject);
  });
}

async function main() {
  const binaryName = getBinaryName();
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryName}`;
  const binDir = path.join(__dirname, "bin");

  console.log(`Downloading starscope-cli v${VERSION}...`);

  try {
    const data = await download(url);
    fs.mkdirSync(binDir, { recursive: true });

    const archivePath = path.join(binDir, binaryName);
    fs.writeFileSync(archivePath, data);

    if (binaryName.endsWith(".tar.gz")) {
      execSync(`tar -xzf "${archivePath}" -C "${binDir}"`, { stdio: "pipe" });
    } else {
      execSync(`unzip -o "${archivePath}" -d "${binDir}"`, { stdio: "pipe" });
    }

    fs.unlinkSync(archivePath);

    const binary = path.join(
      binDir,
      process.platform === "win32" ? "starscope-cli.exe" : "starscope-cli"
    );
    if (process.platform !== "win32") {
      fs.chmodSync(binary, 0o755);
    }

    console.log("starscope-cli installed successfully.");
  } catch (err) {
    console.error(`Failed to install starscope-cli: ${err.message}`);
    process.exit(1);
  }
}

main();
