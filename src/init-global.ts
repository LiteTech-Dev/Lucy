import { exec } from "child_process";
import { dir } from "console";
import * as fs from "fs";
import * as path from "path";

export async function initGlobal(): Promise<void> {
    // 初始化pwd
    await exec("pwd", (err, stdout) => {
        if (err) {
            console.error(err);
            return;
        } else {
            fs.opendir(stdout, (err, dir) => {
                if (err) {
                    console.error(""); // 本地化
                    return;
                } else {
                    pwd = dir;
                }
            });
        }
    });

    //初始化 mcpmDirExists 和 mcpmDirPath
    mcpmDirPath = path.join(pwd.path, mcpmDirName);
    fs.existsSync();
}
