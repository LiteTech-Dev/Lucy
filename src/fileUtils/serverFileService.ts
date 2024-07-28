import path from "path";
import yauzl from "yauzl";
import { cwd } from "process";
import fs from "fs";

export class ServerFileService {
    private static _instance: ServerFileService;

    private constructor() {}

    public static getInstance(): ServerFileService {
        if (!ServerFileService._instance)
            ServerFileService._instance = new ServerFileService();
        return ServerFileService._instance;
    }

    public async readFileFromZip(zip: string, target: string): Promise<string> {
        // Use this method for jar files
        if (!this.exist(zip)) return "";
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
                                return resolve(chunks.concat().toString());
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

    public async readFrom(target: string): Promise<string> {
        if (!this.exist(target)) return "";
        const targetPath = path.join(cwd(), target);
        return new Promise((resolve, reject) => {
            fs.readFile(targetPath, "utf-8", (err, data) => {
                if (err) reject(err);
                resolve(data);
            });
        });
    }

    public exist(...targets: string[]): boolean {
        return targets.every((target) =>
            fs.existsSync(path.join(cwd(), target))
        );
    }
}
