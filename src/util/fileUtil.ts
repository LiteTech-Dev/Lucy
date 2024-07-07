import fs from "fs";
import path from "path";
import yauzl from "yauzl";
import { cwd } from "process";
import { CACHE_PATH, LUCY_PATH } from "../constants.js";

export class FileUtil {
    private static instance: FileUtil;

    private constructor() {}

    public static getInstance(): FileUtil {
        if (!FileUtil.instance) FileUtil.instance = new FileUtil();
        return FileUtil.instance;
    }

    async extractFromJar(
        jarPath: string, // This should be a relative path
        targetPath: string
    ): Promise<string> {
        return new Promise((resolve, reject) => {
            const jarPathAbsolute = path.join(cwd(), jarPath);
            yauzl.open(
                jarPathAbsolute,
                { lazyEntries: true },
                (err, zipFile) => {
                    if (err) return reject(err);
                    zipFile.on("entry", (entry: yauzl.Entry) => {
                        if (entry.fileName === targetPath) {
                            zipFile.openReadStream(entry, (err, readStream) => {
                                if (err) return reject(err);
                                const chunks: Buffer[] = [];
                                readStream.on("data", (chunk) => {
                                    chunks.push(chunk);
                                });
                                readStream.on("end", () =>
                                    resolve(Buffer.concat(chunks).toString())
                                );
                                readStream.on("error", reject);
                            });
                        } else {
                            zipFile.readEntry();
                        }
                    });
                    zipFile.readEntry();
                }
            );
        });
    }

    async wirteToProgramPath(
        filePath: string,
        content: string
    ): Promise<void> {}

    async readFromProgramPath(filePath: string): Promise<string> {
        return new Promise((resolve, reject) => {});
    }

    createLucyDirectory(): void {
        try {
            if (!fs.existsSync(LUCY_PATH)) {
                fs.mkdirSync(LUCY_PATH);
            }
        } catch (err) {
            console.log(err);
        }
    }

    createCacheDirectory(): void {
        try {
            if (!fs.existsSync(CACHE_PATH)) {
                fs.mkdirSync(CACHE_PATH);
            }
        } catch (err) {
            console.log(err);
        }
    }
}
