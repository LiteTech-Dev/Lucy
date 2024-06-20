import path, { resolve } from "path";
import { exec } from "child_process";
import fs from "fs";

export async function getPwd(): Promise<path.ParsedPath> {
    return new Promise((resolve, reject) => {
        exec("pwd", (error, stdout, stderr) => {
            if (stderr != "") {
                console.error(stderr);
            } else if (error) {
                reject(error);
            } else {
                resolve(path.parse(stdout));
            }
        });
    });
}

export function checkMCPMDir(): boolean {
    return fs.existsSync("./.mcpm");
}

export function getMCPMDir(): path.ParsedPath {
    if (checkMCPMDir()) {
        return path.parse("./.mcpm");
    } else throw "";
}
