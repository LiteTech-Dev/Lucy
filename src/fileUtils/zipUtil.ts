import path from "path";
import yauzl from "yauzl";
import { cwd } from "process";
import fs from "fs";
import { FileNotFoundError } from "./fileError.js";

export class ZipUtil {
    private static _instance: ZipUtil;
    // Problem: this.path relies on serverHandler, however, serverHandler will likely call this calss on initialization

    private constructor() {}

    public static getInstance(): ZipUtil {
        if (!ZipUtil._instance) ZipUtil._instance = new ZipUtil();
        return ZipUtil._instance;
    }

    public async readFileFromZip(zip: string, target: string): Promise<string> {
        // Use this method for jar files

        // Check if the file exists
        if (!this.exist(zip)) throw FileNotFoundError;

        const zipPath = path.join(cwd(), zip);
        return new Promise((resolve, reject) => {
            yauzl.open(zipPath, { lazyEntries: true }, (err, zipFile) => {
                if (err) return reject(err);
                zipFile.on("entry", (entry: yauzl.Entry) => {
                    if (entry.fileName === target) {
                        zipFile.openReadStream(entry, (err, readStream) => {
                            if (err) return reject(err);
                            const chunks: Buffer[] = [];
                            readStream.on("data", (chunk) => {
                                chunks.push(chunk);
                            });
                            readStream.on("end", () => {
                                return resolve(
                                    Buffer.concat(chunks).toString()
                                );
                            });
                            readStream.on("error", reject);
                        });
                    } else {
                        zipFile.readEntry();
                    }
                });
                zipFile.readEntry();
            });
        });
    }

    public async fileExistsInZip(
        zip: string,
        target: string
    ): Promise<boolean> {
        if (!this.exist(zip)) return false;

        const zipPath = path.join(cwd(), zip);
        return new Promise((resolve, reject) => {
            yauzl.open(zipPath, { lazyEntries: true }, (err, zipFile) => {
                if (err) return reject(err);
                zipFile.on("entry", (entry: yauzl.Entry) => {
                    if (entry.fileName === target) {
                        return resolve(true);
                    } else {
                        zipFile.readEntry();
                    }
                });
                zipFile.readEntry();
            });
        });
    }

    private exist(...targets: string[]): boolean {
        return targets.every((target) =>
            fs.existsSync(path.join(cwd(), target))
        );
    }
}
